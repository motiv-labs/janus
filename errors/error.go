package errors

type Error struct {
	Code    int
	message string
}

func New(code int, message string) *Error {
	return &Error{code, message}
}

func (e Error) Error() string {
	return e.message
}
