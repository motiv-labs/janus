<p align="center">
  <a href="https://hellofresh.com">
    <img width="120" src="https://www.hellofresh.de/images/hellofresh/press/HelloFresh_Logo.png">
  </a>
</p>

# hellofresh/stats-go

[![Build Status](https://travis-ci.org/hellofresh/stats-go.svg?branch=master)](https://travis-ci.org/hellofresh/stats-go)
[![Coverage Status](https://codecov.io/gh/hellofresh/stats-go/branch/master/graph/badge.svg)](https://codecov.io/gh/hellofresh/stats-go)
[![GoDoc](https://godoc.org/github.com/hellofresh/stats-go?status.svg)](https://godoc.org/github.com/hellofresh/stats-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/hellofresh/stats-go)](https://goreportcard.com/report/github.com/hellofresh/stats-go)

> Generic Stats library written in Go

This is generic stats library that we at HelloFresh use in our projects to collect services' stats and then create monitoring
dashboards to track activity and problems.

## Key Features

* Several stats backends:
  * `log` for development environment
  * `statsd` for production (with fallback to `log` if statsd server is not available)
  * `memory` for testing purpose, to track stats operations in unit tests
  * `noop` for environments that do not require any stats gathering
* Fixed metric sections count for all metrics to allow easy monitoring/alerting setup in `grafana`
* Easy to build HTTP requests metrics - timing and count
* Generalise or modify HTTP Requests metric - e.g. skip ID part
* Hook for [`logrus`](https://github.com/sirupsen/logrus) to monitor application error logs

## Installation

```sh
go get -u github.com/hellofresh/stats-go
```

## Usage

### Instance creation

Connection DSN has the following format: `<type>://<connection params>/<connection path>?<connection options>`.

* `<type>` - one of supported backends: `log`, `statsd`, `memory`, `noop`
* `<connection params>` - used for `statsd` backend only, to defining host and port
* `<connection path>` - used for `statsd` backend only, to define prefix/namespace
* `<connection options>` - the following options are available in the query string format:
  * `unicode` - convert unicode metrics to ASCII, default value is `false` as it takes significant memory allocation number

```go
package main

import (
        "os"

        "github.com/hellofresh/stats-go"
)

func main() {
        // client that tries to connect to statsd service, fallback to debug log backend if fails to connect
        statsdClient, _ := stats.NewClient("statsd://statsd-host:8125/my.app.prefix?unicode=true")
        defer statsdClient.Close()

        // debug log backend for stats
        logClient, _ := stats.NewClient("log://")
        defer logClient.Close()

        // memory backend to track operations in unit tests
        memoryClient, _ := stats.NewClient("memory://")
        defer memoryClient.Close()

        // noop backend to ignore all stats
        noopClient, _ := stats.NewClient("noop://")
        defer noopClient.Close()

        // get settings from env to determine backend and prefix
        statsClient, _ := stats.NewClient(os.Getenv("STATS_DSN"))
        defer statsClient.Close()
}
```

### Count metrics manually

```go
import "github.com/hellofresh/stats-go/bucket"

timing := statsClient.BuildTimer().Start()
operation := bucket.MetricOperation{"orders", "order", "create"}
err := orderService.Create(...)
statsClient.TrackOperation("ordering", operation, timing, err == nil)

statsClient.TrackMetric("requests", operation)

ordersInLast24h := orderService.Count(time.Duration(24)*time.Hour)
statsClient.TrackState("ordering", operations, ordersInLast24h)
```

### Track requests metrics with middleware

```go
package main

import (
        "net/http"
        "os"

        "github.com/go-chi/chi"
        "github.com/hellofresh/stats-go"
        "github.com/hellofresh/stats-go/middleware"
)

func main() {
        statsClient := stats.NewClient(os.Getenv("STATS_DSN"))
        defer statsClient.Close()

        r := chi.NewRouter()
        r.Use(middleware.New(statsClient))

        r.GET("/", func(c *gin.Context) {
                // will produce "<prefix>.get.-.-" metric
                c.JSON(http.StatusOK, "I'm producing stats!")
        })

        http.ListenAndServe(":8080", r)
}
```

### Logging

`hellofresh/stats-go` uses default `log` package for debug and error logging.
If you want to use your own logger - `stats-go/log.SetHandler()` is available.

#### Use `github.com/sirupsen/logrus` for logging

```go
package main

import (
    "github.com/hellofresh/stats-go/log"
    "github.com/sirupsen/logrus"
)

func main() {
    log.SetHandler(func(msg string, fields map[string]interface{}, err error) {
    	entry = logrus.WithFields(logrus.Fields(fields))
    	if err == nil {
    		entry.Debug(msg)
    	} else {
    		entry.WithError(err).Error(msg)
    	}
    })

    // do your application stuff
}
```

#### Use `go.uber.org/zap` for logging

```go
package main

import (
	"github.com/hellofresh/stats-go/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	log.SetHandler(func(msg string, fields map[string]interface{}, err error) {
		fieldsLen := len(fields)
		zapFields := make([]zapcore.Field, fieldsLen)
		
		i := 0
		for name, val := range fields {
			zapFields[i] = zap.Any(name, val)
			i++
		}

		if err != nil { 
			logger.Error(msg, zap.Error(err), zapFields...)
		} else {
			logger.Debug(msg, zapFields...)
		}
	})

	// do your application stuff
}
```

### Usage for error logs monitoring using `github.com/sirupsen/logrus`

```go
package foo

import (
        "github.com/hellofresh/stats-go/client"
        "github.com/hellofresh/stats-go/hooks"
        log "github.com/sirupsen/logrus"
)

const sectionErrors = "errors"

func initErrorsMonitoring(statsClient client.Client) {
        hook := hooks.NewLogrusHook(statsClient, sectionErrors)
        log.AddHook(hook)

        // will not produce any metrics
        log.Debug("debug")
        log.Info("info")
        log.Warn("warn")

        // will produce metrics:
        // <section>.<level>.-.-
        // total.<section>
        log.Error("error")
        log.Panic("panic")
        log.Falat("fatal")
}
```

#### Usage in unit tests

```go
package foo

import (
        "github.com/hellofresh/stats-go/client"
        "github.com/hellofresh/stats-go/bucket"
)

const sectionStatsFoo = "foo"

func DoSomeJob(statsClient client.Client) error {
        tt := statsClient.BuildTimer().Start()
        operation := bucket.MetricOperation{"do", "some", "job"}

        result, err := doSomeRealJobHere()
        statsClient.TrackOperation(sectionStatsFoo, operation, tt, result)

        return err
}
```

```go
package foo

import (
        "testing"

        "github.com/hellofresh/stats-go"
        "github.com/hellofresh/stats-go/client"
        "github.com/stretchr/testify/assert"
)

func TestDoSomeJob(t *testing.T) {
        statsClient, _ := stats.NewClient("memory://") 

        err := DoSomeJob(statsClient)
        assert.Nil(t, err)

        statsMemory, _ := statsClient.(*client.Memory)
        assert.Equal(t, 1, len(statsMemory.TimerMetrics))
        assert.Equal(t, "foo-ok.do.some.job", statsMemory.TimerMetrics[0].Bucket)
        assert.Equal(t, 1, statsMemory.CountMetrics["foo-ok.do.some.job"])
}
```

#### Generalise resources by type and stripping resource ID

In some cases you do not need to collect metrics for all unique requests, but a single metric for requests of the similar type,
e.g. access time to concrete users pages does not matter a lot, but average access time is important.
`hellofresh/stats-go` allows HTTP Request metric modification and supports ID filtering out of the box, so
you can get generic metric `get.users.-id-` instead thousands of metrics like `get.users.1`, `get.users.13`,
`get.users.42` etc. that may make your `graphite` suffer from overloading.

To use metric generalisation by second level path ID, you can pass `stats.bucket.HttpMetricNameAlterCallback` instance to
`stats-go//client.Client.SetHttpMetricCallback()`. Also there is a shortcut function `stats-go//bucket.NewHasIDAtSecondLevelCallback()`
that generates a callback handler for `stats-go//bucket.SectionsTestsMap`, and shortcut function `stats-go//bucket.ParseSectionsTestsMap`,
that generates sections test map from string, so you can get these values from config.
It accepts a list of sections with test callback in the following format: `<section>:<test-callback-name>`.
You can use either double colon or new line character as section-callback pairs separator, so all of the following
forms are correct:

* `<section-0>:<test-callback-name-0>:<section-1>:<test-callback-name-1>:<section-2>:<test-callback-name-2>`
* `<section-0>:<test-callback-name-0>\n<section-1>:<test-callback-name-1>\n<section-2>:<test-callback-name-2>`
* `<section-0>:<test-callback-name-0>:<section-1>:<test-callback-name-1>\n<section-2>:<test-callback-name-2>`

Currently the following test callbacks are implemented:

* `true` - second path level is always treated as ID,
  e.g. `/users/13` -> `users.-id-`, `/users/search` -> `users.-id-`, `/users` -> `users.-id-`
* `numeric` - only numeric second path level is interpreted as ID,
  e.g. `/users/13` -> `users.-id-`, `/users/search` -> `users.search`
* `not_empty` - only not empty second path level is interpreted as ID,
  e.g. `/users/13` -> `users.-id-`, `/users` -> `users.-`

You can register your own test callback functions using the `stats-go/bucket.RegisterSectionTest()` function
before parsing sections map from string.

```go
package main

import (
        "net/http"
        "os"

        "github.com/example/app/middleware"
        "github.com/gin-gonic/gin"
        "github.com/hellofresh/stats-go"
        "github.com/hellofresh/stats-go/bucket"
)

func main() {
        // STATS_IDS=users:not_empty:clients:numeric
        sectionsTestsMap, err := bucket.ParseSectionsTestsMap(os.Getenv("STATS_IDS"))
        if err != nil {
                sectionsTestsMap = map[bucket.PathSection]bucket.SectionTestDefinition{}
        }
        statsClient, _ := stats.NewClient(os.Getenv("STATS_DSN"))
        statsClient.SetHTTPMetricCallback(bucket.NewHasIDAtSecondLevelCallback(&bucket.SecondLevelIDConfig{
                HasIDAtSecondLevel:    sectionsTestsMap,
                AutoDiscoverThreshold: 25,
                AutoDiscoverWhiteList: []string{"products"},
        }))
        defer statsClient.Close()

        router := gin.Default()
        router.Use(middleware.NewStatsRequest(statsClient))

        router.GET("/users", func(c *gin.Context) {
                // will produce "<prefix>.get.users.-" metric
                c.JSON(http.StatusOK, "Get the userslist")
        })
        router.GET("/users/:id", func(c *gin.Context) {
                // will produce "<prefix>.get.users.-id-" metric 
                c.JSON(http.StatusOK, "Get the user ID " + c.Params.ByName("id"))
        })
        router.GET("/clients/:id", func(c *gin.Context) {
                // will produce "<prefix>.get.clients.-id-" metric
                c.JSON(http.StatusOK, "Get the client ID " + c.Params.ByName("id"))
        })
        router.GET("/ingredients/:id", func(c *gin.Context) {
                // will produce "<prefix>.get.ingredients.<id>" metric for the first AutoDiscoverThreshold requests
                // and then will produce "<prefix>.get.ingredients.-id-" metric for the rest of requests
                c.JSON(http.StatusOK, "Get the ingredient ID " + c.Params.ByName("id"))
        })
        router.GET("/products/:id", func(c *gin.Context) {
                // will produce "<prefix>.get.products.<id>" metric
                c.JSON(http.StatusOK, "Get the product ID " + c.Params.ByName("id"))
        })

        router.Run(":8080")
}
```

## Contributing

To start contributing, please check [CONTRIBUTING](CONTRIBUTING.md).

## Documentation

* `hellofresh/stats-go` Docs: https://godoc.org/github.com/hellofresh/stats-go
* Go lang: https://golang.org/
