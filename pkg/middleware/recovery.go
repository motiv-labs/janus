package middleware

import "net/http"

// NewRecovery creates a new instance of Recovery
func NewRecovery(recoverFunc func(w http.ResponseWriter, r *http.Request, err interface{})) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					recoverFunc(w, r, err)
				}
			}()

			handler.ServeHTTP(w, r)
		})
	}
}
