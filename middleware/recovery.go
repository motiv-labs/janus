package nice

import (
	"io"
	
	"github.com/gin-gonic/gin"
	"github.com/urfave/negroni"
)

func Recovery(f func(c *gin.Context, err interface{})) gin.HandlerFunc {
	return RecoveryWithWriter(f, gin.DefaultErrorWriter)
}

func RecoveryWithWriter(f func(c *gin.Context, err interface{}), out io.Writer) negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		defer func() {
			if err := recover(); err != nil {
				f(rw, err)
			}
		}()
		next(rw, r)
	}
}