package extractor

import "net/http"

// HeaderValueExtractor implements RequestExtractor for HTTP header values.
// It extracts and stores header values of a specified string-like type T.
type HeaderValueExtractor[T Value] struct {
	baseValueExtractor[T]
}

// FromRequest implements RequestExtractor.FromRequest by extracting the header value
// using the resolved value name. The header value is converted to type T.
func (r *HeaderValueExtractor[T]) FromRequest(request *http.Request) error {
	name, err := r.resolvedValueName()
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
