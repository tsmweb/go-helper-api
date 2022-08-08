package middleware

import (
	"github.com/tsmweb/go-helper-api/auth"
	"net/http"
)

// Auth validates HTTP requests via token.
type Auth interface {
	RequireTokenAuth(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type _auth struct {
	jwt auth.JWT
}

// NewAuth creates a new instance of Auth and receives an jwt.JWT object as a parameter.
func NewAuth(jwt auth.JWT) Auth {
	return &_auth{jwt}
}

// RequireTokenAuth performs the middleware function by extracting and validating the request token, if the
// token is valid, the request will follow its flow, if the token is invalid, the unauthorized response will be sent.
func (a *_auth) RequireTokenAuth(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	_, err := a.jwt.ExtractToken(r)
	if err == nil {
		next(w, r)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
