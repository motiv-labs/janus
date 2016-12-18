package janus

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

// CorsMiddleware adds CORS headers to a response.
type CorsMiddleware struct {
	corsMeta CorsMeta
}

// Serve is the middleware method.
func (m *CorsMiddleware) Serve(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("CORS middleware started")

		if !m.corsMeta.Enabled {
			log.Debug("CORS is not enabled for this API")
			handler.ServeHTTP(w, r)
		}

		c := cors.New(cors.Options{
			AllowedOrigins:   m.corsMeta.Domains,
			AllowedMethods:   m.corsMeta.Methods,
			AllowedHeaders:   m.corsMeta.RequestHeaders,
			ExposedHeaders:   m.corsMeta.ExposedHeaders,
			AllowCredentials: true,
		})

		innerMiddleware := negroni.New()
		innerMiddleware.Use(c)
		innerMiddleware.ServeHTTP(w, r)

		log.Debug("CORS inner middleware executed")
		handler.ServeHTTP(w, r)
	})
}
