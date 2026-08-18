package binding

import (
	"maps"
	"net/http"
	"strings"

	"github.com/eatmoreapple/hx/encoding/form"
)

// FormBinder handles both application/x-www-form-urlencoded and multipart/form-data
type FormBinder struct{}

// Bind implements the Binder interface for form data.
// It handles both url-encoded forms and multipart forms.
func (f FormBinder) Bind(r *http.Request, dest any) error {
	// Parse the form data first
	if err := r.ParseForm(); err != nil {
		return err
	}

	// For multipart/form-data, also parse the multipart form
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, MIMEMultipartForm) {
		if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB max memory
			return err
		}
	}

	// Combine form values and multipart form values
	values := make(map[string][]string)

	// Add query parameters
	maps.Copy(values, r.URL.Query())

	// Add form values
	maps.Copy(values, r.Form)

	// Add multipart form values if any
	if r.MultipartForm != nil {
		maps.Copy(values, r.MultipartForm.Value)

		// Handle file uploads if the destination struct has multipart.FileHeader fields
		if len(r.MultipartForm.File) > 0 {
			if err := form.UnmarshalFiles(r.MultipartForm.File, dest); err != nil {
				return err
			}
		}
	}

	return form.Unmarshal(values, dest)
}
