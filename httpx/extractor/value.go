package extractor

import (
	"errors"
	"reflect"
	"strconv"
	"uuid"
)

// ValueNameTag is the struct tag used to provide an extractor value name when
// the value type does not implement ValueNamer.
const ValueNameTag = "hx"

// ErrValueNameRequired is returned when an extractor cannot resolve a
// non-empty request value name.
var ErrValueNameRequired = errors.New(`extractor: non-empty value name is required; implement ValueName or set the "hx" struct tag`)

var errUnsupportedValueType = errors.New("unsupported type")

// Value is the set of types a single request value can be converted into.
type Value interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
	~float32 | ~float64 |
	~string |
	~bool |
	uuid.UUID
}

// parse converts a request string into T.
//
// The type switch only matches predeclared types (int, string, uuid.UUID, …).
// Defined types such as `type UserID int` are allowed by Value so they can
// implement ValueNamer, but they do not match those cases; they fall through
// to parseReflect, which converts using the underlying kind.
func parse[T Value](value string) (T, error) {
	switch any(*new(T)).(type) {
	case int:
		v, err := strconv.Atoi(value)
		return any(v).(T), err
	case int8:
		v, err := strconv.ParseInt(value, 10, 8)
		return any(int8(v)).(T), err
	case int16:
		v, err := strconv.ParseInt(value, 10, 16)
		return any(int16(v)).(T), err
	case int32:
		v, err := strconv.ParseInt(value, 10, 32)
		return any(int32(v)).(T), err
	case int64:
		v, err := strconv.ParseInt(value, 10, 64)
		return any(v).(T), err
	case uint:
		v, err := strconv.ParseUint(value, 10, 0)
		return any(uint(v)).(T), err
	case uint8:
		v, err := strconv.ParseUint(value, 10, 8)
		return any(uint8(v)).(T), err
	case uint16:
		v, err := strconv.ParseUint(value, 10, 16)
		return any(uint16(v)).(T), err
	case uint32:
		v, err := strconv.ParseUint(value, 10, 32)
		return any(uint32(v)).(T), err
	case uint64:
		v, err := strconv.ParseUint(value, 10, 64)
		return any(uint64(v)).(T), err
	case float32:
		v, err := strconv.ParseFloat(value, 32)
		return any(float32(v)).(T), err
	case float64:
		v, err := strconv.ParseFloat(value, 64)
		return any(v).(T), err
	case string:
		return any(value).(T), nil
	case bool:
		v, err := strconv.ParseBool(value)
		return any(v).(T), err
	case uuid.UUID:
		v, err := uuid.Parse(value)
		return any(v).(T), err
	default:
		return parseReflect[T](value)
	}
}

// parseReflect handles defined types that the type switch cannot see.
func parseReflect[T Value](value string) (T, error) {
	var dest T
	rv := reflect.ValueOf(&dest).Elem()
	switch rv.Kind() {
	case reflect.String:
		rv.SetString(value)
		return dest, nil
	case reflect.Bool:
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return dest, err
		}
		rv.SetBool(parsed)
		return dest, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsed, err := strconv.ParseInt(value, 10, rv.Type().Bits())
		if err != nil {
			return dest, err
		}
		rv.SetInt(parsed)
		return dest, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		parsed, err := strconv.ParseUint(value, 10, rv.Type().Bits())
		if err != nil {
			return dest, err
		}
		rv.SetUint(parsed)
		return dest, nil
	case reflect.Float32, reflect.Float64:
		parsed, err := strconv.ParseFloat(value, rv.Type().Bits())
		if err != nil {
			return dest, err
		}
		rv.SetFloat(parsed)
		return dest, nil
	default:
		return dest, errUnsupportedValueType
	}
}

// ValueNamer provides the name of the request value to extract.
type ValueNamer interface {
	ValueName() string
}

// NamedValue is a string-like value that provides its own request value name.
type NamedValue interface {
	Value
	ValueNamer
}
