package twix

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const TwixContextKey contextKey = "twixContext"

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

// In the twix package (Router struct file)

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	method := req.Method

	for route, handlers := range r.routes {
		if match, params := matchRoute(route, path); match {
			if handler, ok := handlers[method]; ok {
				ctx := &Context{
					ResponseWriter: w,
					Request:        req,
					Params:         params,
				}
				req = req.WithContext(context.WithValue(req.Context(), TwixContextKey, ctx))
				handler.ServeHTTP(w, req)
				return
			}
		}
	}

	http.NotFound(w, req)
}

func applyMiddlewares(handler http.HandlerFunc, middlewares []func(http.Handler) http.Handler) http.HandlerFunc {
	h := http.HandlerFunc(handler)

	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h).(http.HandlerFunc)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value(TwixContextKey).(*Context)
		if ctx == nil {
			ctx = &Context{
				ResponseWriter: w,
				Request:        r,
				Params:         make(map[string]string),
			}
		}
		newReq := r.WithContext(context.WithValue(r.Context(), TwixContextKey, ctx))
		h.ServeHTTP(w, newReq)
	}
}

// matchRoute matches a route path with the request path and extracts parameters
func matchRoute(route, path string) (bool, map[string]string) {
	routeParts := strings.Split(route, "/")
	pathParts := strings.Split(path, "/")

	if len(routeParts) > len(pathParts) {
		return false, nil
	}

	params := make(map[string]string)
	for i, part := range routeParts {
		if i >= len(pathParts) {
			return false, nil
		}
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

func (g *Group) AddRoute(method, path string, handler http.HandlerFunc) {
	fullPath := g.prefix + path

	allMiddlewares := append(g.middlewares, g.router.middlewares...)
	finalHandler := applyMiddlewares(handler, allMiddlewares)

	g.router.AddRoute(method, fullPath, finalHandler)
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
