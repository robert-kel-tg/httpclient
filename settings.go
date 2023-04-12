package http

import (
	"errors"
	"time"
)

type (
	Settings struct {
		Timeout                int
		MaxConcurrentRequests  int
		RequestVolumeThreshold int
		SleepWindow            int
		ErrorPercentThreshold  int
		RetryAttempt           int
		RetrySleep             time.Duration
	}

	Setting func(settings *Settings) error
)

func WithTimeout(timeout int) Setting {
	return func(settings *Settings) error {
		if timeout < 0 {
			return errors.New("timeout should not be negative")
		}
		settings.Timeout = timeout
		return nil
	}
}

func WithMaxConcurrentRequests(maxConcrReqs int) Setting {
	return func(settings *Settings) error {
		if maxConcrReqs < 0 {
			return errors.New("max requests should not be negative")
		}
		settings.MaxConcurrentRequests = maxConcrReqs
		return nil
	}
}

func WithRequestVolumeThreshold(reqVolThreshold int) Setting {
	return func(settings *Settings) error {
		if reqVolThreshold < 0 {
			return errors.New("threshold should not be negative")
		}
		settings.RequestVolumeThreshold = reqVolThreshold
		return nil
	}
}

func WithSleepWindow(sleepWindow int) Setting {
	return func(settings *Settings) error {
		if sleepWindow < 0 {
			return errors.New("sleepWindow should not be negative")
		}
		settings.SleepWindow = sleepWindow
		return nil
	}
}

func WithErrorPercentThreshold(errPercentThreshold int) Setting {
	return func(settings *Settings) error {
		if errPercentThreshold < 0 {
			return errors.New("errPercentThreshold should not be negative")
		}
		settings.ErrorPercentThreshold = errPercentThreshold
		return nil
	}
}

func WithRetryAttempt(retryAttempt int) Setting {
	return func(settings *Settings) error {
		if retryAttempt < 0 {
			return errors.New("retryAttempt should not be negative")
		}
		settings.RetryAttempt = retryAttempt
		return nil
	}
}

func WithRetrySleep(retrySleep time.Duration) Setting {
	return func(settings *Settings) error {
		settings.RetrySleep = retrySleep
		return nil
	}
}
