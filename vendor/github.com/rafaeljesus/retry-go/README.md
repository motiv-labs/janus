# retry-go

* Retrying made simple and easy for golang.

## Installation
```bash
go get -u github.com/rafaeljesus/retry-go
```

## Usage

### Do
```go
package main

import (
  "time"

  "github.com/rafaeljesus/retry-go"
)

func main() {
  attempts := 3
  sleepTime := time.Second*2
  if err := retry.Do(func() error {
    return work()
  }, attempts, sleepTime); err != nil {
    // Retry failed
  }
}
```

### DoHTTP
```go
package main

import (
  "time"

  "github.com/rafaeljesus/retry-go"
)

func main() {
  attempts := 3
  sleepTime := time.Second*2
  if err := retry.DoHTTP(func() (*http.Response, error) {
    return makeRequest()
  }, attempts, sleepTime); err != nil {
    // Retry failed
  }
}
```

## Contributing
- Fork it
- Create your feature branch (`git checkout -b my-new-feature`)
- Commit your changes (`git commit -am 'Add some feature'`)
- Push to the branch (`git push origin my-new-feature`)
- Create new Pull Request

## Badges

[![Build Status](https://circleci.com/gh/rafaeljesus/retry-go.svg?style=svg)](https://circleci.com/gh/rafaeljesus/retry-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/rafaeljesus/retry-go)](https://goreportcard.com/report/github.com/rafaeljesus/retry-go)
[![Go Doc](https://godoc.org/github.com/rafaeljesus/retry-go?status.svg)](https://godoc.org/github.com/rafaeljesus/retry-go)

---

> GitHub [@rafaeljesus](https://github.com/rafaeljesus) &nbsp;&middot;&nbsp;
> Medium [@_jesus_rafael](https://medium.com/@_jesus_rafael) &nbsp;&middot;&nbsp;
> Twitter [@_jesus_rafael](https://twitter.com/_jesus_rafael)
