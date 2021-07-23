# helper-go
Helper for services written in Golang.

## Features
* Package auth provides authentication and authorization methods using JSON Web Tokens.
* Package cerror provides customized errors that represent status for HTTP responses and specific treatments.
* Package concurrent/executor provides an implementation of the Executor to perform background processing, limiting resource consumption when executing a collection of jobs.
* Package concurrent/flow provides a Flow implementation to perform background processing and notify your subscribers through an Emitter.
* Package ebus implements the event bus design pattern, being an alternative to component communication while maintaining loose coupling and separation of interests principles.  
* Package httputil provides http utility methods.
* Package kafka is a simple wrapper for the kafka-go segmentio library, providing tools for consuming and producing events.
* Package middleware provides settings for CORS, GZIP and JWT token validation.
* Package util/hashutil provides utility functions to generate and validate hash.
 