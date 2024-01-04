package authorization

import (
	"errors"
	"io"
	"net"
	"net/http"
	"time"
)

var (
	ErrTimeout = errors.New("timed out")
)

// doGetRequestWithTimeout sends GET http request on passed URL with specified timeout.
func doGetRequestWithTimeout(url string, timeout time.Duration) ([]byte, error) {
	client := http.Client{Timeout: timeout}

	resp, err := client.Get(url)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return nil, ErrTimeout
		}
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
