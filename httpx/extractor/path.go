package extractor

import "net/http"

// PathValueExtractor implements RequestExtractor for path parameters.
// It extracts named path values from HTTP requests using Go 1.22's Value feature.
type PathValueExtractor[T Value] struct {
	baseValueExtractor[T]
}

// FromRequest implements RequestExtractor.FromRequest by extracting the path value
// using the name provided by ValueName or the containing field's hx tag.
func (r *PathValueExtractor[T]) FromRequest(request *http.Request) error {
	name, err := r.resolvedValueName()
	if err != nil {
		return err
	}
	r.value = T(request.PathValue(name))
	return nil
}
