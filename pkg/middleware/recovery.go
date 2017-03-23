package middleware

import "net/http"

// Recovery represents the recovery middleware
type Recovery struct {
	recoverFunc func(w http.ResponseWriter, r *http.Request, err interface{})
}

// NewRecovery creates a new instance of Recovery
func NewRecovery(recoverFunc func(w http.ResponseWriter, r *http.Request, err interface{})) *Recovery {
	return &Recovery{recoverFunc}
}

// Handler is the middleware function
func (re *Recovery) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				re.recoverFunc(w, r, err)
			}
		}()

		handler.ServeHTTP(w, r)
	})
}
