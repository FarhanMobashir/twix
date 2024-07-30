package twix

import (
	"context"
	"net/http"
	"strings"
)

// Router holds route definitions and middleware
type Router struct {
	routes      map[string]map[string]func(*Context)
	middlewares []func(http.Handler) http.Handler
}

// New creates a new Router instance
func New() *Router {
	return &Router{
		routes:      make(map[string]map[string]func(*Context)),
		middlewares: []func(http.Handler) http.Handler{},
	}
}

// AddRoute adds a route handler for a specific method and path
func (r *Router) AddRoute(method, path string, handler func(*Context)) {
	if r.routes[path] == nil {
		r.routes[path] = make(map[string]func(*Context))
	}
	r.routes[path][method] = handler
}

// Get adds a GET route handler
func (r *Router) Get(path string, handler func(*Context)) {
	r.AddRoute("GET", path, handler)
}

// Post adds a POST route handler
func (r *Router) Post(path string, handler func(*Context)) {
	r.AddRoute("POST", path, handler)
}

// Delete adds a DELETE route handler
func (r *Router) Delete(path string, handler func(*Context)) {
	r.AddRoute("DELETE", path, handler)
}

// Patch adds a PATCH route handler
func (r *Router) Patch(path string, handler func(*Context)) {
	r.AddRoute("PATCH", path, handler)
}

// Put adds a PUT route handler
func (r *Router) Put(path string, handler func(*Context)) {
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
				finalHandler := applyMiddlewares(handler, r.middlewares)
				finalHandler(ctx)
				return
			}
		}
	}

	http.NotFound(w, req)
}

// applyMiddlewares applies middleware functions to a handler
func applyMiddlewares(handler func(*Context), middlewares []func(http.Handler) http.Handler) func(*Context) {
	// Convert the handler function into an http.Handler
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// We'll use a context key to pass the original Context through the middleware chain
		ctx := r.Context()
		origCtx, ok := ctx.Value("twixContext").(*Context)
		if !ok {
			// This shouldn't happen, but just in case
			origCtx = &Context{
				ResponseWriter: w,
				Request:        r,
			}
		}
		handler(origCtx)
	})

	// Apply middlewares in reverse order
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h).(http.HandlerFunc)
	}

	// Return a function that uses ServeHTTP on the resulting http.Handler
	return func(ctx *Context) {
		// Store the original Context in the request's context
		newReq := ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "twixContext", ctx))
		h.ServeHTTP(ctx.ResponseWriter, newReq)
	}
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
