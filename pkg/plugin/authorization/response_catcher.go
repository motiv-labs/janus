package authorization

import "net/http"

type responseCatcherWriter struct {
	http.ResponseWriter
	status int
	body   []byte
}

func (rcw *responseCatcherWriter) Write(b []byte) (int, error) {
	rcw.body = b
	return rcw.ResponseWriter.Write(b)
}

func (rcw *responseCatcherWriter) WriteHeader(statusCode int) {
	rcw.status = statusCode
	rcw.ResponseWriter.WriteHeader(statusCode)
}
