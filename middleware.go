package hx

import (
	"context"
	"net/http"
)

// Middleware represents a function that wraps a HandlerFunc and returns a new HandlerFunc.
// It can be used to add common functionality like logging, authentication, etc.
type Middleware func(HandlerFunc) HandlerFunc

// Chain combines multiple middleware into a single middleware.
// Middleware will be executed in the order they are passed.
// Example:
//
//	handler := Chain(
//	    LoggerMiddleware,
//	    AuthMiddleware,
//	    TimeoutMiddleware,
//	)(finalHandler)
func Chain(middlewares ...Middleware) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// WithValue is a middleware that injects a key-value pair into the request's context.
// The key must be a comparable type (e.g., string, int), and the value can be any type.
// This is useful for passing data (e.g., user information, request IDs) down the middleware chain.
func WithValue[T comparable](key T, value any) Middleware {
	return func(handlerFunc HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			r = r.WithContext(context.WithValue(r.Context(), key, value))
			return handlerFunc(w, r)
		}
	}
}
