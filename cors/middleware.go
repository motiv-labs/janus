package cors

import "github.com/rs/cors"

func NewMiddleware(corsMeta Meta, debug bool) *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   corsMeta.Domains,
		AllowedMethods:   corsMeta.Methods,
		AllowedHeaders:   corsMeta.RequestHeaders,
		ExposedHeaders:   corsMeta.ExposedHeaders,
		AllowCredentials: true,
		Debug:            debug,
	})
}
