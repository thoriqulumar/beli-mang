package panics

import (
	"errors"
	"sync"
)

var (
	// ErrorPanic variable used as global error message
	ErrorPanic = errors.New("Panic happened")

	once sync.Once
)

// CaptureGoroutine wrap function call with goroutines and send notification when there's panic inside it
//
// Receives handle function that will be executed on normal condition and recovery function that will be executed in-case there's panic
func CaptureGoroutine(handleFn func(), recoveryFn func()) {
	defer HandlePanic(func(err error) {
		// TODO: log trimStackTree
		//log.Printf("Panic: %+v\n", err)
		//stack := internal.TrimStackTrace(debug.Stack())
		//os.Stderr.Write(stack)
		//internal.PublishError(err, stack, nil)
		recoveryFn()
	})
	handleFn()
}

// HandlePanic with cb
func HandlePanic(action func(error)) {
	// TODO add handing when panic here
}
