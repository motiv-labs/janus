package retry

import (
	"errors"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"

	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
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

// MarshalJSON is the implementation of the MarshalJSON interface
func (d *Duration) MarshalJSON() ([]byte, error) {
	s := (*time.Duration)(d).String()
	s = strconv.Quote(s)

	return []byte(s), nil
}

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
		Action:   setupRetry,
		Validate: validateConfig,
	})
}

func setupRetry(def *proxy.RouterDefinition, rawConfig plugin.Config) error {
	var config Config
	err := plugin.Decode(rawConfig, &config)
	if err != nil {
		return err
	}

	def.AddMiddleware(NewRetryMiddleware(config))
	return nil
}

func validateConfig(rawConfig plugin.Config) (bool, error) {
	var config Config
	err := plugin.Decode(rawConfig, &config)
	if err != nil {
		return false, err
	}

	return govalidator.ValidateStruct(config)
}
