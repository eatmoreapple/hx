package extractor

import (
	"net/http"
	"reflect"
)

// RequestExtractor defines the interface for types that can extract data from HTTP requests.
// Implementations should handle parsing and validating request data.
type RequestExtractor interface {
	FromRequest(*http.Request) error
}

// FieldRequestExtractor extends RequestExtractor with struct-field context.
// Struct binders use FromRequestField in place of FromRequest when this
// interface is implemented.
type FieldRequestExtractor interface {
	RequestExtractor
	FromRequestField(*http.Request, reflect.StructField) error
}

type Empty struct{}

func (e *Empty) FromRequest(*http.Request) error { return nil }
