package extractor

import (
	"net/http"
	"net/url"
	"reflect"
)

// FormValueExtractor implements RequestExtractor for form values.
// It extracts and stores form values of a specified string-like type T.
type FormValueExtractor[T Value] struct {
	baseValueExtractor[T]
}

// FromRequest implements RequestExtractor.FromRequest by extracting the form value
// using the resolved value name. The form value is converted to type T.
func (r *FormValueExtractor[T]) FromRequest(request *http.Request) error {
	return r.fromRequest(request, "")
}

// FromRequestField extracts a form value with the containing field's hx tag as
// a fallback name.
func (r *FormValueExtractor[T]) FromRequestField(request *http.Request, field reflect.StructField) error {
	return r.fromRequest(request, field.Tag.Get(ValueNameTag))
}

func (r *FormValueExtractor[T]) fromRequest(request *http.Request, fallbackName string) error {
	name, err := r.resolvedValueName(fallbackName)
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
