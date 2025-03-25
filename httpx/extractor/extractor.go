package extractor

import "net/http"

// RequestExtractor defines the interface for types that can extract data from HTTP requests.
// Implementations should handle parsing and validating request data.
type RequestExtractor interface {
	FromRequest(*http.Request) error
}

type Empty struct{}

func (e *Empty) FromRequest(*http.Request) error { return nil }
