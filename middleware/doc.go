/*
Package middleware provides settings for CORS, GZIP and JWT token validation.

CORS sets the accepted headers, permitted sources and methods accepted by the request:

	h := middleware.CORS(h)
	// ...

GZIP compress the http responses:

	h := middleware.GZIP(h)
	// ...

Auth validates HTTP requests via token JWT:

	var jwt *auth.JWT
	// ...

	auth := middleware.NewAuth(jwt)
	auth.RequireTokenAuth(w, r, next)
	// ...

 */
package middleware
