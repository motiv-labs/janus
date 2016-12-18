package middleware

import "net/http"

type Recovery struct {
	recoverFunc func(w http.ResponseWriter, r *http.Request, err interface{})
}

func NewRecovery(recoverFunc func(w http.ResponseWriter, r *http.Request, err interface{})) *Recovery {
	return &Recovery{recoverFunc}
}

func (re *Recovery) Serve(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				re.recoverFunc(w, r, err)
			}
		}()

		handler.ServeHTTP(w, r)
	})
}
