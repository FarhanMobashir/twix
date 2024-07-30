package server

import (
	"net/http"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Default Handler"))
}
