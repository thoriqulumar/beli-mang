// Package callwrapper wrap HTTP & GRPC call with singleflight, metrics, and circuit breaker
package callwrapper

import (
	"beli-mang/pkg/env"
	"beli-mang/pkg/panics"
	"context"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/sync/singleflight"
)

// Config represent callwrapper configuration
type Config struct {
	// in-memory cache config
	InMemCacheConfig *CacheConfig

	// set call context timeout in miliseconds.
	// Notes:
	// - see the `implementation tips` in the readme for the usage of this option with http call
	// - you only need this option if you don't set the timeout of the context
	CallCtxTimeoutMS int64
	// toggle for single flight
	Singleflight    bool
	HystrixCBConfig *HystrixCBConfig

	// UseCapturePanic will capture panic from given fn.
	UseCapturePanic bool
}

// Wrapper represents a  call wrapper
type Wrapper struct {
	sf *singleflight.Group

	cc              *cache
	cb              iBreaker
	reqTimeout      time.Duration
	metricsDetail   MetricsDetail
	metricsCli      MetricsClient
	metricsTags     map[string]string
	useCapturePanic bool

	cbRunHook        func(error) // for testing
	isWhitelistedErr func(error) bool
}

// NewWrapperWithoutMetric creates a new callwrapper without metric functionality.
// You can use this if you already have your own metric definition
func NewWrapperWithoutMetric(cfg Config) *Wrapper {
	return newWrapper(cfg, MetricsDetail{}, nil, nil)
}

// NewWrapper creates a new callwrapper.
//
// metricsCli is the metrics client for this callwrapper.
// You can simply use tdk/go/metrics client (which already support both newrelic and datadog) for this
// or define your own metrics client.
func NewWrapper(cfg Config, metricsDetail MetricsDetail, metricsCli MetricsClient) *Wrapper {
	return newWrapper(cfg, metricsDetail, nil, metricsCli)
}

// WithErrWhitelist is used to prevent whitelisted errors
// from triggering CB and sending success tag false to metrics.
func (w *Wrapper) WithErrWhitelist(errors ...error) *Wrapper {
	w.isWhitelistedErr = func(e error) bool {
		for _, whitelistErr := range errors {
			if e == whitelistErr {
				return true
			}
		}
		return false
	}
	return w
}

// New creates new callwrapper.
//
// Deprecated: use NewWrapper instead
func New(cfg Config, metricsDetail MetricsDetail, ddogCli DatadogClient) *Wrapper {
	return newWrapper(cfg, metricsDetail, ddogCli, nil)
}

func newWrapper(cfg Config, metricsDetail MetricsDetail, ddogCli DatadogClient, metricsCli MetricsClient) *Wrapper {

	var (
		cc         *cache
		sf         *singleflight.Group
		cb         iBreaker = &emptyBreaker{}
		reqTimeout time.Duration
	)

	// memcache
	if cfg.InMemCacheConfig != nil && cfg.InMemCacheConfig.CacheSize > 0 && cfg.InMemCacheConfig.CacheTTLSec > 0 {
		cc = newCache(cfg.InMemCacheConfig.CacheSize, cfg.InMemCacheConfig.CacheTTLSec)
	}

	// circuit breaker
	// default is hystrix CB
	if cfg.HystrixCBConfig != nil {
		cb = newCep21HystrixCB(*cfg.HystrixCBConfig)
	}

	// single flight
	if cfg.Singleflight {
		sf = &singleflight.Group{}
	}

	// monitoring
	// init tags
	tags := map[string]string{
		"source_name":    metricsDetail.SourceName,
		"dest_name":      metricsDetail.DestName,
		"source_host":    metricsDetail.SourceHost,
		"dest_host":      metricsDetail.DestHost,
		"handler":        metricsDetail.HandlerType.toString(),
		"type":           metricsDetail.CallType.toString(),
		"usecase":        metricsDetail.Usecase,
		"env":            env.ServiceEnv(),
		"slack_group_id": metricsDetail.SlackAlertInfo.GroupID,
	}

	for k, v := range tags {
		if v == "" {
			delete(tags, k)
		}
	}

	// setup custom pre-defined tag
	if metricsDetail.CustomTag != nil {
		// tags = append(tags, metricsDetail.CustomTag.Key+":"+metricsDetail.CustomTag.Value)
		tags[metricsDetail.CustomTag.Key] = metricsDetail.CustomTag.Value
	}

	if cfg.CallCtxTimeoutMS != 0 {
		reqTimeout = time.Duration(cfg.CallCtxTimeoutMS) * time.Millisecond
	}

	wrapper := &Wrapper{
		sf:              sf,
		cc:              cc,
		cb:              cb,
		reqTimeout:      reqTimeout,
		metricsDetail:   metricsDetail,
		metricsCli:      metricsCli,
		metricsTags:     tags,
		useCapturePanic: cfg.UseCapturePanic,
		cbRunHook:       func(error) {},
	}
	// send metrics to track the config
	// we don't enable it in development env because it can conflict with our testing purpose
	if !env.IsDevelopment() {
		go wrapper.trackConfig(cfg, time.Now())
	}

	return wrapper
}

const (
	cachedTag    = "cached:true"
	notCachedTag = "cached:false"
	successTag   = "success:true"
	failedTag    = "success:false"
)

// Call : wraps the func call with everything provided.
// [ `fn` ] Please note that every errors returned by the function `fn`
// will be counted as error attempts on CB. Make sure it's all server-failure errors when using the CB.
// [ `requestKey` ] : this key will be used for singleflight / memcache key
func (w *Wrapper) Call(ctx context.Context, requestKey string, fn func(ctx context.Context) (interface{}, error), opts ...CallOption) (interface{}, error) {
	var (
		res interface{}
		err error
		co  = &callOptions{}
	)

	for _, opt := range opts {
		opt.f(co)
	}

	if w.sf == nil {
		// without singleflight
		res, err = w.callDo(ctx, requestKey, co, fn)
	} else {
		// with singleflight
		res, err = w.callDo(ctx, requestKey, co, func(ctx context.Context) (interface{}, error) {
			res, err, _ := w.sf.Do(requestKey, func() (interface{}, error) { return fn(ctx) })
			return res, err
		})
	}
	return res, err
}

// Call wraps the func call.
func (w *Wrapper) callDo(ctx context.Context, requestKey string, callOpts *callOptions, fn func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	var (
		startTime   = time.Now()
		resultTag   = tagResultOK
		histMapTags = map[string]string{
			tagSuccessKey: tagValTrue,
			tagCachedKey:  tagValFalse,
			tagResultKey:  tagResultValOK,
		}
	)

	// call options
	if callOpts.callTags != nil {
		for k, v := range callOpts.callTags {
			histMapTags[k] = v
		}
	}

	// check for context done before passing through
	select {
	case <-ctx.Done():
		switch ctx.Err() {
		case context.DeadlineExceeded:
			resultTag = tagResultCtxTimeout
			histMapTags[tagResultKey] = tagResultValCtxTimeout
		case context.Canceled:
			resultTag = tagResultCtxCanceled
			histMapTags[tagResultKey] = tagResultValCtxCanceled
		default:
			resultTag = "result:ctx_unknown_err"
			histMapTags[tagResultKey] = "ctx_unknown_err"
		}
		histMapTags[tagSuccessKey] = tagValFalse
		histMapTags[tagCachedKey] = tagValFalse

		w.sendHistogram(callMetricKey, startTime, []string{notCachedTag, failedTag, resultTag}, histMapTags, nil)
		return nil, ctx.Err()
	default:
	}

	// check in cache if exists
	if requestKey != "" && w.cc != nil {
		res, ok := w.cc.Get(requestKey)
		if ok {
			histMapTags[tagCachedKey] = tagValTrue
			w.sendHistogram(callMetricKey, startTime, []string{cachedTag, successTag, resultTag}, histMapTags, nil)
			return res, nil
		}
	}

	var (
		res            interface{}
		errCall        error
		whitelistedErr bool
		cancel         context.CancelFunc
		customTags     map[string]string
	)

	if w.reqTimeout != 0 {
		// setup child context for request timeout
		ctx, cancel = context.WithTimeout(ctx, w.reqTimeout)
		defer cancel()
	}

	err := w.cb.Run(func() error {
		// execute the call

		if w.useCapturePanic {
			// using CaptureGoroutine not because want to run the fn in goroutine,
			// but to utilize the recoveryFn.
			// With this, the client will be noticed by errCall if there's panic.
			panics.CaptureGoroutine(func() {
				res, errCall = fn(ctx)
			}, func() {
				errCall = panics.ErrorPanic
			})
		} else {
			res, errCall = fn(ctx)
		}

		if errCall != nil {
			// If the context is canceled, return nil to not trigger the circuit breaker.
			// context canceled should not trigger the CB because the func itself is not failed.
			// it is just the caller don't want to execute this func anymore
			if isCanceled(ctx, errCall) {
				return nil
			}

			// if error is whitelisted then don't include it in cb
			if w.isWhitelistedErr != nil {
				whitelistedErr = w.isWhitelistedErr(errCall)
			}
			if whitelistedErr {
				return nil
			}

			// if it is wrapped error, only trigger CB if told to do so
			if errWrapped, ok := errCall.(*errWrapper); ok {
				customTags = errWrapped.customTags
				errCall = errWrapped.err
				if errWrapped.isTriggerCB {
					return errCall
				}
				return nil
			}
		}
		return errCall
	})
	w.cbRunHook(err)

	// we previously override the errCall inside the cb.Run
	// now we restore it
	if err == nil {
		err = errCall
	}

	if err != nil {
		switch {
		// errors from cb
		case err == ErrBreakerOpen:
			histMapTags[tagResultKey] = tagResultValCBOpen
			resultTag = tagResultCBOpen
		case err == ErrCBConcurrencyLimitReach:
			histMapTags[tagResultKey] = tagResultValCBConcurrencyLimit
			resultTag = tagResultCBConcurrencyLimit
		case err.Error() != errCall.Error():
			histMapTags[tagResultKey] = tagResultValCBError
			resultTag = tagResultCBError

		// errors from client
		case isDeadlineExceeded(ctx, err):
			histMapTags[tagResultKey] = tagResultValCtxTimeout
			resultTag = tagResultCtxTimeout
		case isCanceled(ctx, err):
			histMapTags[tagResultKey] = tagResultValCtxCanceled
			resultTag = tagResultCtxCanceled
		case err == panics.ErrorPanic:
			histMapTags[tagResultKey] = tagResultValPanic
			resultTag = tagResultPanic
		default:
			histMapTags[tagResultKey] = tagResultValUpstreamError
			resultTag = tagResultUpstreamError
		}

		if !whitelistedErr {
			histMapTags[tagSuccessKey] = tagValFalse
		}
		histMapTags[tagCachedKey] = tagValFalse

		w.sendHistogram(callMetricKey, startTime, []string{notCachedTag, failedTag, resultTag}, histMapTags, customTags)
		return nil, err
	}

	// set to cache if ok
	if requestKey != "" && w.cc != nil {
		w.cc.Set(requestKey, res)
	}

	// submit metrics about call duration
	w.sendHistogram(callMetricKey, startTime, []string{notCachedTag, successTag, resultTag}, histMapTags, customTags)
	return res, nil
}

func (w *Wrapper) trackConfig(cfg Config, start time.Time) {
	var (
		singleFlightEnabled = cfg.Singleflight
		memcacheEnabled     = cfg.InMemCacheConfig != nil
		cbEnabled           bool
		cbConfig            = ""
		tags                = map[string]string{
			"env":                  env.ServiceEnv(),
			"lib_version":          libVersion,
			"source_name":          w.metricsDetail.SourceName,
			"dest_name":            w.metricsDetail.DestName,
			"ctx_timeout":          strconv.FormatInt(cfg.CallCtxTimeoutMS, 10),
			"memcache_enabled":     strconv.FormatBool(memcacheEnabled),
			"singleflight_enabled": strconv.FormatBool(singleFlightEnabled),
		}
	)

	if cfg.HystrixCBConfig != nil {
		cbEnabled = true
		// type_maxConcurrent-errorThresholdPercent-minReqThreshold-OnCloseRollingDur-HalfOpenAttemps-ReqConcurrentSuccess-OnOpenSleepDuration
		cbConfig = fmt.Sprintf("%s_%d-%d-%d-%v-%d-%d-%v", cbTypeHystrix,
			cfg.HystrixCBConfig.MaxConcurrentRequest, cfg.HystrixCBConfig.ErrorThresholdPercentage, cfg.HystrixCBConfig.MinRequestThreshold, cfg.HystrixCBConfig.OnCloseRollingDuration,
			cfg.HystrixCBConfig.HalfOpenAttempts, cfg.HystrixCBConfig.RequiredConcurrentSuccessful, cfg.HystrixCBConfig.OnOpenSleepDuration)
	}

	// append cb tags
	tags["cb_enabled"] = strconv.FormatBool(cbEnabled)
	tags["cb_config"] = cbConfig

	// histogram track initialization time
	w.sendHistogram(configMetricKey, start, mapToArrTags(tags), tags, nil)
}

func (w *Wrapper) sendHistogram(key string, startTime time.Time, arrTags []string, mapTags map[string]string,
	customTags map[string]string) {
	duration := time.Since(startTime).Seconds() * 1000
	if w.metricsCli != nil {
		histTags := make(map[string]string)
		for k, v := range w.metricsTags {
			histTags[k] = v
		}
		for k, v := range mapTags {
			histTags[k] = v
		}
		for k, v := range customTags {
			histTags[k] = v
		}
		w.metricsCli.Histogram(key, duration, histTags)
		return
	}
}

// mapToArrTags convert map[string]string to array of `key:value' string
func mapToArrTags(m map[string]string) []string {
	arr := make([]string, 0, len(m))
	for k, v := range m {
		arr = append(arr, k+":"+v)
	}
	return arr
}

// isCanceled returns true if the call is being canceled
// looking by the err and ctx err
func isCanceled(ctx context.Context, err error) bool {
	return err == context.Canceled || ctx.Err() == context.Canceled
}

// isDeadlineExceeded returns true if the call is timed out
// looking by the err and ctx err
func isDeadlineExceeded(ctx context.Context, err error) bool {
	return err == context.DeadlineExceeded || ctx.Err() == context.DeadlineExceeded
}
