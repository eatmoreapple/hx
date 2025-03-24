package extractor

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Value is an interface for types that can be used as path parameters.
// It combines the PathValueName method with the constraint of being a string type.
type Value interface {
	// ValueName returns the name of the path parameter as defined in the route.
	ValueName() string
	~string
}

// baseValueExtractor provides common functionality for value extractors.
// It implements basic operations like value retrieval and JSON marshaling.
type baseValueExtractor[T Value] struct {
	value T // The extracted value after processing
}

// Value returns the extracted value.
// This method should be called after FromRequest has been executed successfully.
func (b baseValueExtractor[T]) Value() T {
	return b.value
}

// MarshalJSON implements json.Marshaler interface to provide JSON serialization
// of the extracted value.
func (b baseValueExtractor[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.value)
}

// Int8 converts the value to int8.
// Returns an error if the value cannot be parsed as an 8-bit integer.
func (b baseValueExtractor[T]) Int8() (int8, error) {
	v, err := strconv.ParseInt(string(b.value), 10, 8)
	return int8(v), err
}

// Int16 converts the value to int16.
// Returns an error if the value cannot be parsed as a 16-bit integer.
func (b baseValueExtractor[T]) Int16() (int16, error) {
	v, err := strconv.ParseInt(string(b.value), 10, 16)
	return int16(v), err
}

// Int32 converts the value to int32.
// Returns an error if the value cannot be parsed as an integer.
func (b baseValueExtractor[T]) Int32() (int32, error) {
	v, err := strconv.ParseInt(string(b.value), 10, 32)
	return int32(v), err
}

// Int64 converts the value to int64.
// Returns an error if the value cannot be parsed as an integer.
func (b baseValueExtractor[T]) Int64() (int64, error) {
	return strconv.ParseInt(string(b.value), 10, 64)
}

// Int converts the value to int.
// Returns an error if the value cannot be parsed as an integer.
func (b baseValueExtractor[T]) Int() (int, error) {
	v, err := strconv.ParseInt(string(b.value), 10, 0)
	return int(v), err
}

// Uint8 converts the value to uint8.
// Returns an error if the value cannot be parsed as an 8-bit unsigned integer.
func (b baseValueExtractor[T]) Uint8() (uint8, error) {
	v, err := strconv.ParseUint(string(b.value), 10, 8)
	return uint8(v), err
}

// Uint16 converts the value to uint16.
// Returns an error if the value cannot be parsed as a 16-bit unsigned integer.
func (b baseValueExtractor[T]) Uint16() (uint16, error) {
	v, err := strconv.ParseUint(string(b.value), 10, 16)
	return uint16(v), err
}

// Uint32 converts the value to uint32.
// Returns an error if the value cannot be parsed as an unsigned integer.
func (b baseValueExtractor[T]) Uint32() (uint32, error) {
	v, err := strconv.ParseUint(string(b.value), 10, 32)
	return uint32(v), err
}

// Uint64 converts the value to uint64.
// Returns an error if the value cannot be parsed as an unsigned integer.
func (b baseValueExtractor[T]) Uint64() (uint64, error) {
	return strconv.ParseUint(string(b.value), 10, 64)
}

// Uint converts the value to uint.
// Returns an error if the value cannot be parsed as an unsigned integer.
func (b baseValueExtractor[T]) Uint() (uint, error) {
	v, err := strconv.ParseUint(string(b.value), 10, 0)
	return uint(v), err
}

// Float64 converts the value to float64.
// Returns an error if the value cannot be parsed as a floating-point number.
func (b baseValueExtractor[T]) Float64() (float64, error) {
	return strconv.ParseFloat(string(b.value), 64)
}

// Float32 converts the value to float32.
// Returns an error if the value cannot be parsed as a floating-point number.
func (b baseValueExtractor[T]) Float32() (float32, error) {
	v, err := strconv.ParseFloat(string(b.value), 32)
	return float32(v), err
}

// Bool converts the value to bool.
// Returns an error if the value cannot be parsed as a boolean.
// Accepts 1, t, T, TRUE, true for true and 0, f, F, FALSE, false for false.
func (b baseValueExtractor[T]) Bool() (bool, error) {
	return strconv.ParseBool(string(b.value))
}

// String returns the value as a string.
// This is a convenience method that simply converts the value to string.
func (b baseValueExtractor[T]) String() string {
	return string(b.value)
}

// FromRequest is a placeholder implementation that should be overridden by embedding types.
// It will panic if called directly.
func (b *baseValueExtractor[T]) FromRequest(*http.Request) error {
	panic("not implemented")
}
