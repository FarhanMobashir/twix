package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/farhanmobashir/twix"
)

// Middleware function
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	// First ServeMux for port 8080
	mux1 := http.NewServeMux()
	router1 := twix.New()

	// Add middleware using Use method
	router1.Use(loggingMiddleware)

	router1.Get("/:name", func(ctx *twix.Context) {
		fmt.Println("API Hit")
		name := ctx.Param("name")
		fmt.Println(name, "name here")
		str := "Hello  Wame " + name
		ctx.Status(http.StatusOK).String(str)
	})
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
