package http

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/robert-kel-tg/httpclient/breaker"

	log "github.com/sirupsen/logrus"
)

type (
	client struct {
		settings Settings
		cb       breaker.CircuitBreaker
	}

	RequestSettings struct {
		Url  string
		Body io.Reader
	}
)

func NewClient(name string, opts ...Setting) *client {

	var settings Settings

	for _, opt := range opts {
		if err := opt(&settings); err != nil {
			return nil
		}
	}

	log.Infof("initializing breaker: %s", name)
	cb := breaker.CircuitBreaker{
		Name:                   name,
		Timeout:                settings.Timeout,
		MaxConcurrentRequests:  settings.MaxConcurrentRequests,
		RequestVolumeThreshold: settings.RequestVolumeThreshold,
		SleepWindow:            settings.SleepWindow,
		ErrorPercentThreshold:  settings.ErrorPercentThreshold,
		RetryAttempt:           settings.RetryAttempt,
		RetrySleep:             settings.RetrySleep,
	}
	breaker.Init(cb)

	return &client{
		settings,
		cb,
	}
}

func (c *client) Get(reqSettings *RequestSettings) (*http.Response, error) {
	return c.do(http.MethodGet, reqSettings)
}

func (c *client) Post(reqSettings *RequestSettings) (*http.Response, error) {
	return c.do(http.MethodPost, reqSettings)
}

func (c *client) do(method string, reqSettings *RequestSettings) (*http.Response, error) {
	req, err := http.NewRequest(method, reqSettings.Url, reqSettings.Body)
	if err != nil {
		return nil, errors.New("request error")
	}

	output := make(chan *http.Response, 1)
	errorCh := c.cb.Execute(func() error {
		var (
			client = &http.Client{
				Timeout: time.Second * 3,
			}
		)
		resp, err := client.Do(req)
		if nil != err {
			return err
		}

		output <- resp
		return nil
	}, func(err error) error {
		if nil != err {
			log.Errorf("In fallback function for breaker, error: %v", err)
		}
		return err
	})

	select {
	case out := <-output:
		return out, nil
	case err := <-errorCh:
		return nil, err
	}
}
