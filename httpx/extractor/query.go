package extractor

import (
	"net/http"
	"net/url"
)

// QueryValueExtractor implements RequestExtractor for query parameters.
// It extracts and stores query values of a specified type T that implements the Value interface.
type QueryValueExtractor[T Value] struct {
	baseValueExtractor[T]
}

// FromRequest implements RequestExtractor.FromRequest by extracting the query value
// using the name provided by ValueName(). The query value is converted to type T.
func (r *QueryValueExtractor[T]) FromRequest(request *http.Request) error {
	r.value = T(request.URL.Query().Get(r.value.ValueName()))
	return nil
}

// FromQuery is a type alias for QueryValueExtractor providing a shorter name
// while maintaining all functionality.
type FromQuery[T Value] = QueryValueExtractor[T]

// QueryExtractor is a type alias for http.URL.Query providing a shorter name
// while maintaining all functionality.
type QueryExtractor url.Values

// FromRequest implements RequestExtractor.FromRequest by extracting the query values
// from the request URL. It populates the QueryExtractor with the query values.
func (r *QueryExtractor) FromRequest(request *http.Request) error {
	*r = QueryExtractor(request.URL.Query())
	return nil
}
