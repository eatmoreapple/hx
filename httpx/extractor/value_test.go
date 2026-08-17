package extractor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

type TestValue string

func (t TestValue) ValueName() string {
	return "test"
}

func TestBaseValueExtractor(t *testing.T) {
	val := TestValue("123")
	extractor := baseValueExtractor[TestValue]{value: val}

	if extractor.Value() != val {
		t.Errorf("expected value %s, got %s", val, extractor.Value())
	}

	if extractor.String() != "123" {
		t.Errorf("expected string %s, got %s", "123", extractor.String())
	}

	// Test JSON Marshaling
	jsonBytes, err := json.Marshal(extractor)
	if err != nil {
		t.Errorf("unexpected error marshaling json: %v", err)
	}
	if string(jsonBytes) != `"123"` {
		t.Errorf("expected json %s, got %s", `"123"`, string(jsonBytes))
	}
}

func TestParseWith(t *testing.T) {
	type UserID int64

	e := baseValueExtractor[TestValue]{value: TestValue("123")}
	id, err := e.ParseWith(func(value string) (UserID, error) {
		parsed, err := strconv.ParseInt(value, 10, 64)
		return UserID(parsed), err
	})
	if err != nil {
		t.Fatalf("ParseWith returned an unexpected error: %v", err)
	}
	if id != UserID(123) {
		t.Fatalf("ParseWith returned %d, want %d", id, UserID(123))
	}

	_, err = e.ParseWith(func(value string) (int, error) {
		return strconv.Atoi(value + "x")
	})
	if err == nil {
		t.Fatal("ParseWith should return parser errors")
	}
}

func TestQueryValueExtractorFromRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/?test=from-request", nil)
	var extractor QueryValueExtractor[TestValue]

	if err := extractor.FromRequest(request); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := extractor.String(); got != "from-request" {
		t.Fatalf("expected %q, got %q", "from-request", got)
	}
}

func TestQueryValueExtractorFromRequestField(t *testing.T) {
	type requestFields struct {
		Query QueryValueExtractor[string] `hx:"query"`
	}

	field, _ := reflect.TypeFor[requestFields]().FieldByName("Query")
	request := httptest.NewRequest(http.MethodGet, "/?query=from-field", nil)
	var extractor QueryValueExtractor[string]

	if err := extractor.FromRequestField(request, field); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := extractor.String(); got != "from-field" {
		t.Fatalf("expected %q, got %q", "from-field", got)
	}
}

func TestValueConversions(t *testing.T) {
	tests := []struct {
		name  string
		value string
		check func(baseValueExtractor[TestValue]) error
	}{
		{
			name:  "Int",
			value: "123",
			check: func(e baseValueExtractor[TestValue]) error {
				v, err := e.Int()
				if err != nil {
					return err
				}
				if v != 123 {
					t.Errorf("expected 123, got %d", v)
				}
				return nil
			},
		},
		{
			name:  "Int8",
			value: "123",
			check: func(e baseValueExtractor[TestValue]) error {
				v, err := e.Int8()
				if err != nil {
					return err
				}
				if v != 123 {
					t.Errorf("expected 123, got %d", v)
				}
				return nil
			},
		},
		{
			name:  "Int16",
			value: "123",
			check: func(e baseValueExtractor[TestValue]) error {
				v, err := e.Int16()
				if err != nil {
					return err
				}
				if v != 123 {
					t.Errorf("expected 123, got %d", v)
				}
				return nil
			},
		},
		{
			name:  "Int32",
			value: "123",
			check: func(e baseValueExtractor[TestValue]) error {
				v, err := e.Int32()
				if err != nil {
					return err
				}
				if v != 123 {
					t.Errorf("expected 123, got %d", v)
				}
				return nil
			},
		},
		{
			name:  "Int64",
			value: "123",
			check: func(e baseValueExtractor[TestValue]) error {
				v, err := e.Int64()
				if err != nil {
					return err
				}
				if v != 123 {
					t.Errorf("expected 123, got %d", v)
				}
				return nil
			},
		},
		{
			name:  "Uint",
			value: "123",
			check: func(e baseValueExtractor[TestValue]) error {
				v, err := e.Uint()
				if err != nil {
					return err
				}
				if v != 123 {
					t.Errorf("expected 123, got %d", v)
				}
				return nil
			},
		},
		{
			name:  "Uint8",
			value: "123",
			check: func(e baseValueExtractor[TestValue]) error {
				v, err := e.Uint8()
				if err != nil {
					return err
				}
				if v != 123 {
					t.Errorf("expected 123, got %d", v)
				}
				return nil
			},
		},
		{
			name:  "Uint16",
			value: "123",
			check: func(e baseValueExtractor[TestValue]) error {
				v, err := e.Uint16()
				if err != nil {
					return err
				}
				if v != 123 {
					t.Errorf("expected 123, got %d", v)
				}
				return nil
			},
		},
		{
			name:  "Uint32",
			value: "123",
			check: func(e baseValueExtractor[TestValue]) error {
				v, err := e.Uint32()
				if err != nil {
					return err
				}
				if v != 123 {
					t.Errorf("expected 123, got %d", v)
				}
				return nil
			},
		},
		{
			name:  "Uint64",
			value: "123",
			check: func(e baseValueExtractor[TestValue]) error {
				v, err := e.Uint64()
				if err != nil {
					return err
				}
				if v != 123 {
					t.Errorf("expected 123, got %d", v)
				}
				return nil
			},
		},
		{
			name:  "Float32",
			value: "123.45",
			check: func(e baseValueExtractor[TestValue]) error {
				v, err := e.Float32()
				if err != nil {
					return err
				}
				if v != 123.45 {
					t.Errorf("expected 123.45, got %f", v)
				}
				return nil
			},
		},
		{
			name:  "Float64",
			value: "123.45",
			check: func(e baseValueExtractor[TestValue]) error {
				v, err := e.Float64()
				if err != nil {
					return err
				}
				if v != 123.45 {
					t.Errorf("expected 123.45, got %f", v)
				}
				return nil
			},
		},
		{
			name:  "Bool",
			value: "true",
			check: func(e baseValueExtractor[TestValue]) error {
				v, err := e.Bool()
				if err != nil {
					return err
				}
				if !v {
					t.Errorf("expected true, got %v", v)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractor := baseValueExtractor[TestValue]{value: TestValue(tt.value)}
			if err := tt.check(extractor); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
