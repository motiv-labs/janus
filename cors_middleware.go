package janus

import "github.com/rs/cors"

func NewCorsMiddleware(corsMeta CorsMeta) *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   corsMeta.Domains,
		AllowedMethods:   corsMeta.Methods,
		AllowedHeaders:   corsMeta.RequestHeaders,
		ExposedHeaders:   corsMeta.ExposedHeaders,
		AllowCredentials: true,
	})
}
