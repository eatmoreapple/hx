package extractor

import (
	"net/http"
	"net/url"
	"reflect"
)

// QueryValueExtractor implements RequestExtractor for query parameters.
// It extracts and stores query values of a specified string-like type T.
type QueryValueExtractor[T Value] struct {
	baseValueExtractor[T]
}

// FromRequest implements RequestExtractor.FromRequest by extracting the query value
// using the resolved value name. The query value is converted to type T.
func (r *QueryValueExtractor[T]) FromRequest(request *http.Request) error {
	return r.fromRequest(request, "")
}

// FromRequestField extracts a query value with the containing field's hx tag
// as a fallback name.
func (r *QueryValueExtractor[T]) FromRequestField(request *http.Request, field reflect.StructField) error {
	return r.fromRequest(request, field.Tag.Get(ValueNameTag))
}

func (r *QueryValueExtractor[T]) fromRequest(request *http.Request, fallbackName string) error {
	name, err := r.resolvedValueName(fallbackName)
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
