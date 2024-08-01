package twix

import (
	"context"
	"net/http"
	"strings"
)

// Router holds route definitions and middleware
type Router struct {
	routes      map[string]map[string]http.HandlerFunc
	middlewares []func(http.Handler) http.Handler
	groups      []*Group
}

// Group represents a routing group
type Group struct {
	prefix      string
	middlewares []func(http.Handler) http.Handler
	router      *Router
}

// New creates a new Router instance
func New() *Router {
	return &Router{
		routes:      make(map[string]map[string]http.HandlerFunc),
		middlewares: []func(http.Handler) http.Handler{},
	}
}

// Group creates a new routing group with a given prefix
func (r *Router) Group(prefix string) *Group {
	group := &Group{
		prefix:      prefix,
		router:      r,
		middlewares: []func(http.Handler) http.Handler{},
	}
	r.groups = append(r.groups, group)
	return group
}

// AddRoute adds a route handler for a specific method and path
func (r *Router) AddRoute(method, path string, handler http.HandlerFunc) {
	if r.routes[path] == nil {
		r.routes[path] = make(map[string]http.HandlerFunc)
	}
	r.routes[path][method] = handler
}

// Get adds a GET route handler
func (r *Router) Get(path string, handler http.HandlerFunc) {
	r.AddRoute("GET", path, handler)
}

// Post adds a POST route handler
func (r *Router) Post(path string, handler http.HandlerFunc) {
	r.AddRoute("POST", path, handler)
}

// Delete adds a DELETE route handler
func (r *Router) Delete(path string, handler http.HandlerFunc) {
	r.AddRoute("DELETE", path, handler)
}

// Patch adds a PATCH route handler
func (r *Router) Patch(path string, handler http.HandlerFunc) {
	r.AddRoute("PATCH", path, handler)
}

// Put adds a PUT route handler
func (r *Router) Put(path string, handler http.HandlerFunc) {
	r.AddRoute("PUT", path, handler)
}

// Use adds middleware to the router
func (r *Router) Use(middleware func(http.Handler) http.Handler) {
	r.middlewares = append(r.middlewares, middleware)
}

// ServeHTTP processes HTTP requests
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	method := req.Method

	for route, handlers := range r.routes {
		if match, params := matchRoute(route, path); match {
			ctx := &Context{
				ResponseWriter: w,
				Request:        req,
				Params:         params,
			}
			if handler, ok := handlers[method]; ok {
				finalHandler := applyMiddlewares(handler, r.middlewares, ctx)
				finalHandler.ServeHTTP(w, req)
				return
			}
		}
	}

	http.NotFound(w, req)
}

type contextKey string

const TwixContextKey contextKey = "twixContext"

// applyMiddlewares applies middleware functions to a handler
func applyMiddlewares(handler http.HandlerFunc, middlewares []func(http.Handler) http.Handler, ctx *Context) http.Handler {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	})

	// Apply middlewares in reverse order
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h).(http.HandlerFunc)
	}

	// Return a function that uses ServeHTTP on the resulting http.Handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Store the original Context in the request's context
		newReq := r.WithContext(context.WithValue(r.Context(), TwixContextKey, ctx))
		h.ServeHTTP(w, newReq)
	})
}

// matchRoute matches a route path with the request path and extracts parameters
func matchRoute(route, path string) (bool, map[string]string) {
	routeParts := strings.Split(route, "/")
	pathParts := strings.Split(path, "/")

	if len(routeParts) != len(pathParts) {
		return false, nil
	}

	params := make(map[string]string)
	for i, part := range routeParts {
		if strings.HasPrefix(part, ":") {
			params[part[1:]] = pathParts[i]
		} else if part != pathParts[i] {
			return false, nil
		}
	}

	return true, params
}

func URLParam(r *http.Request, param string) string {
	ctx, ok := r.Context().Value(TwixContextKey).(*Context)
	if !ok {
		return ""
	}
	return ctx.Params[param]
}

// Group functions

// AddRoute adds a route handler for a specific method and path within the group
func (g *Group) AddRoute(method, path string, handler http.HandlerFunc) {
	fullPath := g.prefix + path
	g.router.AddRoute(method, fullPath, handler)
}

// Get adds a GET route handler within the group
func (g *Group) Get(path string, handler http.HandlerFunc) {
	g.AddRoute("GET", path, handler)
}

// Post adds a POST route handler within the group
func (g *Group) Post(path string, handler http.HandlerFunc) {
	g.AddRoute("POST", path, handler)
}

// Delete adds a DELETE route handler within the group
func (g *Group) Delete(path string, handler http.HandlerFunc) {
	g.AddRoute("DELETE", path, handler)
}

// Patch adds a PATCH route handler within the group
func (g *Group) Patch(path string, handler http.HandlerFunc) {
	g.AddRoute("PATCH", path, handler)
}

// Put adds a PUT route handler within the group
func (g *Group) Put(path string, handler http.HandlerFunc) {
	g.AddRoute("PUT", path, handler)
}

// Use adds middleware to the group
func (g *Group) Use(middleware func(http.Handler) http.Handler) {
	g.middlewares = append(g.middlewares, middleware)
}
