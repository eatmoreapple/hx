// Package httpx provides HTTP request handling utilities and interfaces.
// It offers a comprehensive set of extractors for obtaining data from
// different parts of an HTTP request such as path parameters, headers,
// query parameters, form values, and cookies.
package httpx

import (
	"reflect"
	"sync"

	"github.com/eatmoreapple/hx/httpx/extractor"
)

// RequestExtractor is an alias for extractor.RequestExtractor interface,
// which defines methods for extracting data from HTTP requests.
type RequestExtractor = extractor.RequestExtractor

// RequestExtractorType holds the reflection Type of the RequestExtractor interface.
// This is used for runtime type checking and reflection-based operations
// when determining if a type implements the RequestExtractor interface.
var RequestExtractorType = reflect.TypeOf((*RequestExtractor)(nil)).Elem()

// implementsRequestExtractorTypeMap is a synchronized map that caches results
// of interface implementation checks to improve performance.
// Keys are reflect.Type objects, values are booleans indicating whether
// the type implements RequestExtractor.
var implementsRequestExtractorTypeMap = sync.Map{}

// isRequestExtractorType is an internal function that checks if a type implements
// the RequestExtractor interface. If the type is not a pointer, it converts it to
// a pointer type before checking.
func isRequestExtractorType(t reflect.Type) bool {
	if t.Kind() != reflect.Ptr {
		t = reflect.PointerTo(t)
	}
	return t.Implements(RequestExtractorType)
}

// IsRequestExtractorType checks if the given type implements the RequestExtractor interface.
// It uses a cache to avoid repeated checks for the same type, improving performance.
// If the type is not a pointer, it creates a pointer to the type before checking.
func IsRequestExtractorType(t reflect.Type) bool {
	if value, exists := implementsRequestExtractorTypeMap.Load(t); exists {
		return value.(bool)
	}

	result := isRequestExtractorType(t)
	implementsRequestExtractorTypeMap.Store(t, result)
	return result
}

// Type aliases for various extractor types.
// These provide convenient access to the underlying extractor implementations
// while maintaining the package's cohesive API.
type (
	// PathValueExtractor extracts values from URL path parameters
	PathValueExtractor[T extractor.Value] = extractor.PathValueExtractor[T]
	// FromPath is a shorthand for PathValueExtractor
	FromPath[T extractor.Value] = extractor.FromPath[T]

	// HeaderValueExtractor extracts values from HTTP headers
	HeaderValueExtractor[T extractor.Value] = extractor.HeaderValueExtractor[T]
	// FromHeader is a shorthand for HeaderValueExtractor
	FromHeader[T extractor.Value] = extractor.FromHeader[T]

	// QueryValueExtractor extracts values from URL query parameters
	QueryValueExtractor[T extractor.Value] = extractor.QueryValueExtractor[T]
	// FromQuery is a shorthand for QueryValueExtractor
	FromQuery[T extractor.Value] = extractor.FromQuery[T]

	// FormValueExtractor extracts values from form submissions
	FormValueExtractor[T extractor.Value] = extractor.FormValueExtractor[T]
	// FromForm is a shorthand for FormValueExtractor
	FromForm[T extractor.Value] = extractor.FromForm[T]

	// CookieValueExtractor extracts values from HTTP cookies
	CookieValueExtractor[T extractor.Value] = extractor.CookieValueExtractor[T]
	// FromCookie is a shorthand for CookieValueExtractor
	FromCookie[T extractor.Value] = extractor.FromCookie[T]
)

// Additional type aliases for complete extractors that handle
// collections of values rather than single named values.
type (
	// HeaderExtractor provides access to all HTTP headers in a request
	HeaderExtractor = extractor.HeaderExtractor
	// CookieExtractor provides access to all cookies in a request
	CookieExtractor = extractor.CookieExtractor
	// QueryExtractor provides access to all query parameters in a request
	QueryExtractor = extractor.QueryExtractor
	// FormExtractor provides access to all form values in a request
	FormExtractor = extractor.FormExtractor
)

// Empty is a no-op extractor that always succeeds without extracting any values.
// It can be used as a placeholder when an extractor is required but no extraction is needed.
type Empty = extractor.Empty
