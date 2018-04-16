package retry

import (
	"strconv"
	"time"

	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
)

type (
	// Config represents the Body Limit configuration
	Config struct {
		Attempts  int      `json:"attempts"`
		Backoff   Duration `json:"backoff"`
		Predicate string   `json:"predicate"`
	}

	// Duration is a wrapper for time.Duration so we can use huma readable configs
	Duration time.Duration
)

// UnmarshalJSON is the implementation of the UnmarshalJSON interface
func (d *Duration) UnmarshalJSON(data []byte) (err error) {
	s := string(data)
	if s == "null" {
		return
	}

	s, err = strconv.Unquote(s)
	if err != nil {
		return
	}

	t, err := time.ParseDuration(s)
	if err != nil {
		return
	}

	*d = Duration(t)
	return
}

func init() {
	plugin.RegisterPlugin("retry", plugin.Plugin{
		Action: setupRetry,
	})
}

func setupRetry(route *proxy.Route, rawConfig plugin.Config) error {
	var config Config
	err := plugin.Decode(rawConfig, &config)
	if err != nil {
		return err
	}

	route.AddInbound(NewRetryMiddleware(config))
	return nil
}
