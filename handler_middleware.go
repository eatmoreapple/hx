package hx

import (
	"context"
	"slices"
)

// HandlerMiddleware runs around a typed handler, after the request is bound.
// Call next with the context and request that the rest of the chain should see.
// Skip next to stop the chain.
type HandlerMiddleware[Request, Response any] func(ctx context.Context, req Request, next TypedHandlerFunc[Request, Response]) (Response, error)

// HandlerChain is handler middleware executed from left to right.
// The first middleware is the outermost call: it runs first, and code after
// next runs last.
type HandlerChain[Request, Response any] []HandlerMiddleware[Request, Response]

// Pipe returns a chain that runs middlewares from left to right.
// The returned chain is independent of the caller's slice.
func Pipe[Request, Response any](middlewares ...HandlerMiddleware[Request, Response]) HandlerChain[Request, Response] {
	if len(middlewares) == 0 {
		return nil
	}
	chain := make(HandlerChain[Request, Response], len(middlewares))
	copy(chain, middlewares)
	return chain
}

// Handle ends the chain with h.
// Middleware runs first, then h. An empty chain returns h unchanged.
func (c HandlerChain[Request, Response]) Handle(h TypedHandlerFunc[Request, Response]) TypedHandlerFunc[Request, Response] {
	if len(c) == 0 {
		return h
	}
	next := h
	for _, middleware := range slices.Backward(c) {
		inner := next
		next = func(ctx context.Context, req Request) (Response, error) {
			return middleware(ctx, req, inner)
		}
	}
	return next
}
