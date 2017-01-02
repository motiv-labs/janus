package errors

// Error is a custom error that implements the `error` interface.
// When creating errors you should provide a code (could be and http status code)
// and a message, this way we can handle the errors in a centralized place.
type Error struct {
	Code    int
	message string
}

// New creates a new instance of Error
func New(code int, message string) *Error {
	return &Error{code, message}
}

func (e Error) Error() string {
	return e.message
}
