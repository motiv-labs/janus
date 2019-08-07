package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/magiconair/properties/assert"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
)

func TestDebugTrace(t *testing.T) {
	format := &b3.HTTPFormat{}

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	tests := []struct {
		testName            string
		debugHeader         string
		expectedTraceHeader bool
	}{
		{
			testName:            "'X-Debug-Trace: secret-key' produces response debug header",
			debugHeader:         "secret-key",
			expectedTraceHeader: true,
		},
		{
			testName:            "'X-Debug-Trace: 0' does not produce response debug header",
			debugHeader:         "0",
			expectedTraceHeader: false,
		},
		{
			testName:            "'X-Debug-Trace: true' does not response debug header",
			debugHeader:         "true",
			expectedTraceHeader: false,
		},
		{
			testName:            "'X-Debug-Trace: ' does not produce response debug header",
			debugHeader:         "",
			expectedTraceHeader: false,
		},
	}

	middleware := DebugTrace(format, "secret-key")

	for _, test := range tests {
		req, _ := http.NewRequest("GET", "http://hello-world", nil)
		req.Header.Add("X-Debug-Trace", test.debugHeader)

		rr := httptest.NewRecorder()
		middleware(testHandler).ServeHTTP(rr, req)
		hasDebugHeader := rr.Header().Get("X-Debug-Trace") != ""
		assert.Equal(t, hasDebugHeader, test.expectedTraceHeader, test.testName)
	}
}
