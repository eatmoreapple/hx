package binding

import (
	"cmp"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

// Common errors that can occur during binding
var (
	ErrPointerRequired = errors.New("binding: destination must be a pointer")
	ErrStructRequired  = errors.New("binding: destination must be a struct")
	ErrUnsupportedType = errors.New("binding: unsupported type")
	ErrTooManyFields   = errors.New("binding: too many fields")
)

const (
	maxFields = 1000 // Maximum number of fields to prevent DOS attacks
)

// mapTo maps url.Values to a struct using reflection.
// The struct fields should be tagged with "form" tags.
// If a field's tag is "-", it will be skipped.
func mapTo(values url.Values, dest any) error {
	if len(values) > maxFields {
		return ErrTooManyFields
	}

	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		return ErrPointerRequired
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return ErrStructRequired
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := cmp.Or(f.Tag.Get("form"), f.Name)
		if tag == "-" { // skip this field
			continue
		}
		if value, ok := values[tag]; ok {
			if err := setTo(v.Field(i), value); err != nil {
				return fmt.Errorf("binding field %q: %w", f.Name, err)
			}
		}
	}
	return nil
}

// setTo sets a reflect.Value from a slice of strings
func setTo(field reflect.Value, value []string) error {
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		field = field.Elem()
	}

	switch field.Kind() {
	case reflect.Slice:
		return bindSlice(field, value)
	default:
		if len(value) == 0 {
			field.Set(reflect.Zero(field.Type()))
			return nil
		}
		return setValue(field, value[0])
	}
}

// bindSlice handles binding of slice types
func bindSlice(field reflect.Value, formValue []string) error {
	if len(formValue) > maxFields {
		return ErrTooManyFields
	}

	if field.Type().Elem().Kind() == reflect.String {
		field.Set(reflect.ValueOf(formValue))
		return nil
	}

	if field.Type().Elem().Kind() == reflect.Ptr {
		return bindPtrSlice(field, formValue)
	}
	return bindValueSlice(field, formValue)
}

// bindPtrSlice handles binding of slices of pointers
func bindPtrSlice(field reflect.Value, formValue []string) error {
	slice := reflect.MakeSlice(field.Type(), len(formValue), len(formValue))
	for i, v := range formValue {
		ptr := reflect.New(field.Type().Elem().Elem())
		if err := setValue(ptr.Elem(), v); err != nil {
			return fmt.Errorf("binding slice element %d: %w", i, err)
		}
		slice.Index(i).Set(ptr)
	}
	field.Set(slice)
	return nil
}

// bindValueSlice handles binding of slices of values
func bindValueSlice(field reflect.Value, formValue []string) error {
	slice := reflect.MakeSlice(field.Type(), len(formValue), len(formValue))
	for i, v := range formValue {
		if err := setValue(slice.Index(i), v); err != nil {
			return fmt.Errorf("binding slice element %d: %w", i, err)
		}
	}
	field.Set(slice)
	return nil
}

// bindInt binds a string to an int field
func bindInt(field reflect.Value, formValue string, bitSize int) error {
	if formValue == "" {
		field.SetInt(0)
		return nil
	}
	v, err := strconv.ParseInt(formValue, 10, bitSize)
	if err != nil {
		return fmt.Errorf("parsing int: %w", err)
	}
	field.SetInt(v)
	return nil
}

// bindUint binds a string to a uint field
func bindUint(field reflect.Value, formValue string, bitSize int) error {
	if formValue == "" {
		field.SetUint(0)
		return nil
	}
	v, err := strconv.ParseUint(formValue, 10, bitSize)
	if err != nil {
		return fmt.Errorf("parsing uint: %w", err)
	}
	field.SetUint(v)
	return nil
}

// bindFloat binds a string to a float field
func bindFloat(field reflect.Value, formValue string, bitSize int) error {
	if formValue == "" {
		field.SetFloat(0)
		return nil
	}
	v, err := strconv.ParseFloat(formValue, bitSize)
	if err != nil {
		return fmt.Errorf("parsing float: %w", err)
	}
	field.SetFloat(v)
	return nil
}

// bindBool binds a string to a bool field
func bindBool(field reflect.Value, formValue string) error {
	if formValue == "" {
		field.SetBool(false)
		return nil
	}
	v, err := strconv.ParseBool(formValue)
	if err != nil {
		return fmt.Errorf("parsing bool: %w", err)
	}
	field.SetBool(v)
	return nil
}

// setValue sets a field's value from a string
func setValue(field reflect.Value, formValue string) error {
	if formValue == "" {
		field.Set(reflect.Zero(field.Type()))
		return nil
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(formValue)
	case reflect.Int:
		return bindInt(field, formValue, 0)
	case reflect.Int8:
		return bindInt(field, formValue, 8)
	case reflect.Int16:
		return bindInt(field, formValue, 16)
	case reflect.Int32:
		return bindInt(field, formValue, 32)
	case reflect.Int64:
		return bindInt(field, formValue, 64)
	case reflect.Uint:
		return bindUint(field, formValue, 0)
	case reflect.Uint8:
		return bindUint(field, formValue, 8)
	case reflect.Uint16:
		return bindUint(field, formValue, 16)
	case reflect.Uint32:
		return bindUint(field, formValue, 32)
	case reflect.Uint64:
		return bindUint(field, formValue, 64)
	case reflect.Float32:
		return bindFloat(field, formValue, 32)
	case reflect.Float64:
		return bindFloat(field, formValue, 64)
	case reflect.Bool:
		return bindBool(field, formValue)
	default:
		return fmt.Errorf("%w: %s", ErrUnsupportedType, field.Kind())
	}
	return nil
}
