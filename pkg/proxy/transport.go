package proxy

import (
	"net/http"
)

// Transport defines the basic methods for transporters
type Transport interface {
	GetRoundTripper(roundTripper http.RoundTripper) http.RoundTripper
}
