package extractor

import (
	"mime"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

const (
	mimeMultipartForm = "multipart/form-data"
	defaultMaxMemory  = 32 << 20
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

// FromRequestField extracts a form value using the containing field's tags or
// Go name as a fallback value name.
func (r *FormValueExtractor[T]) FromRequestField(request *http.Request, field reflect.StructField) error {
	return r.fromRequest(request, valueNameFromField(field, "form"))
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
	contentType, _, err := mime.ParseMediaType(request.Header.Get("Content-Type"))
	if err == nil && strings.EqualFold(contentType, mimeMultipartForm) {
		if err := request.ParseMultipartForm(defaultMaxMemory); err != nil {
			return err
		}
	} else if err := request.ParseForm(); err != nil {
		return err
	}

	*r = FormExtractor(request.Form)
	return nil
}
