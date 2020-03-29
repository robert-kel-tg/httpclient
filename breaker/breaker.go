package breaker

import (
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/rafaeljesus/retry-go"
	log "github.com/sirupsen/logrus"
)

type (
	// Holder for CB properties
	CircuitBreaker struct {
		Name                   string
		Timeout                int
		MaxConcurrentRequests  int
		RequestVolumeThreshold int
		SleepWindow            int
		ErrorPercentThreshold  int
		RetryAttempt           int
		RetrySleep             time.Duration
	}
)

// Init applies settings for the circuit breaker
func Init(c CircuitBreaker) {
	hystrix.ConfigureCommand(c.Name, hystrix.CommandConfig{
		Timeout:                c.Timeout,
		RequestVolumeThreshold: c.RequestVolumeThreshold,
		ErrorPercentThreshold:  c.ErrorPercentThreshold,
		MaxConcurrentRequests:  c.MaxConcurrentRequests,
		SleepWindow:            c.SleepWindow,
	})
	hystrix.SetLogger(log.StandardLogger())
}

// Execute wraps the function passed to it with a circuit breaker and a retry
func (c *CircuitBreaker) Execute(f func() error, fallback func(err error) error) chan error {
	errChan := hystrix.Go(c.Name, // the name of the circuit breaker.
		//the inlined func to run inside the breaker.
		func() error {
			err := retry.Do(f, c.RetryAttempt, c.RetrySleep)
			return err
		},
		// the fallback func. logging and return the error.
		fallback,
	)
	return errChan
}

// IsOpen returns true when the circuit breaker is disabled, otherwise it is enabled ( closed or half-open )
func (c *CircuitBreaker) IsOpen() bool {
	cb, _, _ := hystrix.GetCircuit(c.Name)
	return cb.IsOpen()
}
