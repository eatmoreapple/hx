// Package httpx provides HTTP request handling utilities and interfaces.
package httpx

import (
	"encoding/json"
	"net/http"
	"reflect"
)

// RequestExtractor defines the interface for types that can extract data from HTTP requests.
// Implementations should handle parsing and validating request data.
type RequestExtractor interface {
	FromRequest(*http.Request) error
}

// RequestExtractorType holds the reflection Type of the RequestExtractor interface.
// This is useful for runtime type checking and reflection-based operations.
var RequestExtractorType = reflect.TypeOf((*RequestExtractor)(nil)).Elem()

// Value is an interface for types that can be used as path parameters.
// It combines the PathValueName method with the constraint of being a string type.
type Value interface {
	// ValueName returns the name of the path parameter as defined in the route.
	ValueName() string
	~string
}

// baseValueExtractor provides common functionality for value extractors.
// It implements basic operations like value retrieval and JSON marshaling.
type baseValueExtractor[T Value] struct {
	value T // The extracted value after processing
}

// Value returns the extracted value.
// This method should be called after FromRequest has been executed successfully.
func (b baseValueExtractor[T]) Value() T {
	return b.value
}

// MarshalJSON implements json.Marshaler interface to provide JSON serialization
// of the extracted value.
func (b baseValueExtractor[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.value)
}

// FromRequest is a placeholder implementation that should be overridden by embedding types.
// It will panic if called directly.
func (b *baseValueExtractor[T]) FromRequest(request *http.Request) error {
	panic("not implemented")
}

// PathValueExtractor implements RequestExtractor for path parameters.
// It extracts named path values from HTTP requests using Go 1.22's Value feature.
type PathValueExtractor[T Value] struct {
	baseValueExtractor[T]
}

// FromRequest implements RequestExtractor.FromRequest by extracting the path value
// from the request using the name provided by ValueName().
func (r *PathValueExtractor[T]) FromRequest(request *http.Request) error {
	r.value = T(request.PathValue(r.value.ValueName()))
	return nil
}

// FromPath is a type alias for PathValueExtractor providing a shorter name
// while maintaining all functionality.
type FromPath[T Value] = PathValueExtractor[T]

// HeaderValueExtractor implements RequestExtractor for HTTP header values.
// It extracts and stores header values of a specified type T that implements the Value interface.
type HeaderValueExtractor[T Value] struct {
	baseValueExtractor[T]
}

// FromRequest implements RequestExtractor.FromRequest by extracting the header value
// using the name provided by ValueName(). The header value is converted to type T.
func (r *HeaderValueExtractor[T]) FromRequest(request *http.Request) error {
	r.value = T(request.Header.Get(r.value.ValueName()))
	return nil
}

// FromHeader is a type alias for HeaderValueExtractor providing a shorter name
// while maintaining all functionality.
type FromHeader[T Value] = HeaderValueExtractor[T]

// QueryValueExtractor implements RequestExtractor for query parameters.
// It extracts and stores query values of a specified type T that implements the Value interface.
type QueryValueExtractor[T Value] struct {
	baseValueExtractor[T]
}

// FromRequest implements RequestExtractor.FromRequest by extracting the query value
// using the name provided by ValueName(). The query value is converted to type T.
func (r *QueryValueExtractor[T]) FromRequest(request *http.Request) error {
	r.value = T(request.URL.Query().Get(r.value.ValueName()))
	return nil
}

// FromQuery is a type alias for QueryValueExtractor providing a shorter name
// while maintaining all functionality.
type FromQuery[T Value] = QueryValueExtractor[T]

// FormValueExtractor implements RequestExtractor for form values.
// It extracts and stores form values of a specified type T that implements the Value interface.
type FormValueExtractor[T Value] struct {
	baseValueExtractor[T]
}

// FromRequest implements RequestExtractor.FromRequest by extracting the form value
// using the name provided by ValueName(). The form value is converted to type T.
func (r *FormValueExtractor[T]) FromRequest(request *http.Request) error {
	r.value = T(request.FormValue(r.value.ValueName()))
	return nil
}

// FromForm is a type alias for FormValueExtractor providing a shorter name
// while maintaining all functionality.
type FromForm[T Value] = FormValueExtractor[T]

// CookieValueExtractor implements RequestExtractor for cookie values.
// It extracts and stores cookie values of a specified type T that implements the Value interface.
type CookieValueExtractor[T Value] struct {
	baseValueExtractor[T]
}

// FromRequest implements RequestExtractor.FromRequest by extracting the cookie value
// using the name provided by ValueName(). The cookie value is converted to type T.
func (r *CookieValueExtractor[T]) FromRequest(request *http.Request) error {
	cookie, err := request.Cookie(r.value.ValueName())
	if err != nil {
		return err
	}
	r.value = T(cookie.Value)
	return nil
}

// FromCookie is a type alias for CookieValueExtractor providing a shorter name
// while maintaining all functionality.
type FromCookie[T Value] = CookieValueExtractor[T]
