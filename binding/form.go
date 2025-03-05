package binding

import (
	"cmp"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"
)

var (
	// fileHeaderType is the reflect type for *multipart.FileHeader.
	fileHeaderType = reflect.TypeFor[*multipart.FileHeader]()

	// fileHeaderSliceType is the reflect type for []*multipart.FileHeader.
	fileHeaderSliceType = reflect.TypeFor[[]*multipart.FileHeader]()
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
	for k, v := range r.URL.Query() {
		values[k] = v
	}

	// Add form values
	for k, v := range r.Form {
		values[k] = v
	}

	// Add multipart form values if any
	if r.MultipartForm != nil {
		for k, v := range r.MultipartForm.Value {
			values[k] = v
		}

		// Handle file uploads if the destination struct has multipart.FileHeader fields
		if len(r.MultipartForm.File) > 0 {
			if err := handleFileUploads(r.MultipartForm.File, dest); err != nil {
				return err
			}
		}
	}

	return mapTo(values, dest)
}

// handleFileUploads processes file uploads in multipart forms
func handleFileUploads(files map[string][]*multipart.FileHeader, dest any) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		return ErrPointerRequired
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return ErrStructRequired
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		if field.Type == fileHeaderType || field.Type == fileHeaderSliceType {
			tag := cmp.Or(field.Tag.Get("form"), field.Name)
			if file, ok := files[tag]; ok {
				if field.Type == fileHeaderType {
					v.Field(i).Set(reflect.ValueOf(file[0]))
				} else {
					v.Field(i).Set(reflect.ValueOf(file))
				}
			}
		}
	}
	return nil
}
