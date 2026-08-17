package extractor

import (
	"net/http"
	"net/url"
)

// QueryValueExtractor implements RequestExtractor for query parameters.
// It extracts and stores query values of a specified string-like type T.
type QueryValueExtractor[T Value] struct {
	baseValueExtractor[T]
}

// FromRequest implements RequestExtractor.FromRequest by extracting the query value
// using the resolved value name. The query value is converted to type T.
func (r *QueryValueExtractor[T]) FromRequest(request *http.Request) error {
	name, err := r.resolvedValueName()
	if err != nil {
		return err
	}
	r.value = T(request.URL.Query().Get(name))
	return nil
}

// QueryExtractor is a type alias for http.URL.Query providing a shorter name
// while maintaining all functionality.
type QueryExtractor url.Values

// FromRequest implements RequestExtractor.FromRequest by extracting the query values
// from the request URL. It populates the QueryExtractor with the query values.
func (r *QueryExtractor) FromRequest(request *http.Request) error {
	*r = QueryExtractor(request.URL.Query())
	return nil
}
