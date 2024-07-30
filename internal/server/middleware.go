package server

import (
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func (r *Router) Use(mw Middleware) {
	for path, handler := range r.routes {
		r.routes[path] = mw(handler)
	}
}
