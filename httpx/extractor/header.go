package extractor

import (
	"net/http"
	"reflect"
)

// HeaderValueExtractor implements RequestExtractor for HTTP header values.
// It extracts and stores header values of a specified string-like type T.
type HeaderValueExtractor[T Value] struct {
	baseValueExtractor[T]
}

// FromRequest implements RequestExtractor.FromRequest by extracting the header value
// using the resolved value name. The header value is converted to type T.
func (r *HeaderValueExtractor[T]) FromRequest(request *http.Request) error {
	return r.fromRequest(request, "")
}

// FromRequestField extracts a header value using the containing field's tags
// or Go name as a fallback value name.
func (r *HeaderValueExtractor[T]) FromRequestField(request *http.Request, field reflect.StructField) error {
	return r.fromRequest(request, valueNameFromField(field, "header"))
}

func (r *HeaderValueExtractor[T]) fromRequest(request *http.Request, fallbackName string) error {
	name, err := r.resolvedValueName(fallbackName)
	if err != nil {
		return err
	}
	r.value = T(request.Header.Get(name))
	return nil
}

type HeaderExtractor http.Header

func (r *HeaderExtractor) FromRequest(request *http.Request) error {
	*r = HeaderExtractor(request.Header)
	return nil
}
