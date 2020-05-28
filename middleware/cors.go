package middleware

import (
	"github.com/gorilla/handlers"
	"net/http"
)

// CORS sets the accepted headers, permitted sources and methods accepted by the request.
func CORS(h http.Handler) http.Handler {
	headersOK := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"})
	originsOK := handlers.AllowedOrigins([]string{"*"})
	methodsOK := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT", "DELETE"})

	handler := handlers.CORS(originsOK, headersOK, methodsOK)(h)

	return handler
}
