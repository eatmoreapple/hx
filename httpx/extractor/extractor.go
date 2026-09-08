package extractor

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"reflect"
)

// RequestExtractor defines the interface for types that can extract data from HTTP requests.
// Implementations should handle parsing and validating request data.
type RequestExtractor interface {
	FromRequest(*http.Request) error
}

// FieldRequestExtractor extends RequestExtractor with struct-field context.
// Struct binders use FromRequestField in place of FromRequest when this
// interface is implemented.
type FieldRequestExtractor interface {
	RequestExtractor
	FromRequestField(*http.Request, reflect.StructField) error
}

type Empty struct{}

func (e *Empty) FromRequest(*http.Request) error { return nil }

// baseValueExtractor provides common functionality for value extractors.
// It implements basic operations like value retrieval and JSON marshaling.
type baseValueExtractor[T Value] struct {
	raw   string
	value T
}

// valueNameFromField resolves a request value name from struct metadata.
// The generic hx tag takes precedence over the extractor-specific tag, with
// the Go field name used as the final fallback.
func valueNameFromField(field reflect.StructField, extractorTag string) string {
	if name := field.Tag.Get(ValueNameTag); name != "" {
		return name
	}
	if name := field.Tag.Get(extractorTag); name != "" {
		return name
	}
	return field.Name
}

func allowEmptyFromField(field reflect.StructField) bool {
	return field.Tag.Get("required") != "true"
}

// resolvedValueName returns the name supplied by the value type or the
// caller-provided fallback name.
func (b baseValueExtractor[T]) resolvedValueName(fallback string) (string, error) {
	if namer, ok := any(*new(T)).(ValueNamer); ok {
		name := namer.ValueName()
		if name == "" {
			return "", ErrValueNameRequired
		}
		return name, nil
	}
	if fallback != "" {
		return fallback, nil
	}
	return "", ErrValueNameRequired
}

func (b *baseValueExtractor[T]) set(s string, allowEmpty bool) error {
	v, err := parse[T](s, allowEmpty)
	if err != nil {
		return err
	}
	b.raw = s
	b.value = v
	return nil
}

// Value returns the extracted value.
// This method should be called after FromRequest has been executed successfully.
func (b baseValueExtractor[T]) Value() T {
	return b.value
}

// MarshalJSON implements json.Marshaler interface to provide JSON serialization
// of the extracted value.
func (b baseValueExtractor[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Value())
}

func (b *baseValueExtractor[T]) UnmarshalJSON(_ []byte) error {
	return nil
}

// UnmarshalXML implements xml.Unmarshaler by consuming and ignoring the
// extractor's XML element. Request extractors are populated from the HTTP
// request rather than the request body.
func (b *baseValueExtractor[T]) UnmarshalXML(decoder *xml.Decoder, _ xml.StartElement) error {
	return decoder.Skip()
}

// UnmarshalXMLAttr implements xml.UnmarshalerAttr by ignoring the extractor's
// XML attribute.
func (b *baseValueExtractor[T]) UnmarshalXMLAttr(xml.Attr) error {
	return nil
}

// UnmarshalForm ignores form and query values so the extractor remains
// responsible for populating itself from the HTTP request.
func (b *baseValueExtractor[T]) UnmarshalForm([]string) error {
	return nil
}

// String returns the value as a string.
// This is a convenience method that simply converts the value to string.
func (b baseValueExtractor[T]) String() string {
	return b.raw
}

// FromRequest is a placeholder implementation that should be overridden by embedding types.
// It returns an error indicating that the method is not supported.
func (b *baseValueExtractor[T]) FromRequest(*http.Request) error {
	return errors.ErrUnsupported
}
