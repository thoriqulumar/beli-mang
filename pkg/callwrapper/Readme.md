# callwrapper
Callwrapper is external dependency call wrapper that consist these features:
- [callwrapper](#callwrapper)
    - [Status](#status)
    - [Usage](#usage)
    - [Features](#features)
        - [In-memory-cache](#in-memory-cache)
        - [ContextTimeout](#contexttimeout)
        - [SingleFlight](#singleflight)
        - [CircuitBreaker](#circuitbreaker)
        - [Error Whitelist](#error-whitelist)
        - [Metrics](#metrics)
            - [Custom Metrics](#custom-metrics)
        - [Wrap Error and Tags](#wrap-error-and-tags)
        - [Whitelist Error](#whitelist-error)
    - [Implementation Tips](#implementation-tips)

## Status
- Code: **ready** :white_check_mark:
- Docs: **ready** :white_check_mark:

## Usage

Example of call usecase:
1. Some of big-identical-request might happening at the same time. Put **single flight** might be good idea
2. Number of RPS can't be predicted. Choosing %-err-based **CB** (Hystrix)
3. Response is quite big (can be 200MB-500MB), so prefer to not use the **memcache** for now
4. p99 latency for the call is around `4,5-4,8s`, so setting up timeout to 5s is quite make sense

**On App Init #1 Option ( If Service Owner need Metric Feature )**
```go
// init metrics client

metricsClient := ..... // init metrics client using pkg/metrics package or other client with compatible interface
// init callwrapper

// init callwrapper
cw := callwrapper.NewWrapper(
    //For detail config struct, can see comment in the struct itself
    Config{
      // Toggleable, to turn off just left with false
      Singleflight: true,
      // Usual timeout for context passed
      // - see the `implementation tips` in the readme for the usage of this option with http call
	    // - you only need this option if you don't set the timeout of the context
      CallCtxTimeoutMS: 5000,
      HystrixCBConfig: &HystrixCBConfig{
        // define max concurency limit for requests
        // default is unlimited
        MaxConcurrentRequest: 2000, 
        // % of error that triggers circuit to open
        ErrorThresholdPercentage:     80, // set non-sensitive CB because of RPS can't be predicted
        // minimum request for circuit to open based on `rolling window`.
	// for instance if `rolling duration 10s` then min request in 10s
        MinRequestThreshold:          10,
        // how many attempts to allow per OnOpenSleepDuration
        HalfOpenAttempts:             2,
	// how may consecutive passing requests are required before the circuit is closed
        RequiredConcurrentSuccessful: 2,
	// amount of time, after tripping the circuit, to reject requests before allowing attempts again

        OnOpenSleepDuration:          time.Duration(3000) * time.Millisecond,
      },
    },
    MetricsDetail{ 
      SourceName:  "myservice",
      DestName:    "s3",
      SourceHost:  GetHostname(), // Preferrably, put ip here. If there's not, consul host is okay
      DestHost:    "someofawshost.com",
      CallType:    callwrapper.CallTypeExternal,
      HandlerType: callwrapper.HandlerTypeHTTP,
      Usecase:     "/v1/endpoint_name",
      SlackAlertInfo: struct { // currently used for personalized datadog alerting 
        DdogChannelTag string
        TeamTag    string
      }{
        // Format: @slack-{username}-{slack channel}. Try search in Ddog Integrations Tab. Ref:https://docs.datadoghq.com/integrations/slack/?tab=slackapplicationus
        ChannelTag: "slack_alert", 
        // Get by open slack via web-browser, click your team tag, and see the URL param. It will contains group_id
        TeamTag:    "SBWMSBQJ1", 
      }
	}
    }, metricsClient)
```

**On App Init #2 Option ( If Service Owner has their own metric, use simpler constructor )**
```go
// init callwrapper
cw := callwrapper.NewWrapperWithoutMetric(
    //For detail config struct, can see comment in the struct itself
    Config{
      // Toggleable, to turn off just left with false
      Singleflight: true,
      // Usual timeout for context passed
      // - see the `implementation tips` in the readme for the usage of this option with http call
	    // - you only need this option if you don't set the timeout of the context
      CallCtxTimeoutMS: 5000,
      HystrixCBConfig: &HystrixCBConfig{
        // define max concurency limit for requests
        // default is unlimited
        MaxConcurrentRequest: 2000, 
        // % of error that triggers circuit to open
        ErrorThresholdPercentage:     80, // set non-sensitive CB because of RPS can't be predicted
        // minimum request for circuit to open based on `rolling window`.
	// for instance if `rolling duration 10s` then min request in 10s
        MinRequestThreshold:          10,
        // how many attempts to allow per OnOpenSleepDuration
        HalfOpenAttempts:             2,
	// how may consecutive passing requests are required before the circuit is closed
        RequiredConcurrentSuccessful: 2,
	// amount of time, after tripping the circuit, to reject requests before allowing attempts again

        OnOpenSleepDuration:          time.Duration(3000) * time.Millisecond,
      },
    })
```

**On Call**

```go

// on call
// make sure context is well-set and propagated until the last function
res, err := cw.Call(ctx, fmt.Sprintf("do_action_%d",userID), func(ctx context.Context) (interface{},error) {
	// every errors that returned from this function will be counted as error attempt on CB
    // be careful to return errors. Provide server-failure errors only
    return s3.DownloadBatch(ctx, req) 
})
if err != nil {
	log.Println(err) // might be good to put log here or outside layer. Lib doesn't produce logs
    return nil,err
}
resp, ok := res.([]*s3.S3Object)
if !ok {
	return nil,errors.New("failed parsing response") // you define it
}
return resp, nil
```



## Features

### In-memory-cache
We are using in-memory LRU cache with TTL [library](https://github.com/karlseguin/ccache). Key is using the `requestKey` param specified in `Call` method

### ContextTimeout
Provide context timeout to specify call timeout. Make sure that timeout set here is not exceeding the **parent context**.

### SingleFlight
For single flight using this common [library](https://golang.org/x/sync/singleflight). Key is using the `requestKey` param specified in `Call` method

### CircuitBreaker
For circuit breaker:
1. [Cep21Hystrix](https://github.com/cep21/circuit) (**recommended**)
    - Customizeable wider config options
    - %-based errors to trigger CB open

### Error Whitelist
Whitelist errors that you don't want to be sent as `success:false` in your metrics(e.g. `sql.ErrNoRows`).
### Metrics

Provide built-in metrics in centralized format with fixed key **`tdk.dependency_call`**

**User-defined Metrics**
```go
type MetricsDetail struct {
    SourceName  string 		// name is name of the service
    DestName    string
    SourceHost  string 		// host can be consul host / public host / IP
    DestHost    string
    CallType    CallType    // initialized using defined constants (CallTypeInternal, CallTypeExternal, etc)
    HandlerType HandlerType // initialized using defined constants (HandlerTypeHTTP, HandlerTypeGRPC)
    Usecase     string      // provide call usecase (endpoint / grpc function / query usecase)
    CustomTag   *struct {   // custom tag is provided for any specific usecase for further debugging / analysis
        Key   string
        Value string
    }
    SlackAlertInfo struct { // currently used for personalized datadog alerting to slack
	ChannelTag string // Format: @slack-{username}-{slack channel}.
	GroupID        string // Get by open slack via web-browser, click your team tag, and see the URL param. It will contains group_id
    }
}
```

**Server-provided Metrics**
```
Env -> development / staging / production based on $TKPENV
Status -> success / failed
Cached -> true / false 
Result -> ok / cb_open / cb_error / cb_concurrency_limit / ctx_canceled / ctx_timeout / upstream_error
```

**Metric Example**
```
key : tdk.dependency_call
tags: 
  env: production
  source_name: accounts
  dest_name: postgres_user
  source_host: accounts.service.consul
  dest_host: user_slave.pgbouncer.service.consul
  type: database
  handler: http
  usecase: getUserInfo
  status: failed
  cached: false
  result: upstream_error
```

#### Custom Metrics

While the predefined metrics already a lot, it can't meet all teams requirements.
So we provide `WithMetricsTags` optional func for the team to specify their own call tags
```go
opts = append(opts, callwrapper.WithMetricsTags(map[string]string{
            "agg_usecase": aggUsecase,
        }))
    }
    return cw.Call(ctx, requestKey, fn, opts...)
```
### Wrap Error and Tags

By default (except in the context cancelled case) wrapped call that returns error will trigger the Circuit breaker.
We provide mechanism for the wrapped call to return error but not trigger the CB and add additional tagging.

```go
var specialErr1 := errors.New("this error should not trigger cb")
var specialErr2 := errors.New("this error needs additional tags")

// cw is callwrapper.Wrapper object
result, err := cw.Call(ctx, requestKey, func(ctx context.Context) (interface{}, error) {
  got, err := DoSomething()
  if err == specialErr1 {
    isTriggerCB := false // set to false to not trigger circuit breaker (can be inlined, just for explanation)
    err = callwrapper.WrapErr(err, isTriggerCB, map[string]string{})
  }
  if err == specialErr2 {
    err = callwrapper.WrapErr(err, true, map[string]string{
      "tag1":"hello",
      "tag2":"world",
    })
  }
  return got, err
})
```

### Whitelist Error

Alternatively, you can whitelist errors that you want to exclude from both circuit breaker & metrics tag `success:false`.

```go
// init metrics client
metricsClient := ...
myErr := errors.New("my custom error")

// init new wrapper
cw := callwrapper.NewWrapper(
  Config{}, 
  MetricsDetail{}, 
  metricsClient).
WithErrWhitelist(
  sql.ErrNoRows, 
  redigo.ErrNil, goredis.Nil, // choose according to redis engine that you used
  myErr,
)
```

Upstream calls that will return the whitelisted err will now have `success:true` in their tags, and not trigger CB. Do note that it _will still_ update `result` tag appropriately.

**P.S**
- Always welcome any contributor to `Make It Better` :muscle:
- This whitelist doesn't work if you wrap the error (ex: [wrap](#wrap-error-and-tags) method above). Wrapping the error will cause the error to be unequal when compared

## Implementation Tips

1. Use percentage based CB
   Percentage based CB will be more maintainable because it doesn't need to be adjusted when the call RPS changed,
   unlike the error count based CB.

2. Start with the high CB error percentage.
    - DB & Redis: 90-95%
    - HTTP & gRPC call: 80%

3. All errors will contribute to CB open (except context cancel).
   So for the cases below, don't forget to [wrap](#wrap-error-and-tags) or [whitelist](#whitelist-error) the error
    - err == sql.ErrNoRows
    - redis.IsErrNil(err) == true


4. If you don't want to use memcache, only wrap the actual call.
   As mentioned in the previous point, all of the errors returned by wrapped call will contribute to the CB open. And it is easy for the developer to incidentally include validation mechanism which can return error and forget to wrap it.

5. If it is HTTP call and use `CallCtxTimeoutMS` option, include the reading of response body into the callwrapper.
   Because if you don't, the response body reading will failed with `context canceled` error.

6. Don't assign to variable declared outside of callwrapper function if you use singleflight or memcache.
   Because the wrapped function is not executed if it is served by singleflight or memcache.

don't do this
   ```go
  var result resultStruct // variable outside the callwrapper
  _, err = callwrapper.Call(ctx, key,
	  func(ctx context.Context) (interface{}, error) {
		  err := db.SelectContext(ctx, &result, query, args...)
		  return nil, err
  })
  ```

do this
  ```go
  resp, err = callwrapper.Call(ctx, key,
	  func(ctx context.Context) (interface{}, error) {
		  var result resultStruct
		  err := db.SelectContext(ctx, &result, query, args...)
		  return result, err
  })
  result, _ := resp.(resultStruct)
  ```
7. Test thoroughly both unit test (add multiple test cases positive, negative, extreme, etc) and regression
8. Teams usually has metrics client which prepend the metrics key with their service name.
   Please don't do this to the metrics client of the callwrapper because we want to centralize the callwrapper metrics.
