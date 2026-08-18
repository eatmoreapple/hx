package form

import (
	"cmp"
	"mime/multipart"
	"reflect"
)

var (
	fileHeaderType      = reflect.TypeFor[*multipart.FileHeader]()
	fileHeaderSliceType = reflect.TypeFor[[]*multipart.FileHeader]()
)

// UnmarshalFiles decodes multipart file headers into matching struct fields.
// Supported field types are *multipart.FileHeader and
// []*multipart.FileHeader. Field names are resolved from the form tag, then
// the Go field name.
func UnmarshalFiles(files map[string][]*multipart.FileHeader, dest any) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Pointer {
		return ErrPointerRequired
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return ErrStructRequired
	}

	for structField, field := range v.Fields() {
		if structField.Type != fileHeaderType && structField.Type != fileHeaderSliceType {
			continue
		}

		name := cmp.Or(structField.Tag.Get("form"), structField.Name)
		if name == "-" {
			continue
		}
		matchedFiles, ok := files[name]
		if !ok || len(matchedFiles) == 0 {
			continue
		}

		if structField.Type == fileHeaderType {
			field.Set(reflect.ValueOf(matchedFiles[0]))
		} else {
			field.Set(reflect.ValueOf(matchedFiles))
		}
	}

	return nil
}
