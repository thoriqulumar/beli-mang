package callwrapper

// CallOption defines option can that can be passed to a call
type CallOption struct {
	f func(*callOptions)
}

type callOptions struct {
	callTags map[string]string
}

// WithMetricsTags will add tags to the call
func WithMetricsTags(tags map[string]string) CallOption {
	return CallOption{
		f: func(co *callOptions) {
			co.callTags = tags
		},
	}
}
