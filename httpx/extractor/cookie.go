package extractor

import "net/http"

// CookieValueExtractor implements RequestExtractor for cookie values.
// It extracts and stores cookie values of a specified type T that implements the Value interface.
type CookieValueExtractor[T Value] struct {
	baseValueExtractor[T]
}

// FromRequest implements RequestExtractor.FromRequest by extracting the cookie value
// using the name provided by ValueName(). The cookie value is converted to type T.
func (r *CookieValueExtractor[T]) FromRequest(request *http.Request) error {
	cookie, err := request.Cookie(r.value.ValueName())
	if err != nil {
		return err
	}
	r.value = T(cookie.Value)
	return nil
}

type CookieExtractor []*http.Cookie

func (r *CookieExtractor) FromRequest(request *http.Request) error {
	*r = request.Cookies()
	return nil
}
