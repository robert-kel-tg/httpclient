package http

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/robertke/httpclient/breaker"

	log "github.com/sirupsen/logrus"
)

type (
	Client interface {
		Get(reqSettings *RequestSettings) (*http.Response, error)
		Post(reqSettings *RequestSettings) (*http.Response, error)
	}

	client struct {
		clSettings ClientSettings
		cb         breaker.CircuitBreaker
	}

	RequestSettings struct {
		Url  string
		Body io.Reader
	}

	ClientSettings struct {
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

func NewClient(clSettings ClientSettings) Client {

	log.Infof("initializing breaker: %s", clSettings.Name)
	cb := breaker.CircuitBreaker{
		Name:                   clSettings.Name,
		Timeout:                clSettings.Timeout,
		MaxConcurrentRequests:  clSettings.MaxConcurrentRequests,
		RequestVolumeThreshold: clSettings.RequestVolumeThreshold,
		SleepWindow:            clSettings.SleepWindow,
		ErrorPercentThreshold:  clSettings.ErrorPercentThreshold,
		RetryAttempt:           clSettings.RetryAttempt,
		RetrySleep:             clSettings.RetrySleep,
	}
	breaker.Init(cb)

	return &client{
		clSettings,
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
