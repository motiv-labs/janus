package retry

import (
	"strconv"
	"time"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/pkg/errors"
)

const (
	strNull = "null"
)

type (
	// Config represents the Body Limit configuration
	Config struct {
		Attempts  int      `json:"attempts"`
		Backoff   Duration `json:"backoff"`
		Predicate string   `json:"predicate"`
	}

	// Duration is a wrapper for time.Duration so we can use human readable configs
	Duration time.Duration
)

// UnmarshalJSON is the implementation of the UnmarshalJSON interface
func (d *Duration) UnmarshalJSON(data []byte) error {
	s := string(data)
	if s == strNull {
		return errors.New("invalid time duration")
	}

	s, err := strconv.Unquote(s)
	if err != nil {
		return err
	}

	t, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	*d = Duration(t)
	return nil
}

func init() {
	plugin.RegisterPlugin("retry", plugin.Plugin{
		Action: setupRetry,
	})
}

func setupRetry(def *api.Definition, rawConfig plugin.Config) error {
	var config Config
	err := plugin.Decode(rawConfig, &config)
	if err != nil {
		return err
	}

	def.Proxy.AddMiddleware(NewRetryMiddleware(config))
	return nil
}
