package janus

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

// CorsMiddleware adds CORS headers to a response.
type CorsMiddleware struct {
	*Middleware
}

// ProcessRequest is the middleware method.
func (m *CorsMiddleware) ProcessRequest(rw http.ResponseWriter, req *http.Request) (int, error) {
	log.Debug("CORS middleware started")

	if !m.Spec.CorsMeta.Enabled {
		log.Debug("CORS is not enabled for this API")
		return http.StatusOK, nil
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   m.Spec.CorsMeta.Domains,
		AllowedMethods:   m.Spec.CorsMeta.Methods,
		AllowedHeaders:   m.Spec.CorsMeta.RequestHeaders,
		ExposedHeaders:   m.Spec.CorsMeta.ExposedHeaders,
		AllowCredentials: true,
	})

	innerMiddleware := negroni.New()
	innerMiddleware.Use(c)
	innerMiddleware.ServeHTTP(rw, req)

	log.Debug("CORS inner middleware executed")
	return http.StatusOK, nil
}
