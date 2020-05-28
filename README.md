# helper-go
Helper for services written in Golang.

## Features
* Package auth provides authentication and authorization methods using JSON Web Tokens.
* Package cerror provides customized errors that represent status for HTTP responses and specific treatments.
* Package concurrent provides a Flow implementation to perform background processing and notify your subscribers through an Emitter.
* Package controller implements a Controller with utility methods.
* Package integration provides a custom Kafka wrapper to work with the outbox pattern.
* Package middleware provides settings for CORS, GZIP and JWT token validation.
* Package hashutil provides utility functions to generate and validate hash.
