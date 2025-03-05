// Package hx provides a lightweight and type-safe HTTP handler framework with generic support.
package hx

import (
	"fmt"
	"net/http"
	"path"
	"strings"
)

// Router is the main router structure that handles HTTP request routing and error handling.
// It wraps the standard http.ServeMux and adds additional functionality like middleware support,
// route groups, and custom error handling.
type Router struct {
	// ErrHandler handles any errors returned by handlers
	ErrHandler ErrorHandler

	// mux is the underlying HTTP request multiplexer
	mux *http.ServeMux

	// basePath is the base path for all routes in this router
	basePath string

	// middleware stack for this router
	middleware []Middleware
}

// RouterOption defines a function type for configuring a Router instance.
type RouterOption func(*Router)

// WithErrorHandler sets a custom error handler for the router.
func WithErrorHandler(handler ErrorHandler) RouterOption {
	return func(r *Router) {
		r.ErrHandler = handler
	}
}

// WithMiddleware adds middleware to the router.
func WithMiddleware(middleware ...Middleware) RouterOption {
	return func(r *Router) {
		r.middleware = append(r.middleware, middleware...)
	}
}

// New creates a new Router instance with the given options.
// If no error handler is provided, it uses a default one that returns 500 Internal Server Error.
func New(options ...RouterOption) *Router {
	r := &Router{
		mux:      http.NewServeMux(),
		basePath: "/",
		ErrHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}

	for _, opt := range options {
		opt(r)
	}

	return r
}

// Group creates a new router group with the given path prefix.
// All routes registered on the group will be prefixed with the group's path.
// The group inherits the middleware stack from its parent.
func (r *Router) Group(prefix string) *Router {
	return &Router{
		mux:        r.mux,
		basePath:   path.Join(r.basePath, prefix),
		ErrHandler: r.ErrHandler,
		middleware: append([]Middleware{}, r.middleware...),
	}
}

// Use adds middleware to the router's middleware stack.
// Middleware will be executed in the order they are added.
func (r *Router) Use(middleware ...Middleware) {
	r.middleware = append(r.middleware, middleware...)
}

// Handle registers a new route with the given method and path.
// The handler will be wrapped with the router's middleware stack.
func (r *Router) Handle(method, path string, handler HandlerFunc) {
	// Ensure path starts with /
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// Combine base path with route path
	fullPath := joinPath(r.basePath, path)
	pattern := fmt.Sprintf("%s %s", method, fullPath)

	// Apply middleware stack
	if len(r.middleware) > 0 {
		handler = Chain(r.middleware...)(handler)
	}

	// Register the route
	r.mux.HandleFunc(pattern, func(w http.ResponseWriter, req *http.Request) {
		if err := handler(w, req); err != nil {
			r.ErrHandler(w, req, err)
		}
	})
}

// Common HTTP method handlers
// These methods provide a convenient way to register routes for specific HTTP methods.

// GET registers a new GET route.
func (r *Router) GET(path string, handler HandlerFunc) {
	r.Handle(http.MethodGet, path, handler)
}

// POST registers a new POST route.
func (r *Router) POST(path string, handler HandlerFunc) {
	r.Handle(http.MethodPost, path, handler)
}

// PUT registers a new PUT route.
func (r *Router) PUT(path string, handler HandlerFunc) {
	r.Handle(http.MethodPut, path, handler)
}

// DELETE registers a new DELETE route.
func (r *Router) DELETE(path string, handler HandlerFunc) {
	r.Handle(http.MethodDelete, path, handler)
}

// PATCH registers a new PATCH route.
func (r *Router) PATCH(path string, handler HandlerFunc) {
	r.Handle(http.MethodPatch, path, handler)
}

// OPTIONS registers a new OPTIONS route.
func (r *Router) OPTIONS(path string, handler HandlerFunc) {
	r.Handle(http.MethodOptions, path, handler)
}

// HEAD registers a new HEAD route.
func (r *Router) HEAD(path string, handler HandlerFunc) {
	r.Handle(http.MethodHead, path, handler)
}

// ServeHTTP implements the http.Handler interface.
// This method is called by the HTTP server to handle incoming requests.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// joinPath joins two path segments ensuring there is exactly one slash between them.
func joinPath(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	default:
		return a + b
	}
}
