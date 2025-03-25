package extractor

import "net/http"

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
