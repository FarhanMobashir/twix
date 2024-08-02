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
	log.Println("nameHandler started")
	log.Printf("Request context: %+v", r.Context())
	ctx, ok := r.Context().Value(twix.TwixContextKey).(*twix.Context)
	if !ok {
		log.Printf("Context type: %T", r.Context().Value(twix.TwixContextKey))
		log.Println("Invalid context")
		http.Error(w, "Invalid context", http.StatusInternalServerError)
		return
	}
	log.Printf("Twix context: %+v", ctx)
	log.Printf("Params: %v", ctx.Params)
	name := ctx.Params["name"]
	if name == "" {
		log.Println("Name parameter is missing")
		http.Error(w, "Name parameter is missing", http.StatusBadRequest)
		return
	}
	// Prepare the response
	str := "Hello, " + name
	// Set the content type and status code
	log.Println("Setting content type and status code")
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	// Write the response body
	_, err := w.Write([]byte(str))
	if err != nil {
		log.Println("Error writing response:", err)
		// Handle potential error when writing the response
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	log.Println("nameHandler finished")
}

func main() {
	// Create a new router
	router := twix.New()

	// Define CORS configuration
	corsConfig := middlewares.CorsConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}

	// Define Rate Limit Configuration
	rateLimitConfig := middlewares.RateLimitConfig{
		RequestLimit: 50,
		WindowSize:   time.Second * 15,
	}

	// Define JWT Configuration
	// jwtConfig := middlewares.JWTConfig{
	// 	SecretKey:   []byte("hello"),
	// 	TokenSource: middlewares.Header, // or middlewares.Cookie if you prefer
	// 	CookieName:  "jwt_token",
	// }

	// Add middleware using Use method
	router.Use(middlewares.CorsMiddleware(corsConfig))
	router.Use(middlewares.RecoveryMiddleware)
	router.Use(middlewares.RateLimit(rateLimitConfig))
	router.Use(middlewares.LoggingMiddleware)

	router.Get("/", func(http.ResponseWriter, *http.Request) { panic("foo") })

	// Create a routing group with the prefix /api
	apiGroup := router.Group("/api")

	// apiGroup.Use(middlewares.JWTAuth(jwtConfig))

	// Define routes within the /api group
	apiGroup.Get("/hello/:name", nameHandler)

	// Define the server
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("Starting server on :8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
