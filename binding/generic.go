package binding

import (
	"net/http"
	"reflect"

	"github.com/eatmoreapple/hx/httpx"
)

// GenericBinder is a utility for binding HTTP request data to a struct.
// It iterates over the fields of the provided struct and checks if they implement
// the `httpx.RequestExtractor` interface. If a field implements the interface,
// the `FromRequest` method is called to populate the field with data from the HTTP request.
type GenericBinder struct{}

// Bind processes the HTTP request and populates the provided struct (`a`) with data.
// It uses reflection to inspect the struct fields and checks if they implement the
// `httpx.RequestExtractor` interface. If a field implements the interface, the
// `FromRequest` method is invoked to extract and set the data from the request.
//
// Parameters:
//   - r: The HTTP request containing the data to be bound.
//   - a: A pointer to the struct that will be populated with the request data.
//
// Returns:
//   - An error if any field implementing `httpx.RequestExtractor` fails to extract data.
//   - nil if the binding process completes successfully.
func (g GenericBinder) Bind(r *http.Request, a any) error {
	// Use reflection to get the underlying value of the struct.
	v := reflect.Indirect(reflect.ValueOf(a))
	// If the provided value is not a struct, return early.
	if v.Kind() != reflect.Struct {
		return nil
	}

	// Iterate over each field in the struct.
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		isPointer := field.Kind() == reflect.Ptr
		isImplementedRequestExtractor := httpx.IsRequestExtractorType(field.Type())

		// If the field implements `httpx.RequestExtractor`, process it.
		if isImplementedRequestExtractor {
			// If the field is a pointer and is nil, initialize it with a new instance of its type.
			if isPointer && field.IsNil() {
				field.Set(reflect.New(field.Type().Elem()))
			} else {
				// If the field is not a pointer, convert it to a pointer.
				field = field.Addr()
			}
			// Call the `FromRequest` method to extract data from the request and populate the field.
			if err := field.Interface().(httpx.RequestExtractor).FromRequest(r); err != nil {
				return err
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
