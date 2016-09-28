package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
)

type CorsMiddleware struct {
	*Middleware
}

func (m CorsMiddleware) ProcessRequest(req *http.Request, c *gin.Context) (error, int) {
	m.Logger.Debug("CORS middleware started")

	if !m.Spec.CorsMeta.Enabled {
		m.Logger.Debug("CORS is not enabled")
		return nil, http.StatusOK
	}

	innerMiddleware := cors.Middleware(cors.Config{
		Origins:         strings.Join(m.Spec.CorsMeta.Domains, ","),
		Methods:         strings.Join(m.Spec.CorsMeta.Methods, ","),
		RequestHeaders:  strings.Join(m.Spec.CorsMeta.RequestHeaders, ","),
		ExposedHeaders:  strings.Join(m.Spec.CorsMeta.ExposedHeaders, ","),
		Credentials:     true,
		ValidateHeaders: false,
	})

	innerMiddleware(c)
	m.Logger.Debug("CORS inner middleware executed")

	return nil, http.StatusOK
}
