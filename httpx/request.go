// Package httpx provides HTTP request handling utilities and interfaces.
package httpx

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"sync"
)

// RequestExtractor defines the interface for types that can extract data from HTTP requests.
// Implementations should handle parsing and validating request data.
type RequestExtractor interface {
	FromRequest(*http.Request) error
}

// RequestExtractorType holds the reflection Type of the RequestExtractor interface.
// This is useful for runtime type checking and reflection-based operations.
var RequestExtractorType = reflect.TypeOf((*RequestExtractor)(nil)).Elem()

// implementsRequestExtractorTypeMap is a synchronized map used to cache the results of whether a given type implements the RequestExtractor interface.
// The key is the reflect.Type, and the value is a boolean indicating whether the type implements the interface.
var implementsRequestExtractorTypeMap = sync.Map{}

func isRequestExtractorType(t reflect.Type) bool {
	if t.Kind() != reflect.Ptr {
		t = reflect.PointerTo(t)
	}
	return t.Implements(RequestExtractorType)
}

// IsRequestExtractorType checks if the given type implements the RequestExtractor interface.
// If the type is not a pointer, it creates a pointer to the type and then checks
// if the resulting type implements the RequestExtractor interface.
func IsRequestExtractorType(t reflect.Type) bool {
	if value, exists := implementsRequestExtractorTypeMap.Load(t); exists {
		return value.(bool)
	}

	result := isRequestExtractorType(t)
	implementsRequestExtractorTypeMap.Store(t, result)
	return result
}

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

// Int8 converts the value to int8.
// Returns an error if the value cannot be parsed as an 8-bit integer.
func (b baseValueExtractor[T]) Int8() (int8, error) {
	v, err := strconv.ParseInt(string(b.value), 10, 8)
	return int8(v), err
}

// Int16 converts the value to int16.
// Returns an error if the value cannot be parsed as a 16-bit integer.
func (b baseValueExtractor[T]) Int16() (int16, error) {
	v, err := strconv.ParseInt(string(b.value), 10, 16)
	return int16(v), err
}

// Int32 converts the value to int32.
// Returns an error if the value cannot be parsed as an integer.
func (b baseValueExtractor[T]) Int32() (int32, error) {
	v, err := strconv.ParseInt(string(b.value), 10, 32)
	return int32(v), err
}

// Int64 converts the value to int64.
// Returns an error if the value cannot be parsed as an integer.
func (b baseValueExtractor[T]) Int64() (int64, error) {
	return strconv.ParseInt(string(b.value), 10, 64)
}

// Int converts the value to int.
// Returns an error if the value cannot be parsed as an integer.
func (b baseValueExtractor[T]) Int() (int, error) {
	v, err := strconv.ParseInt(string(b.value), 10, 0)
	return int(v), err
}

// Uint8 converts the value to uint8.
// Returns an error if the value cannot be parsed as an 8-bit unsigned integer.
func (b baseValueExtractor[T]) Uint8() (uint8, error) {
	v, err := strconv.ParseUint(string(b.value), 10, 8)
	return uint8(v), err
}

// Uint16 converts the value to uint16.
// Returns an error if the value cannot be parsed as a 16-bit unsigned integer.
func (b baseValueExtractor[T]) Uint16() (uint16, error) {
	v, err := strconv.ParseUint(string(b.value), 10, 16)
	return uint16(v), err
}

// Uint32 converts the value to uint32.
// Returns an error if the value cannot be parsed as an unsigned integer.
func (b baseValueExtractor[T]) Uint32() (uint32, error) {
	v, err := strconv.ParseUint(string(b.value), 10, 32)
	return uint32(v), err
}

// Uint64 converts the value to uint64.
// Returns an error if the value cannot be parsed as an unsigned integer.
func (b baseValueExtractor[T]) Uint64() (uint64, error) {
	return strconv.ParseUint(string(b.value), 10, 64)
}

// Uint converts the value to uint.
// Returns an error if the value cannot be parsed as an unsigned integer.
func (b baseValueExtractor[T]) Uint() (uint, error) {
	v, err := strconv.ParseUint(string(b.value), 10, 0)
	return uint(v), err
}

// Float64 converts the value to float64.
// Returns an error if the value cannot be parsed as a floating-point number.
func (b baseValueExtractor[T]) Float64() (float64, error) {
	return strconv.ParseFloat(string(b.value), 64)
}

// Float32 converts the value to float32.
// Returns an error if the value cannot be parsed as a floating-point number.
func (b baseValueExtractor[T]) Float32() (float32, error) {
	v, err := strconv.ParseFloat(string(b.value), 32)
	return float32(v), err
}

// Bool converts the value to bool.
// Returns an error if the value cannot be parsed as a boolean.
// Accepts 1, t, T, TRUE, true for true and 0, f, F, FALSE, false for false.
func (b baseValueExtractor[T]) Bool() (bool, error) {
	return strconv.ParseBool(string(b.value))
}

// String returns the value as a string.
// This is a convenience method that simply converts the value to string.
func (b baseValueExtractor[T]) String() string {
	return string(b.value)
}

// FromRequest is a placeholder implementation that should be overridden by embedding types.
// It will panic if called directly.
func (b *baseValueExtractor[T]) FromRequest(*http.Request) error {
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

// Empty is a no-op implementation of RequestExtractor.
type Empty struct{}

func (e Empty) FromRequest(*http.Request) error { return nil }
