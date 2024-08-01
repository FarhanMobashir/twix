package main

import (
	"log"
	"net/http"
	"time"

	"github.com/farhanmobashir/twix"
	"github.com/farhanmobashir/twix/internal/middlewares"
)

// Handler function for the route
func nameHandler(w http.ResponseWriter, r *http.Request) {
	name := twix.URLParam(r, "name")
	if name == "" {
		http.Error(w, "Name parameter is missing", http.StatusBadRequest)
		return
	}
	// Prepare the response
	str := "Hello, " + name
	// Set the content type and status code
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	// Write the response body
	_, err := w.Write([]byte(str))
	if err != nil {
		// Handle potential error when writing the response
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
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
		RequestLimit: 5,
		WindowSize:   time.Second * 15,
	}

	// Add middleware using Use method
	router.Use(middlewares.CorsMiddleware(corsConfig))
	router.Use(middlewares.RateLimit(rateLimitConfig))
	router.Use(middlewares.LoggingMiddleware)

	// Create a routing group with the prefix /api
	apiGroup := router.Group("/api")

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
