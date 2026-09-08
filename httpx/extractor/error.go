package extractor

import "errors"

// ExtractError wraps a failure from FromRequest or FromRequestField so callers
// can distinguish extractor/parse errors from other binding or handler errors.
type ExtractError struct {
	Field string
	Err   error
}

func (e *ExtractError) Error() string {
	if e == nil || e.Err == nil {
		return "extractor: extract"
	}
	if e.Field != "" {
		return "extractor: field " + e.Field + ": " + e.Err.Error()
	}
	return "extractor: " + e.Err.Error()
}

func (e *ExtractError) Unwrap() error { return e.Err }

// WrapExtractError returns err wrapped as an ExtractError.
// A nil error is returned unchanged. An error that already unwraps to
// ExtractError is returned as-is so nested binding does not double-wrap.
func WrapExtractError(field string, err error) error {
	if err == nil {
		return nil
	}
	if _, ok := errors.AsType[*ExtractError](err); ok {
		return err
	}
	return &ExtractError{Field: field, Err: err}
}
