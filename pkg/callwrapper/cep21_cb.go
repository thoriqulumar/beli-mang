package callwrapper

import (
	"beli-mang/pkg/defaults"
	"context"
	"time"

	"github.com/cep21/circuit/closers/hystrix"
	circuit "github.com/cep21/circuit/v3"
)

const defaultName = "hystrix-circuit"

type hystrixCircuit struct {
	cb *circuit.Circuit
}

// HystrixCBConfig : using cep21/circuit lib as basis. Configs are pretty flexible for every state of circuit breaker
type HystrixCBConfig struct {
	// default values set based on `hystrix default`

	// max concurrent requests that cb can handle, more than that will be rejected.
	// 0 for unlimited concurrent request
	MaxConcurrentRequest int64 `default:"999999999"`

	// ---- configuration for circuit opener

	// % of error that triggers circuit to open
	ErrorThresholdPercentage int64 `default:"80"`
	// minimum request for circuit to open based on `rolling window`.
	// for instance if `rolling duration 10s` then min request in 10s
	MinRequestThreshold int64 `default:"20"`
	// duration needed to reroll circuit stats while circuit is closed
	OnCloseRollingDuration time.Duration `default:"10s"`

	// ---- configuration for circuit closer

	// amount of time, after tripping the circuit, to reject requests before allowing attempts again
	OnOpenSleepDuration time.Duration `default:"5s"`
	// how many attempts to allow per OnOpenSleepDuration
	HalfOpenAttempts int64 `default:"1"`
	// how may consecutive passing requests are required before the circuit is closed
	RequiredConcurrentSuccessful int64 `default:"1"`

	// amount of time the request is expected to finish, if not then it will be considered as timeout
	Timeout time.Duration `default:"4s"`
}

func newCep21HystrixCB(cfg HystrixCBConfig) *hystrixCircuit {
	// set default values based on tag `default`
	defaults.SetDefault(&cfg)

	// initialize hystrix configuration
	configuration := hystrix.Factory{
		// Hystrix open logic is to open the circuit after an % of errors
		ConfigureOpener: hystrix.ConfigureOpener{
			ErrorThresholdPercentage: cfg.ErrorThresholdPercentage,
			RequestVolumeThreshold:   cfg.MinRequestThreshold,
			RollingDuration:          cfg.OnCloseRollingDuration,
			// The default values match what hystrix does by default
		},
		// Hystrix close logic is to sleep then check
		ConfigureCloser: hystrix.ConfigureCloser{
			SleepWindow:                  cfg.OnOpenSleepDuration,
			HalfOpenAttempts:             cfg.HalfOpenAttempts,
			RequiredConcurrentSuccessful: cfg.RequiredConcurrentSuccessful,
			// The default values match what hystrix does by default
		},
	}

	h := circuit.Manager{
		// Tell the manager to use this configuration factory whenever it makes a new circuit
		DefaultCircuitProperties: []circuit.CommandPropertiesConstructor{configuration.Configure},
	}

	// define circuit config.
	config := circuit.Config{
		Execution: circuit.ExecutionConfig{
			MaxConcurrentRequests: cfg.MaxConcurrentRequest,
			Timeout:               cfg.Timeout,
		},
	}

	return &hystrixCircuit{
		// "name" specified here isn't used yet, because we are using 1 circuit for each wrapper only for now
		cb: h.MustCreateCircuit(defaultName, config),
	}
}

func (hc *hystrixCircuit) Run(work func() error) error {
	// bypass context usage for now, context handled on implementation layer
	// not using fallback function for now
	err := hc.cb.Execute(
		context.Background(),
		func(ctx context.Context) error {
			return work()
		}, nil)

	if err != nil {
		// parse errors from cep21
		// normalize using defined error
		if errCB, ok := err.(circuit.Error); ok {
			if errCB.CircuitOpen() {
				return ErrBreakerOpen
			} else if errCB.ConcurrencyLimitReached() {
				return ErrCBConcurrencyLimitReach
			}
		}
	}

	return err
}
