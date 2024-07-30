package main

import (
	"fmt"
	"net/http"

	"github.com/farhanmobashir/twix/internal/server"
)

func main() {
	router := server.NewRouter()
	router.AddRoute("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	fmt.Println("Server starting now")
	http.ListenAndServe(":8080", router)
}
