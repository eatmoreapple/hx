package binding

import (
	"mime"
	"net/http"
	"strings"
)

// Common MIME types for request Content-Type
// These constants define the standard MIME types used for HTTP request content identification.
const (
	MIMEJSON          = "application/json"                  // MIMEJSON represents JSON content type
	MIMEMultipartForm = "multipart/form-data"               // MIMEMultipartForm represents multipart form data (typically used for file uploads)
	MIMEPOSTForm      = "application/x-www-form-urlencoded" // MIMEPOSTForm represents URL-encoded form data
	XMLMIME           = "application/xml"                   // XMLMIME represents XML content type
)

// Common binders for common MIME types
// These pre-initialized binder instances are used to avoid creating new binders for each request.
var (
	jsonBinder  = JSONBinder{}  // jsonBinder handles binding of JSON request bodies
	xmlBinder   = XMLBinder{}   // xmlBinder handles binding of XML request bodies
	formBinder  = FormBinder{}  // formBinder handles binding of form data (both multipart and URL-encoded)
	queryBinder = QueryBinder{} // queryBinder handles binding of URL query parameters
)

type Binder interface {
	Bind(*http.Request, any) error
}

// Default returns the appropriate binder based on the HTTP method and Content-Type header.
// Content-Type parsing follows RFC 7231, section 3.1.1.1 and RFC 2045.
// Examples of valid Content-Type values:
//   - application/json
//   - application/x-www-form-urlencoded
//   - multipart/form-data; boundary=something
//
// If the Content-Type header is invalid or not provided, it defaults to QueryBinder.
// GET requests always use QueryBinder regardless of Content-Type.
func Default(method, contentType string) Binder {
	// GET requests always use query parameters
	if method == http.MethodGet {
		return QueryBinder{}
	}

	// Parse media type according to RFC 7231
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return QueryBinder{} // Invalid Content-Type defaults to query
	}

	// Media type comparison should be case-insensitive (RFC 2045)
	mediaType = strings.ToLower(mediaType)
	switch mediaType {
	case MIMEJSON:
		return jsonBinder
	case XMLMIME:
		return xmlBinder
	case MIMEMultipartForm, MIMEPOSTForm:
		return formBinder // Both form types use the same binder
	default:
		return queryBinder
	}
}
