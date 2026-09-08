package extractor

import (
	"errors"
	"net/http"
	"reflect"
)

// CookieValueExtractor implements RequestExtractor for cookie values.
// It extracts and stores cookie values of a specified string-like type T.
type CookieValueExtractor[T Value] struct {
	baseValueExtractor[T]
}

// FromRequest implements RequestExtractor.FromRequest by extracting the cookie value
// using the resolved value name. The cookie value is converted to type T.
func (r *CookieValueExtractor[T]) FromRequest(request *http.Request) error {
	return r.fromRequest(request, "", false)
}

// FromRequestField extracts a cookie value using the containing field's tags
// or Go name as a fallback value name.
func (r *CookieValueExtractor[T]) FromRequestField(request *http.Request, field reflect.StructField) error {
	return r.fromRequest(request, valueNameFromField(field, "cookie"), allowEmptyFromField(field))
}

func (r *CookieValueExtractor[T]) fromRequest(request *http.Request, fallbackName string, allowEmpty bool) error {
	name, err := r.resolvedValueName(fallbackName)
	if err != nil {
		return err
	}
	cookie, err := request.Cookie(name)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return r.set("", allowEmpty)
		}
		return err
	}
	return r.set(cookie.Value, allowEmpty)
}

type CookieExtractor []*http.Cookie

func (r *CookieExtractor) FromRequest(request *http.Request) error {
	*r = request.Cookies()
	return nil
}
