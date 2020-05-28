package middleware

import (
	"github.com/tsmweb/helper-go/auth"
	"net/http"
)

// Auth validates HTTP requests via token.
type Auth struct {
	jwt *auth.JWT
}

// NewAuth creates a new instance of Auth and receives an jwt.JWT object as a parameter.
func NewAuth(jwt *auth.JWT) *Auth {
	return &Auth{jwt}
}

// RequireTokenAuth performs the middleware function by extracting and validating the request token, if the
// token is valid, the request will follow its flow, if the token is invalid, the unauthorized response will be sent.
func (a *Auth) RequireTokenAuth(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token, err := a.jwt.Token(r)
	if err == nil && token.Valid {
		next(w, r)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
