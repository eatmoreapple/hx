package binding

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/eatmoreapple/hx/httpx"
)

// GenericBinder is a utility for binding HTTP request data to a struct.
// It uses FromRequestField when an extractor supports struct-field context and
// falls back to FromRequest otherwise.
type GenericBinder struct{}

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
func (g GenericBinder) Bind(r *http.Request, a any) (err error) {
	// Use reflection to get the underlying value of the struct.
	v := reflect.Indirect(reflect.ValueOf(a))
	// If the provided value is not a struct, return early.
	if v.Kind() != reflect.Struct {
		return nil
	}

	for structField, field := range v.Fields() {
		isImplementedRequestExtractor := httpx.IsRequestExtractorType(field.Type())
		// If the field implements `httpx.RequestExtractor`, process it.
		if isImplementedRequestExtractor {
			isPointer := field.Kind() == reflect.Pointer

			// If the field is a pointer and is nil, initialize it with a new instance of its type.
			if isPointer {
				field.Set(reflect.New(field.Type().Elem()))
			} else {
				// If the field is not a pointer, convert it to a pointer.
				field = field.Addr()
			}
			// Prefer field-aware extraction when the extractor supports it.
			extractor, _ := reflect.TypeAssert[httpx.RequestExtractor](field)

			if fieldExtractor, ok := extractor.(httpx.FieldRequestExtractor); ok {
				err = fieldExtractor.FromRequestField(r, structField)
			} else {
				err = extractor.FromRequest(r)
			}
			if err != nil {
				return fmt.Errorf("binding field %q: %w", structField.Name, err)
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
