package callwrapper

// MetricsDetail : details of the metrics set by client in wrapper initialization
type MetricsDetail struct {
	SourceName  string // name is name of the service
	DestName    string
	SourceHost  string // host can be consul host / public host / IP
	DestHost    string
	CallType    CallType    // initialized using defined constants (CallTypeInternal, CallTypeExternal, etc)
	HandlerType HandlerType // initialized using defined constants (HandlerTypeHTTP, HandlerTypeGRPC)
	Usecase     string      // provide call usecase (endpoint / grpc function / query usecase)
	CustomTag   *struct {   // custom tag is provided for any specific usecase for debugging / analysis
		Key   string
		Value string
	}
	SlackAlertInfo struct { // currently used for personalized datadog alerting to slack
		ChannelTag string // Format: @slack-{username}-{slack channel}. Try search in Ddog Integrations Tab. Ref:https://docs.datadoghq.com/integrations/slack/?tab=slackapplicationus
		GroupID    string // Get by open slack via web-browser, click your team tag, and see the URL param. It will contains group_id
	}
}

// DatadogClient is the interface of the datadog client used by the callwrapper
//
//go:generate mockgen -package=callwrapper -source=callwrapper.go -destination=callwrapper_mock_test.go
type DatadogClient interface {
	Histogram(name string, value float64, tags []string, rate float64) error
}

// MetricsClient is the interface of metrics client used by the callwrapper.
//
// you can simply use tdk/go/metrics client that support this interface.
// the client currently support both newrelic and datadog
type MetricsClient interface {
	Histogram(name string, value float64, tags map[string]string)
}
