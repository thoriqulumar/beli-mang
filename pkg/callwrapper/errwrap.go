package callwrapper

import "fmt"

type errWrapper struct {
	err         error
	isTriggerCB bool
	customTags  map[string]string
}

func (ew *errWrapper) Error() string {
	return fmt.Sprintf("callwrapper err:%v, trigger CB: %v, custom tags: %v", ew.err, ew.isTriggerCB, ew.customTags)
}

// WrapErr wraps the error returned from the function call with various informations
// to customize how the callwrapper will treat the error.
//
// isTriggerCB specify whether we want to trigger the circuit breaker or not.
//
// customTags give additional tags to send to the metrics.
func WrapErr(err error, isTriggerCB bool, customTags map[string]string) error {
	return &errWrapper{
		err:         err,
		isTriggerCB: isTriggerCB,
		customTags:  customTags,
	}
}
