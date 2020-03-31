# httpclient

* Http Client with Circuit Breaker and Retry mechanisms

## Installation
```bash
go get -u github.com/robertke/httpclient
```

## Usage

```go
package main

import (
	"bytes"
	"github.com/robertke/httpclient"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"time"
)

func main() {
	client := http.NewClient(
		http.ClientSettings{
			Name:          "TestHttpClient",
			Timeout:                3000,
			MaxConcurrentRequests:  3,
			RequestVolumeThreshold: 3,
			SleepWindow:            1000,
			ErrorPercentThreshold:  1,
			RetryAttempt:           3,
			RetrySleep:             time.Millisecond * 5,
		},
	)

	body := "That's body!"

	resp, err := client.Post(&http.RequestSettings{
		Url:  "https://httpbin.org/post",
		Body: bytes.NewReader([]byte(body)),
	})

	if err != nil {
		log.Errorf("Error reading response %v", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Errorf("Error closing body %v", err)
		}
	}()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading body %v", err)
	}
	log.Printf("Printing response body \n %s", string(bodyBytes))
}
```
