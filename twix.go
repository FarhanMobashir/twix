package twix

import (
	"net/http"

	"github.com/farhanmobashir/twix/internal/server"
)

func CreateRouter() *server.Router {
	return server.NewRouter()
}

func AddRoute(r *server.Router, path string, handler http.HandlerFunc) {
	r.AddRoute(path, handler)
}
