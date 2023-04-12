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
	"github.com/robert-kel-tg/httpclient"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"time"
)

func main() {
  	// With default values applied
	client := http.NewClient("TestHttpClient")

  	// Or you can specify one or more values
	client := NewClient(
		"TestHttpClient",
		WithTimeout(3000),
		WithMaxConcurrentRequests(3),
		WithRequestVolumeThreshold(3),
		WithSleepWindow(1000),
		WithErrorPercentThreshold(1),
		WithRetryAttempt(3),
		WithRetrySleep(time.Millisecond*5),
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
