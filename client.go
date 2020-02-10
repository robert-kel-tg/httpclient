package http

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/rafaeljesus/retry-go"

	log "github.com/sirupsen/logrus"
	"github.com/sony/gobreaker"
)

type (
	Client interface {
		Get(reqSettings *RequestSettings) (*http.Response, error)
		Post(reqSettings *RequestSettings) (*http.Response, error)
	}

	client struct {
		clSettings ClientSettings
		cb         *gobreaker.CircuitBreaker
	}

	RequestSettings struct {
		Url  string
		Body io.Reader
	}

	ClientSettings struct {
		Name          string
		MaxRequests   uint32
		Interval      time.Duration
		Timeout       time.Duration
		CountRequests uint32
		FailureRatio  float64
		RetryNumber   uint32
		RetryTimeout  time.Duration
	}
)

func NewClient(clSettings ClientSettings) Client {
	cb := gobreaker.NewCircuitBreaker(
		gobreaker.Settings{
			Name:        clSettings.Name,
			MaxRequests: clSettings.MaxRequests,
			Interval:    clSettings.Interval * time.Second,
			Timeout:     clSettings.Timeout,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
				return counts.Requests >= clSettings.CountRequests && failureRatio >= clSettings.FailureRatio
			},
			OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
				log.WithFields(
					log.Fields{
						"name":      name,
						"fromState": from,
						"toState":   to,
					},
				)
			},
		},
	)

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

	body, err := c.cb.Execute(func() (interface{}, error) {
		return retry.DoHTTP(func() (*http.Response, error) {

			var (
				client = &http.Client{
					Timeout: time.Second * 3,
				}
			)

			resp, err := client.Do(req)
			return resp, err
		},
			int(c.clSettings.RetryNumber), c.clSettings.RetryTimeout)
	})

	var response *http.Response

	if r, ok := body.(*http.Response); ok {
		response = r
	}

	return response, err
}
