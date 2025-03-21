// Package serializer provides functionality for serializing and deserializing data
// between Go data structures and various formats like JSON.
// It defines a generic interface for serialization and provides a default implementation
// using Go's standard encoding/json package.
package serializer

import (
	"encoding/json"
	"io"
)

// Serializer defines an interface for encoding and decoding data.
// Implementations of this interface can support different serialization formats,
// such as JSON, XML, or custom binary formats.
type Serializer interface {
	// Serialize encodes the value v into the specified format and writes it to the provided writer.
	// The value v can be any Go data structure that the serializer supports.
	// The writer w is where the serialized data will be written.
	// Returns an error if the serialization process fails.
	Serialize(v any, w io.Writer) error

	// Deserialize reads data from the provided reader and decodes it into the value pointed to by v.
	// The v parameter should be a pointer to the target object where the deserialized data will be stored.
	// Returns an error if the deserialization process fails.
	Deserialize(r io.Reader, v any) error
}

// StdJSONSerializer implements the Serializer interface using Go's standard
// encoding/json package for JSON serialization and deserialization.
type StdJSONSerializer struct{}

// Serialize encodes the value v as JSON and writes it to the provided writer w.
// This method uses Go's standard JSON encoder to perform the serialization.
// Returns an error if the encoding process fails.
func (s *StdJSONSerializer) Serialize(v any, w io.Writer) error {
	return json.NewEncoder(w).Encode(v)
}

// Deserialize reads JSON data from the provided reader r and decodes it into the value pointed to by v.
// This method uses Go's standard JSON decoder to perform the deserialization.
// Returns an error if the decoding process fails.
func (s *StdJSONSerializer) Deserialize(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}

// stdJSONSerializer is a singleton instance of StdJSONSerializer.
// This instance is used as the default JSON serializer for the package.
var stdJSONSerializer Serializer = &StdJSONSerializer{}

// JSONSerializer returns a singleton instance of StdJSONSerializer
// that implements the Serializer interface using the standard JSON library.
// This function provides a convenient way to access the default JSON serializer.
func JSONSerializer() Serializer {
	return stdJSONSerializer
}

// SetJSONSerializer sets the global JSON serializer instance to the provided serializer s.
// This function allows customization of the JSON serialization behavior by replacing
// the default StdJSONSerializer with a custom implementation.
// Panics if the provided serializer is nil, as a nil serializer is not valid.
func SetJSONSerializer(s Serializer) {
	if s == nil {
		panic("serializer cannot be nil")
	}
	stdJSONSerializer = s
}
