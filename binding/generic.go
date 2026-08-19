package binding

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/eatmoreapple/hx/httpx"
)

// GenericBinder is a utility for binding HTTP request data to a struct.
// It uses FromRequestField when an extractor supports struct-field context and
// falls back to FromRequest otherwise. Anonymous value-struct fields that do
// not implement RequestExtractor are recursively bound.
type GenericBinder struct{}

// ErrGenericBinderTarget indicates that Bind received anything other than a
// non-nil pointer to a struct.
var ErrGenericBinderTarget = errors.New("generic binder target must be a non-nil pointer to a struct")

// Bind processes the HTTP request and populates the provided struct (`a`) with data.
// It uses reflection to inspect the struct fields and checks if they implement the
// `httpx.RequestExtractor` interface. Field-aware implementations receive the
// containing struct field while regular implementations receive only the request.
//
// Parameters:
//   - r: The HTTP request containing the data to be bound.
//   - a: A pointer to the struct that will be populated with the request data.
//
// Returns:
//   - An error if any field implementing `httpx.RequestExtractor` fails to extract data.
//   - nil if the binding process completes successfully.
func (g GenericBinder) Bind(r *http.Request, a any) error {
	value := reflect.ValueOf(a)
	if !value.IsValid() || value.Kind() != reflect.Pointer || value.IsNil() || value.Elem().Kind() != reflect.Struct {
		return ErrGenericBinderTarget
	}

	// Use reflection to get the underlying value of the struct.
	v := value.Elem()
	return g.bindValue(r, v)
}

// bindValue recursively binds request data into a reflected struct value.
// Keeping the recursion at the reflect.Value level avoids converting embedded
// structs to interfaces and then reflecting them again.
func (g GenericBinder) bindValue(r *http.Request, v reflect.Value) error {
	for structField, field := range v.Fields() {
		// Unexported fields cannot be safely initialized or addressed through
		// reflection, so leave them untouched.
		if !field.CanSet() {
			continue
		}

		isImplementedRequestExtractor := httpx.IsRequestExtractorType(field.Type())
		// If the field implements `httpx.RequestExtractor`, process it.
		if isImplementedRequestExtractor {
			isPointer := field.Kind() == reflect.Pointer

			// If the field is a pointer and is nil, initialize it with a new instance of its type.
			if isPointer {
				if field.IsNil() {
					field.Set(reflect.New(field.Type().Elem()))
				}
			} else {
				// If the field is not a pointer, convert it to a pointer.
				field = field.Addr()
			}
			// Prefer field-aware extraction when the extractor supports it.
			extractor, _ := reflect.TypeAssert[httpx.RequestExtractor](field)

			var extractErr error
			if fieldExtractor, ok := extractor.(httpx.FieldRequestExtractor); ok {
				extractErr = fieldExtractor.FromRequestField(r, structField)
			} else {
				extractErr = extractor.FromRequest(r)
			}
			if extractErr != nil {
				return fmt.Errorf("binding field %q: %w", structField.Name, extractErr)
			}
			continue
		}

		// Treat an anonymous value struct as an embedded binding group. A
		// pointer-to-struct is intentionally not handled here: allocating and
		// recursively traversing nil embedded pointers can create an infinite
		// recursion for recursive types.
		if structField.Anonymous && field.Kind() == reflect.Struct {
			if bindErr := g.bindValue(r, field); bindErr != nil {
				return fmt.Errorf("binding embedded field %q: %w", structField.Name, bindErr)
			}
		}
	}

	return nil
}

// generic is a singleton instance of GenericBinder.
// It's used as a shared instance since GenericBinder has no state.
var generic = &GenericBinder{}

// Generic returns a shared instance of GenericBinder.
// Since GenericBinder is stateless, this singleton pattern is safe for concurrent use.
func Generic() Binder {
	return generic
}
