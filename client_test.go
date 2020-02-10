package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {

	rBody := `{"title":"Hi!"}`

	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, bytes.NewBuffer([]byte(rBody)))
	})

	testServer := httptest.NewServer(handlerFunc)
	defer testServer.Close()

	httpClient := NewClient(
		ClientSettings{
			Name:          "Test Client",
			MaxRequests:   1,
			Interval:      time.Duration(3),
			Timeout:       time.Duration(3),
			CountRequests: 1,
			FailureRatio:  0.6,
			RetryNumber:   3,
			RetryTimeout:  time.Millisecond * 5,
		},
	)

	res, _ := httpClient.Post(
		&RequestSettings{
			Url: testServer.URL,
		},
	)

	body, _ := ioutil.ReadAll(res.Body)
	bodyString := string(body)

	assert.Equal(t, rBody, bodyString, "Expected response body to be %v, got %v", rBody, bodyString)
}

func TestTimeoutClient(t *testing.T) {

	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 3)
	})

	testServer := httptest.NewServer(handlerFunc)

	expError := `Post ` + testServer.URL + `: net/http: request canceled (Client.Timeout exceeded while awaiting headers)`

	defer testServer.Close()

	httpClient := NewClient(
		ClientSettings{
			Name:          "Test Client",
			MaxRequests:   1,
			Interval:      time.Duration(3),
			Timeout:       time.Duration(3),
			CountRequests: 1,
			FailureRatio:  0.6,
			RetryNumber:   3,
			RetryTimeout:  time.Millisecond * 5,
		},
	)

	_, err := httpClient.Post(
		&RequestSettings{
			// https://httpbin.org/delay/5000
			//Url: "https://httpstat.us/200?sleep=5000",
			Url: testServer.URL,
		},
	)

	var (
		currentErr string
	)

	if err != nil {
		currentErr = err.Error()
	}

	assert.Equal(t, expError, currentErr, "Expected response body to be %v, got %v", expError, currentErr)
}
