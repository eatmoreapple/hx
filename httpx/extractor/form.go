package extractor

import (
	"net/http"
	"net/url"
)

// FormValueExtractor implements RequestExtractor for form values.
// It extracts and stores form values of a specified string-like type T.
type FormValueExtractor[T Value] struct {
	baseValueExtractor[T]
}

// FromRequest implements RequestExtractor.FromRequest by extracting the form value
// using the resolved value name. The form value is converted to type T.
func (r *FormValueExtractor[T]) FromRequest(request *http.Request) error {
	name, err := r.resolvedValueName()
	if err != nil {
		return err
	}
	r.value = T(request.FormValue(name))
	return nil
}

// FormExtractor is a type alias for http.Request.Form
type FormExtractor url.Values

func (r *FormExtractor) FromRequest(request *http.Request) error {
	*r = FormExtractor(request.Form)
	return nil
}
