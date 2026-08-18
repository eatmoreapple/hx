package extractor

import (
	"net/http"
	"reflect"
)

// PathValueExtractor implements RequestExtractor for path parameters.
// It extracts named path values from HTTP requests using Go 1.22's Value feature.
type PathValueExtractor[T Value] struct {
	baseValueExtractor[T]
}

// FromRequest implements RequestExtractor.FromRequest by extracting the path value
func (r *PathValueExtractor[T]) FromRequest(request *http.Request) error {
	return r.fromRequest(request, "")
}

// FromRequestField extracts a path value using the containing field's tags or
// Go name as a fallback value name.
func (r *PathValueExtractor[T]) FromRequestField(request *http.Request, field reflect.StructField) error {
	return r.fromRequest(request, valueNameFromField(field, "path"))
}

func (r *PathValueExtractor[T]) fromRequest(request *http.Request, fallbackName string) error {
	name, err := r.resolvedValueName(fallbackName)
	if err != nil {
		return err
	}
	r.value = T(request.PathValue(name))
	return nil
}
