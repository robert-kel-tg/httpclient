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
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"time"
)

func main() {
	client := http.NewClient(
		http.ClientSettings{
			Name:          "Test Client",
			MaxRequests:   100,
			Interval:      time.Duration(3),
			Timeout:       time.Duration(5),
			CountRequests: 3,
			FailureRation: 0.6,
		},
	)

	body := "That's body!"

	resp, err := client.Post(&http.RequestSettings{
		Url:  "https://httpbin.org/post",
		Body: bytes.NewReader([]byte(body)),
	})

	if err != nil {
		logrus.Errorf("Error reading response %v", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logrus.Errorf("Error closing body %v", err)
		}
	}()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Error reading body %v", err)
	}
	logrus.Printf("Printing response body \n %s", string(bodyBytes))
}
```