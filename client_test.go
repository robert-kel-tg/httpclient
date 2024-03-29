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
		"TestHttpClient",
		WithTimeout(3000),
		WithMaxConcurrentRequests(3),
		WithRequestVolumeThreshold(3),
		WithSleepWindow(1000),
		WithErrorPercentThreshold(1),
		WithRetryAttempt(3),
		WithRetrySleep(time.Millisecond*5),
	)

	res, err := httpClient.Post(
		&RequestSettings{
			Url: testServer.URL,
		},
	)

	if err != nil {
		t.Errorf("error was not expected %v", err)
		t.FailNow()
	}

	body, _ := ioutil.ReadAll(res.Body)
	bodyString := string(body)

	assert.Equal(t, rBody, bodyString, "Expected response body to be %v, got %v", rBody, bodyString)
}

func TestTimeoutClient(t *testing.T) {

	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 3)
	})

	testServer := httptest.NewServer(handlerFunc)

	defer testServer.Close()

	httpClient := NewClient(
		"TestHttpClient",
		WithTimeout(3000),
		WithMaxConcurrentRequests(3),
		WithRequestVolumeThreshold(3),
		WithSleepWindow(2000),
		WithErrorPercentThreshold(1),
		WithRetryAttempt(3),
		WithRetrySleep(time.Millisecond*5),
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
	
	expError := `fallback failed with 'hystrix: timeout'. run error was 'hystrix: timeout'`

	assert.Equal(t, expError, currentErr, "Expected response body to be %v, got %v", expError, currentErr)
}
