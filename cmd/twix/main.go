package main

import (
	"log"
	"net/http"

	"github.com/farhanmobashir/twix"
	"github.com/farhanmobashir/twix/internal/utils/middlewares"
)

// Handler function for the route
func nameHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the name parameter from the request context
	name := twix.URLParam(r, "name")
	if name == "" {
		// If the name parameter is missing, return a bad request error
		http.Error(w, "Name parameter is missing", http.StatusBadRequest)
		return
	}

	// Prepare the response
	str := "Hello Wame " + name

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
	// First ServeMux for port 8080
	mux1 := http.NewServeMux()
	router1 := twix.New()

	// Add middleware using Use method
	router1.Use(middlewares.LoggingMiddleware)

	// Define the route with the updated handler function
	router1.Get("/:name", nameHandler)
	mux1.Handle("/", router1)

	// Define the servers
	server1 := &http.Server{
		Addr:    ":8080",
		Handler: mux1,
	}

	log.Println("Starting server on :8080")
	err := server1.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
