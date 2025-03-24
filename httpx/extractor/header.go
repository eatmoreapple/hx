package extractor

import "net/http"

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

type HeaderExtractor http.Header

func (r *HeaderExtractor) FromRequest(request *http.Request) error {
	*r = HeaderExtractor(request.Header)
	return nil
}
