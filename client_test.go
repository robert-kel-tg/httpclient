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
		fmt.Fprint(w, bytes.NewBuffer([]byte(rBody)))
	})

	testServer := httptest.NewServer(handlerFunc)
	defer testServer.Close()

	httpClient := NewClient(
		ClientSettings{
			Name:          "Test Client",
			MaxRequests:   100,
			Interval:      time.Duration(3),
			Timeout:       time.Duration(5),
			CountRequests: 3,
			FailureRation: 0.6,
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
