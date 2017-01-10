package proxy

import (
	"net/http"
)

type Transport interface {
	GetRoundTripper(roundTripper http.RoundTripper) http.RoundTripper
}
