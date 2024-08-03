# Twix

Twix is a lightweight, modular router designed to simplify the process of building HTTP services in Go. With an intuitive API and support for middleware, Twix helps you create clean and maintainable web applications.

## Features

- Simple and intuitive routing
- Middleware support
- Routing groups for modular organization
- Context management for request data
- CORS, logging, rate limiting, and JWT authentication middleware

## Documentation

Comprehensive documentation for Twix can be found at the [Twix Documentation Site](https://twix-go.netlify.app/).

## Installation

To install Twix, run:

```bash
go get github.com/farhanmobashir/twix
```

## Quick Start

Here's a quick example to get you started with Twix:

```go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/farhanmobashir/twix"
	"github.com/farhanmobashir/twix/middlewares"
)

// Handler function for the route
func nameHandler(w http.ResponseWriter, r *http.Request) {
	ctx, ok := r.Context().Value(twix.TwixContextKey).(*twix.Context)
	if !ok {
		http.Error(w, "Invalid context", http.StatusInternalServerError)
		return
	}

	name := ctx.Param("name")
	if name == "" {
		http.Error(w, "Name parameter is missing", http.StatusBadRequest)
		return
	}
	str := "Hello, " + name
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(str))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func main() {
	router := twix.New()

	corsConfig := middlewares.CorsConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}

	rateLimitConfig := middlewares.RateLimitConfig{
		RequestLimit: 5,
		WindowSize:   time.Second * 15,
	}

	jwtConfig := middlewares.JWTConfig{
		SecretKey:   []byte("hello"),
		TokenSource: middlewares.Header,
		CookieName:  "jwt_token",
	}

	router.Use(middlewares.CorsMiddleware(corsConfig))
	router.Use(middlewares.RecoveryMiddleware)
	router.Use(middlewares.RateLimit(rateLimitConfig))
	router.Use(middlewares.LoggingMiddleware)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		panic("foo")
	})

	apiGroup := router.Group("/api")
	apiGroup.Use(middlewares.JWTAuth(jwtConfig))
	apiGroup.Get("/hello/:name", nameHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
```

## Middleware

Twix comes with several built-in middleware:

### CORS Middleware

Allows cross-origin requests.

```go
corsConfig := middlewares.CorsConfig{
	AllowedOrigins:   []string{"*"},
	AllowedMethods:   []string{"GET"},
	AllowedHeaders:   []string{"Content-Type", "Authorization"},
	AllowCredentials: true,
}
router.Use(middlewares.CorsMiddleware(corsConfig))
```

### Logging Middleware

Logs incoming requests with colored output.

```go
router.Use(middlewares.LoggingMiddleware)
```

### Rate Limiting Middleware

Limits the number of requests from a single IP address within a specified time window.

```go
rateLimitConfig := middlewares.RateLimitConfig{
	RequestLimit: 5,
	WindowSize:   time.Second * 15,
}
router.Use(middlewares.RateLimit(rateLimitConfig))
```

### JWT Authentication Middleware

Handles JWT authentication, supporting both header and cookie token sources.

```go
jwtConfig := middlewares.JWTConfig{
	SecretKey:   []byte("hello"),
	TokenSource: middlewares.Header, // or middlewares.Cookie
	CookieName:  "jwt_token",
}
router.Use(middlewares.JWTAuth(jwtConfig))
```

### Recovery Middleware

Recovers from panics and returns a 500 Internal Server Error.

```go
router.Use(middlewares.RecoveryMiddleware)
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request with any improvements or features you'd like to add.
