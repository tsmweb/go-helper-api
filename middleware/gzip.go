package middleware

import (
	"net/http"

	"github.com/xi2/httpgzip"
)

// GZIP compress the http responses.
func GZIP(h http.Handler) http.Handler {
	handler := httpgzip.NewHandler(h, nil)
	return handler
}
