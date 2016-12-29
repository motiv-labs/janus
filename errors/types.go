package errors

import "net/http"

var (
	// ErrInvalidID represents an invalid indentifier
	ErrInvalidID = New(http.StatusBadRequest, "Please provide a valid ID")
	// ErrProxyExists occurs when you try to register an already registered proxy
	ErrProxyExists = New(http.StatusBadRequest, "Proxy already registered")
)
