package retry

import (
	"math/rand"
	"net/http"
	"time"
)

var (
	defaultSleep = 500 * time.Millisecond
)

// Func is the function to be executed and eventually retried.
type Func func() error

// HTTPFunc is the function to be executed and eventually retried.
// The only difference from Func is that it expects an *http.Response on the first returning argument.
type HTTPFunc func() (*http.Response, error)

// Do runs the passed function until the number of retries is reached.
// Whenever Func returns err it will sleep and Func will be executed again in a recursive fashion.
// The sleep value is slightly modified on every retry (exponential backoff) to prevent the thundering herd problem (https://en.wikipedia.org/wiki/Thundering_herd_problem).
// If no value is given to sleep it will defaults to 500ms.
func Do(fn Func, retries int, sleep time.Duration) error {
	if sleep == 0 {
		sleep = defaultSleep
	}

	if err := fn(); err != nil {
		retries--
		if retries <= 0 {
			return err
		}

		// preventing thundering herd problem (https://en.wikipedia.org/wiki/Thundering_herd_problem)
		sleep += (time.Duration(rand.Int63n(int64(sleep)))) / 2
		time.Sleep(sleep)

		return Do(fn, retries, 2*sleep)
	}

	return nil
}

// DoHTTP wraps Func and returns *http.Response and error as returning arguments.
func DoHTTP(fn HTTPFunc, retries int, sleep time.Duration) (*http.Response, error) {
	var res *http.Response

	err := Do(func() error {
		var err error
		res, err = fn()
		return err
	}, retries, sleep)

	return res, err
}
