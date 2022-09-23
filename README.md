# helper-go
Helper for services written in Golang.

## Features
* Package auth provides authentication and authorization methods using JSON Web Tokens.
* Package cerror provides customized errors that represent status for HTTP responses and specific treatments.
* Package concurrent/executor provides an implementation of the Executor to perform background processing, limiting resource consumption when executing a collection of tasks.
* Package concurrent/flow provides a Flow implementation to perform background processing and notify your subscribers through an Emitter.
* Package concurrent/gopool provides an implementation of tools to reuse goroutine and limit resource consumption when running a collection of tasks.
* Package ebus implements the event bus design pattern, being an alternative to component communication while maintaining loose coupling and separation of interests principles.  
* Package httputil provides http utility methods.
* Package kafka is a simple wrapper for the kafka-go segmentio library, providing tools for consuming and producing events.
* Package middleware provides settings for CORS, GZIP and JWT token validation.
* Package observability/event implements routines to produce event log for a topic in Apache Kafka.
* Package observability/metric implements routines to collect metrics from localhost and send to a topic in Apache Kafka. The metrics collected are: "uptime", "os", "total memory", "memory used", "cpu count", "cpu user", "cpu system", "cpu idle" and "num goroutines".
* Package util/hashutil provides utility functions to generate and validate hash.
 