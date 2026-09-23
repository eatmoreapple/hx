// Package hx provides a lightweight and type-safe HTTP handler framework with generic support.
// It focuses on request handling, response rendering, and middleware composition.
package hx

import (
	"context"
	"net/http"
	"reflect"
	"unsafe"

	"github.com/eatmoreapple/hx/binding"
	"github.com/eatmoreapple/hx/httpx"
)

// ErrorHandler is a function type that handles errors occurred during request processing.
// It receives the ResponseWriter to write the error response, the original Request that caused the error,
// and the error itself. This allows for custom error handling and formatting across the application.
type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

// HandlerFunc is the standard handler type for processing HTTP requests in evo.
// It follows a similar pattern to http.HandlerFunc but returns an error instead of void.
// This allows for better error handling and middleware composition.
// If an error is returned, it will be passed to the ErrorHandler for processing.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// Warp is a convenience function that wraps an http.HandlerFunc into a HandlerFunc.
func Warp(h http.HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		h(w, r)
		return nil
	}
}

// Generic creates a type-safe handler with specified Request and Response types.
// It's a type assertion function that ensures the handler conforms to the TypedHandlerFunc interface.
// This function is particularly useful when you want to explicitly declare the types of your handler
// or when type inference isn't sufficient.
//
// Example:
//
//	handler := Generic[UserRequest, UserResponse](func(ctx context.Context, req UserRequest) (UserResponse, error) {
//	    return UserResponse{Name: req.Name}, nil
//	})
func Generic[Request, Response any](h TypedHandlerFunc[Request, Response]) TypedHandlerFunc[Request, Response] {
	return h
}

// G is a shortcut for Generic
func G[Request, Response any](h TypedHandlerFunc[Request, Response]) TypedHandlerFunc[Request, Response] {
	return Generic(h)
}

// JSON converts a typed handler into a HandlerFunc that serializes its response as JSON.
// Request and Response are inferred from h.
func JSON[Request, Response any](h TypedHandlerFunc[Request, Response]) HandlerFunc {
	return h.JSON()
}

// Render is a generic handler function that processes requests of type Request
// and returns responses of type httpx.ResponseRender. It operates within a context and may return an error.
//
// Example:
//
//	handler := Renderer[UserRequest](func(ctx context.Context, req UserRequest) (httpx.ResponseRender, error) {
//	    return httpx.JSONResponse{Data: req}, nil
//	})
func Render[Request any](h TypedHandlerFunc[Request, httpx.ResponseRender]) HandlerFunc {
	var handler requestHandler[Request] = func(ctx context.Context, req Request) (httpx.ResponseRender, error) {
		responseRender, err := h(ctx, req)
		if err != nil {
			return nil, err
		}
		return responseRender, nil
	}
	return handler.asHandlerFunc()
}

// R is a shortcut for Renderer function.
// It provides the same functionality as Renderer but with a more concise name.
func R[Request any](h TypedHandlerFunc[Request, httpx.ResponseRender]) HandlerFunc {
	return Render(h)
}

// E is a convenience function for creating a handler that doesn't require any request data.
// It takes a function that accepts a context and returns a response of type Response.
func E[Response any](h func(ctx context.Context) (Response, error)) TypedHandlerFunc[Empty, Response] {
	return func(ctx context.Context, req Empty) (Response, error) {
		return h(ctx)
	}
}

// TypedHandlerFunc is a generic handler function that processes requests of type Request
// and returns responses of type Response. It operates within a context and may return an error.
type TypedHandlerFunc[Request, Response any] func(context.Context, Request) (Response, error)

// JSON converts the handler into a JSON response handler.
// The response will be automatically serialized to JSON format.
func (h TypedHandlerFunc[Request, Response]) JSON() HandlerFunc {
	var handler requestHandler[Request] = func(ctx context.Context, req Request) (httpx.ResponseRender, error) {
		resp, err := h(ctx, req)
		if err != nil {
			return nil, err
		}
		return httpx.JSONResponse{Data: resp}, nil
	}
	return handler.asHandlerFunc()
}

// String converts the handler into a string response handler.
// This method panics if the Response type is not string.
func (h TypedHandlerFunc[Request, Response]) String() HandlerFunc {
	if _, ok := any((*Response)(nil)).(*string); !ok {
		panic("String() only supports string response type")
	}
	var handler requestHandler[Request] = func(ctx context.Context, req Request) (httpx.ResponseRender, error) {
		resp, err := h(ctx, req)
		if err != nil {
			return nil, err
		}
		str := *(*string)(unsafe.Pointer(&resp))
		return httpx.StringResponse{Data: str}, nil
	}
	return handler.asHandlerFunc()
}

// XML converts the handler into an XML response handler.
// The response will be automatically serialized to XML format.
func (h TypedHandlerFunc[Request, Response]) XML() HandlerFunc {
	var handler requestHandler[Request] = func(ctx context.Context, req Request) (httpx.ResponseRender, error) {
		resp, err := h(ctx, req)
		if err != nil {
			return nil, err
		}
		return httpx.XMLResponse{Data: resp}, nil
	}
	return handler.asHandlerFunc()
}

// Pipe composes the handler with a series of middleware functions.
// Each middleware function has the signature func(ctx context.Context, req Request) error
// and can perform operations such as logging, authentication, or request modification.
// The middleware functions are executed in the order they are provided before the main handler is called.
//
// Note: When chaining multiple Pipe calls, the execution order is reversed.
// For example, with hx.G(app).Pipe(app2).Pipe(app3), app3 will execute first,
// followed by app2, and finally app (the original handler).
func (h TypedHandlerFunc[Request, Response]) Pipe(steps ...func(ctx context.Context, req Request) error) TypedHandlerFunc[Request, Response] {
	if len(steps) == 0 {
		return h
	}
	return func(ctx context.Context, request Request) (resp Response, err error) {
		// Execute middleware functions in order
		for _, middleware := range steps {
			if err := middleware(ctx, request); err != nil {
				return resp, err
			}
		}

		// Execute the final handler
		return h(ctx, request)
	}
}

// requestHandler is an internal type that handles the processing of requests
// and produces a ResponseRender for rendering the response.
type requestHandler[Request any] func(context.Context, Request) (httpx.ResponseRender, error)

// call executes the handler with the given request and writes the response.
func (h requestHandler[Request]) call(w http.ResponseWriter, r *http.Request, req Request) error {
	resp, err := h(r.Context(), req)
	if err != nil {
		return err
	}
	return resp.IntoResponse(w)
}

// asHandlerFunc converts the requestHandler into a standard HandlerFunc.
// RequestExtractor types are extracted from the request. Structs and pointers
// to structs are bound. Any other request type panics here, when the handler
// is built, instead of failing on the first request.
func (h requestHandler[Request]) asHandlerFunc() HandlerFunc {
	requestType := reflect.TypeFor[Request]()
	if httpx.IsRequestExtractorType(requestType) {
		return h.extractAndHandle()
	}
	if !isStructType(requestType) {
		panic("hx: request type " + requestType.String() + " must be a struct, a pointer to a struct, or implement FromRequest")
	}
	return h.bindAndHandle()
}

// isStructType reports whether t is a struct or a pointer to a struct.
func isStructType(t reflect.Type) bool {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t.Kind() == reflect.Struct
}

// createHandler builds a HandlerFunc that allocates a zero Request for each
// call, fills it through extractFunc, and passes that same value to the handler.
// Pointer and value requests are split once, while the handler is built.
// A pointer Request is already the address extractFunc writes through.
// A value Request is allocated with new so extractFunc can write through *Request;
// the handler receives the filled value after extraction.
func (h requestHandler[Request]) createHandler(extractFunc func(any, *http.Request) error) HandlerFunc {
	requestType := reflect.TypeFor[Request]()
	isPointer := requestType.Kind() == reflect.Pointer

	if isPointer {
		// Request is *T. reflect.New allocates T and returns *T.
		// Passing &req would bind **T, which the struct binder rejects.
		elemType := requestType.Elem()

		return func(w http.ResponseWriter, r *http.Request) error {
			req, _ := reflect.TypeAssert[Request](reflect.New(elemType))
			if err := extractFunc(req, r); err != nil {
				return err
			}
			return h.call(w, r, req)
		}
	}
	// Request is T. Bind *T, then pass the filled value.
	return func(w http.ResponseWriter, r *http.Request) error {
		req := new(Request)
		if err := extractFunc(req, r); err != nil {
			return err
		}
		return h.call(w, r, *req)
	}
}

// extractAndHandle creates a HandlerFunc that extracts request data using the RequestExtractor interface.
func (h requestHandler[Request]) extractAndHandle() HandlerFunc {
	requestType := reflect.TypeFor[Request]()
	if requestType.Kind() == reflect.Pointer {
		requestType = requestType.Elem()
	}
	field := requestType.Name()
	return h.createHandler(func(target any, r *http.Request) error {
		return httpx.WrapExtractError(field, target.(httpx.RequestExtractor).FromRequest(r))
	})
}

// bindAndHandle creates a HandlerFunc that binds request data using the ShouldBind function.
func (h requestHandler[Request]) bindAndHandle() HandlerFunc {
	return h.createHandler(func(target any, r *http.Request) error {
		return ShouldBind(r, target)
	})
}

// ShouldBind binds the request data to the given interface.
// It first tries to bind using the default binder based on Content-Type,
// then attempts to bind using the GenericBinder if the type implements RequestExtractor.
func ShouldBind(r *http.Request, e any) error {
	binder := binding.Default(r.Method, r.Header.Get("Content-Type"))
	if err := binder.Bind(r, e); err != nil {
		return err
	}
	// if each field has implemented RequestExtractor
	return binding.Generic().Bind(r, e)
}
