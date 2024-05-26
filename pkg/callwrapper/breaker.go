package callwrapper

type iBreaker interface {
	Run(work func() error) error
}

type emptyBreaker struct{}

func (emptyBreaker) Run(work func() error) error {
	return work()
}
