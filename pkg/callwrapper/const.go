package callwrapper

import "errors"

const (
	// lib version
	libVersion = "1.0"

	// metrics key
	callMetricKey   = "dependency_call"
	configMetricKey = "dependency_call.config"

	// cb type
	cbTypeHystrix = "hystrix_cb"
)

// define unified errors
var (
	ErrBreakerOpen             = errors.New("error breaker open")
	ErrCBConcurrencyLimitReach = errors.New("error cb concurrency limit reach")
	ErrUnknownCallType         = errors.New("error unknown call type")
	ErrUnknownHandlerType      = errors.New("error unknown handler type")
)

// define ddog result tag
const (
	tagResultCtxTimeout         = "result:ctx_timeout"
	tagResultCtxCanceled        = "result:ctx_canceled"
	tagResultUpstreamError      = "result:upstream_error"
	tagResultPanic              = "result:panic"
	tagResultCBOpen             = "result:cb_open"
	tagResultCBError            = "result:cb_error"
	tagResultCBConcurrencyLimit = "result:cb_concurrency_limit"
	tagResultOK                 = "result:ok"
)

// metrics tags key & value
const (
	tagResultKey                   = "result"
	tagResultValCtxTimeout         = "ctx_timeout"
	tagResultValCtxCanceled        = "ctx_canceled"
	tagResultValUpstreamError      = "upstream_error"
	tagResultValPanic              = "panic"
	tagResultValCBOpen             = "cb_open"
	tagResultValCBError            = "cb_error"
	tagResultValCBConcurrencyLimit = "cb_concurrency_limit"
	tagResultValOK                 = "ok"

	tagCachedKey  = "cached"
	tagSuccessKey = "success"

	tagValTrue  = "true"
	tagValFalse = "false"
)

// CallType define call types
type CallType string

const (
	// CallTypeInternal : internal / local services dependencies
	CallTypeInternal CallType = "internal"
	// CallTypeExternal : 3rd party dependencies
	CallTypeExternal CallType = "external"
	// CallTypeDatabase : databases dependencies
	CallTypeDatabase CallType = "database"
	// CallTypeCache : cache dependencies
	CallTypeCache CallType = "cache"
	// CallTypeUnknown : unknown dependencies
	CallTypeUnknown CallType = "unknown"
)

func (c CallType) toString() string {
	return string(c)
}

// TranslateCallType get CallType based on string input
func TranslateCallType(s string) (c CallType, err error) {
	switch s {
	case CallTypeInternal.toString():
		c = CallTypeInternal
	case CallTypeExternal.toString():
		c = CallTypeExternal
	case CallTypeDatabase.toString():
		c = CallTypeDatabase
	case CallTypeCache.toString():
		c = CallTypeCache
	default:
		c = CallTypeUnknown
		err = ErrUnknownCallType
	}

	return
}

// HandlerType define handler types
type HandlerType string

const (
	// HandlerTypeHTTP : call using http method
	HandlerTypeHTTP HandlerType = "http"
	// HandlerTypeGRPC : call using grpc method
	HandlerTypeGRPC HandlerType = "grpc"
	// HandlerTypeUnknown : unknown method handler
	HandlerTypeUnknown HandlerType = "unknown"
)

func (c HandlerType) toString() string {
	return string(c)
}

// TranslateHandlerType get HandlerType based on string input
func TranslateHandlerType(s string) (h HandlerType, err error) { //nolint not to change the HandlerType as public
	switch s {
	case HandlerTypeHTTP.toString():
		h = HandlerTypeHTTP
	case HandlerTypeGRPC.toString():
		h = HandlerTypeGRPC
	default:
		h = HandlerTypeUnknown
		err = ErrUnknownHandlerType
	}

	return
}
