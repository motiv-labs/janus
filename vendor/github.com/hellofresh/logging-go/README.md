<p align="center">
  <a href="https://hellofresh.com">
    <img width="120" src="https://www.hellofresh.de/images/hellofresh/press/HelloFresh_Logo.png">
  </a>
</p>

# hellofresh/logging-go

[![Build Status](https://travis-ci.org/hellofresh/logging-go.svg?branch=master)](https://travis-ci.org/hellofresh/logging-go)
[![Coverage Status](https://coveralls.io/repos/github/hellofresh/logging-go/badge.svg?branch=master)](https://coveralls.io/github/hellofresh/logging-go?branch=master)
[![GoDoc](https://godoc.org/github.com/hellofresh/logging-go?status.svg)](https://godoc.org/github.com/hellofresh/logging-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/hellofresh/logging-go)](https://goreportcard.com/report/github.com/hellofresh/logging-go)

> Generic logging configuration library written in Go

This is generic logging configuration library that we at HelloFresh use in our projects to write applications logs
to different log collecting solutions.

## Key Features

* Uses [logrus](https://github.com/sirupsen/logrus) as logging library
* Allows applying logging configuration from config file or environment variables,
  uses [viper](https://github.com/spf13/viper) under the hood
* The following hooks/writers are available:
  * `stderr`
  * `stdout`
  * `discard` aka `/dev/null`
  * `logstash`
  * `syslog`
  * `graylog`

## Installation

```sh
go get -u github.com/hellofresh/logging-go
```

## Usage

### Standalone logging configuring

#### YAML config file example

```yaml
level: info
format: logstash
formatSettings:
  type: MyService
  ts: RFC3339Nano
writer: stderr
hooks:
- format: logstash
  settings: {network: udp, host: logstash.mycompany.io, port: 8911, type: MyService, ts: RFC3339Nano}
- format: syslog
  settings: {network: udp, host: localhost, port: 514, tag: MyService, facility: LOG_LOCAL0, severity: LOG_INFO}
- format: graylog
  settings: {host: graylog.mycompany.io, port: 9000}
```

#### Environment variable config example

```bash
export LOG_LEVEL="info"
export LOG_FORMAT="logstash"
export LOG_FORMAT_SETTINGS="type=MyService,ts:RFC3339Nano"
export LOG_WRITER="stderr"
export LOG_HOOKS='[{"format":"logstash", "settings":{"type":"MyService","ts":"RFC3339Nano", "network": "udp","host":"logstash.mycompany.io","port": "8911"}},{"format":"syslog","settings":{"network": "udp", "host":"localhost", "port": "514", "tag": "MyService", "facility": "LOG_LOCAL0", "severity": "LOG_INFO"}},{"format":"graylog","settings":{"host":"graylog.mycompany.io","port":"9000"}}]'
```

#### Loading and applying configuration

```go
package main

import (
        "github.com/hellofresh/logging-go"
        log "github.com/sirupsen/logrus"
        "github.com/spf13/viper"
)

func main() {
      logging.InitDefaults(viper.GetViper(), "")
      logConfig, err := logging.Load(viper.GetViper(), "/path/to/config.yml")
      if nil != err {
            panic(err)
      }

      err = logConfig.Apply()
      if nil != err {
            panic(err)
      }
      defer logConfig.Flush()

      log.Info("Logger successfully initialised!")
}
```

### Logging as part of the global config

#### YAML config file example

```yaml
foo: rule
bar: 34
log:
    level: info
    format: json
    output: stderr
    hooks:
    - format: logstash
      settings: {network: udp, host: logstash.mycompany.io, port: 8911, type: MyService, ts: RFC3339Nano}
    - format: syslog
      settings: {network: udp, host: localhost, port: 514, tag: MyService, facility: LOG_LOCAL0, severity: LOG_INFO}
    - format: graylog
      settings: {host: graylog.mycompany.io, port: 9000}
```

#### Environment variable config example

```bash
export APP_FOO="rule"
export APP_BAR="34"
export LOG_LEVEL="info"
export LOG_FORMAT="json"
export LOG_WRITER="stderr"
export LOG_HOOKS='[{"format":"logstash", "settings":{"type":"MyService","ts":"RFC3339Nano", "network": "udp","host":"logstash.mycompany.io","port": "8911"}},{"format":"syslog","settings":{"network": "udp", "host":"localhost", "port": "514", "tag": "MyService", "facility": "LOG_LOCAL0", "severity": "LOG_INFO"}},{"format":"graylog","settings":{"host":"graylog.mycompany.io","port":"9000"}}]'
```

#### Loading and applying configuration

```go
package main

import (
        "github.com/hellofresh/logging-go"
        log "github.com/sirupsen/logrus"
        "github.com/spf13/viper"
)

func init() {
        viper.SetDefault("foo", "foo")
        viper.SetDefault("bar", 42)
        logging.InitDefaults(viper.GetViper(), "log")
}

type AppConfig struct {
        Foo string `envconfig:"APP_FOO"`
        Bar int    `envconfig:"APP_BAR"`

        Log logging.LogConfig
}

func LoadAppConfig() (*AppConfig, error) {
        var instance AppConfig
        
        ...
        
        return &instance, nil
}

func main() {
        appConfig, err := LoadAppConfig()
        if nil != err {
                panic(err)
        }

        err = appConfig.Log.Apply()
        if nil != err {
                panic(err)
        }
        defer appConfig.Log.Flush()

        log.Info("Application successfully initialised!")
}
```

## Configuration

### Base logger

* `level` (env `LOG_LEVEL`, default: `info`): `panic`, `fatal`, `error`, `warn` (`warning`), `info`, `debug` (see [`logrus.ParseLevel()`](https://github.com/sirupsen/logrus/blob/master/logrus.go#L36))
* `format` (env `LOG_FORMAT`, default: `json`):
  * `text` - plain text
  * `json` - all fields encoded into JSON string
  * `logstash` - same as `json` but includes additional logstash fields (e.g. `@version`) and format settings (see bellow)
* `formatSettings` (env `LOG_FORMAT_SETTINGS`):
    * `type` (used only for `logstash` format) - any valid string field that will be added to all log entries
    * `ts` (used only for `logstash` format) - `timestamp` field format, the following values are available: `RFC3339`, `RFC3339Nano` (become `time.RFC3339` and  `time.RFC3339Nano` [`time` package constants](https://golang.org/pkg/time/#pkg-constants))
* `writer` (env `LOG_WRITER`, default: `stderr`): `stderr`, `stdout`, `discard`
* `hooks` (env `LOG_HOOKS`) - each hook has te following fields: `format` and `settings`. Currently te following formats are available:

#### `logstash`

Uses [`github.com/bshuster-repo/logrus-logstash-hook` implementation](https://github.com/bshuster-repo/logrus-logstash-hook)

| Setting   | Required | Description                      |
|-----------|----------|----------------------------------|
| `host`    | **YES**  | Logstash host name or IP address |
| `port`    | **YES**  | Logstash host port               |
| `network` | **YES**  | `udp` or `tcp`                   |
| `type`    | no       | same as `formatSettings.type`    |
| `ts`      | no       | same as `formatSettings.type`    |


#### `syslog`

Not supported on Windows.
Uses [`logstash` implementation of `log/syslog`](https://github.com/Sirupsen/logrus/blob/master/hooks/syslog/syslog.go)

| Setting    | Required | Description                                                                                                      |
|------------|----------|------------------------------------------------------------------------------------------------------------------|
| `host`     | **YES**  | Syslog host name or IP address                                                                                   |
| `port`     | **YES**  | Syslog host port                                                                                                 |
| `network`  | **YES**  | `udp` or `tcp`                                                                                                   |
| `severity` | **YES**  | severity part of [syslog priority](https://golang.org/pkg/log/syslog/#Priority) (`LOG_INFO`, `LOG_NOTICE`, etc.) |
| `facility` | **YES**  | facility part of [syslog priority](https://golang.org/pkg/log/syslog/#Priority) (`LOG_LOCAL0`, `LOG_CRON`, etc.) |
| `tag`      | no       | any valid string that will be sent to syslog as tag                                                              |


#### `graylog`

Uses [`github.com/gemnasium/logrus-graylog-hook` implementation](https://github.com/gemnasium/logrus-graylog-hook)

| Setting | Required | Description                                                                                                                                          |
|---------|----------|------------------------------------------------------------------------------------------------------------------------------------------------------|
| `host`  | **YES**  | Graylog host name or IP address                                                                                                                      |
| `port`  | **YES**  | Graylog host port                                                                                                                                    |
| `async` | no       | send log messages to Graylog in synchronous or asynchronous mode, string value must be [parsable to bool](https://golang.org/pkg/strconv/#ParseBool) |


## Contributing

To start contributing, please check [CONTRIBUTING](CONTRIBUTING.md).

## Documentation

* `hellofresh/logging-go` Docs: https://godoc.org/github.com/hellofresh/logging-go
* Go lang: https://golang.org/
