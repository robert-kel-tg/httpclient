# httpclient

* Http Client with cirquit breaker and retry mechanism

## Installation
```bash
go get -u github.com/robertke/httpclient
```

## Usage

```go
package main

import (
  "github.com/robertke/httpclient"
)

func main() {
  client := httpClient.NewClient(
		httpClient.ClientSettings{
			Name:          "Test Client",
			MaxRequests:   100,
			Interval:      time.Duration(3),
			Timeout:       time.Duration(5),
			CountRequests: 3,
			FailureRation: 0.6,
		},
    )
    
    body := "Thats body!"

    resp, err := client.Post(&httpClient.RequestSettings{
		Url:  "https://httpbin.org/post",
		Body: bytes.NewReader([]byte(body)),
	})
}
```