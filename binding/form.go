package binding

import (
	"mime"
	"net/http"
	"strings"

	"github.com/eatmoreapple/hx/encoding/form"
)

// FormBinder handles both application/x-www-form-urlencoded and multipart/form-data
type FormBinder struct{}

// Bind implements the Binder interface for form data.
// It handles both url-encoded forms and multipart forms.
func (f FormBinder) Bind(r *http.Request, dest any) error {
	// ParseMultipartForm also populates r.Form, so choose one parsing path
	// instead of calling ParseForm before it for multipart requests.
	contentType, _, mediaTypeErr := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if mediaTypeErr == nil && strings.EqualFold(contentType, MIMEMultipartForm) {
		if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB max memory
			return err
		}
	} else if err := r.ParseForm(); err != nil {
		return err
	}

	if r.MultipartForm != nil {
		// Handle file uploads if the destination struct has multipart.FileHeader fields
		if len(r.MultipartForm.File) > 0 {
			if err := form.UnmarshalFiles(r.MultipartForm.File, dest); err != nil {
				return err
			}
		}
	}

	return form.Unmarshal(r.Form, dest)
}
