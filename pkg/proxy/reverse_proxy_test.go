package proxy

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestRequest() *http.Request {
	return &http.Request{
		Method:           "",
		URL:              nil,
		Proto:            "",
		ProtoMajor:       0,
		ProtoMinor:       0,
		Header:           nil,
		Body:             nil,
		GetBody:          nil,
		ContentLength:    0,
		TransferEncoding: nil,
		Close:            false,
		Host:             "",
		Form:             nil,
		PostForm:         nil,
		MultipartForm:    nil,
		Trailer:          nil,
		RemoteAddr:       "",
		RequestURI:       "",
		TLS:              nil,
		Cancel:           nil,
		Response:         nil,
	}
}
func TestStripPathWithParams(t *testing.T) {
	t.Run("properly strips path - params and listenPath", func(t *testing.T) {
		req := newTestRequest()
		path := "/prepath/my-service/endpoint"
		listenPath := "/prepath/{service}/*"
		paramNames := []string{"service"}

		old := chiURLParam
		defer func() { chiURLParam = old }()

		chiURLParam = func(r *http.Request, key string) string {
			return "my-service"
		}
		returnPath := stripPathWithParams(req, path, listenPath, paramNames)

		assert.Equal(t, "/endpoint", returnPath)
	})

	t.Run("check that strip logic is correct if value is not in path", func(t *testing.T) {
		req := newTestRequest()
		path := "/prepath/my-service/endpoint"
		listenPath := "/prepath/{service}/*"
		paramNames := []string{"service"}

		old := chiURLParam
		defer func() { chiURLParam = old }()

		chiURLParam = func(r *http.Request, key string) string {
			return "other-value"
		}
		returnPath := stripPathWithParams(req, path, listenPath, paramNames)

		assert.Equal(t, "/my-service/endpoint", returnPath)
	})
}
