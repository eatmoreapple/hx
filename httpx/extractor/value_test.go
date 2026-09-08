package extractor

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type TestValue string

func (t TestValue) ValueName() string {
	return "test"
}

func TestBaseValueExtractor(t *testing.T) {
	extractor := baseValueExtractor[string]{raw: "123", value: "123"}

	if extractor.Value() != "123" {
		t.Errorf("expected value %s, got %s", "123", extractor.Value())
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

func TestQueryValueExtractorFromRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/?test=from-request", nil)
	var extractor QueryValueExtractor[TestValue]

	if err := extractor.FromRequest(request); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := extractor.String(); got != "from-request" {
		t.Fatalf("expected %q, got %q", "from-request", got)
	}
	if got := extractor.Value(); got != TestValue("from-request") {
		t.Fatalf("expected value %q, got %q", "from-request", got)
	}
}

type namedPage int

func (namedPage) ValueName() string { return "page" }

func TestQueryValueExtractorFromRequestDefinedInt(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/?page=7", nil)
	var extractor QueryValueExtractor[namedPage]

	if err := extractor.FromRequest(request); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := extractor.Value(); got != 7 {
		t.Fatalf("expected 7, got %d", got)
	}
}

func TestQueryValueExtractorFromRequestRequiresName(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/?query=value", nil)
	var extractor QueryValueExtractor[string]

	if err := extractor.FromRequest(request); !errors.Is(err, ErrValueNameRequired) {
		t.Fatalf("expected ErrValueNameRequired, got %v", err)
	}
}

func TestQueryValueExtractorFromRequestInt(t *testing.T) {
	type requestFields struct {
		Page QueryValueExtractor[int] `hx:"page"`
	}

	field, _ := reflect.TypeFor[requestFields]().FieldByName("Page")
	request := httptest.NewRequest(http.MethodGet, "/?page=42", nil)
	var extractor QueryValueExtractor[int]

	if err := extractor.FromRequestField(request, field); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := extractor.Value(); got != 42 {
		t.Fatalf("expected 42, got %d", got)
	}
	if got := extractor.String(); got != "42" {
		t.Fatalf("expected raw %q, got %q", "42", got)
	}
}

func TestQueryValueExtractorFromRequestIntEmpty(t *testing.T) {
	type requestFields struct {
		Page QueryValueExtractor[int] `hx:"page"`
	}

	field, _ := reflect.TypeFor[requestFields]().FieldByName("Page")
	var extractor QueryValueExtractor[int]

	for _, rawURL := range []string{"/", "/?page=", "/?other=1"} {
		request := httptest.NewRequest(http.MethodGet, rawURL, nil)
		if err := extractor.FromRequestField(request, field); err != nil {
			t.Fatalf("%s: unexpected error: %v", rawURL, err)
		}
		if got := extractor.Value(); got != 0 {
			t.Fatalf("%s: expected 0, got %d", rawURL, got)
		}
	}
}

func TestQueryValueExtractorFromRequestIntRequiredEmpty(t *testing.T) {
	type requestFields struct {
		Page QueryValueExtractor[int] `hx:"page" required:"true"`
	}

	field, _ := reflect.TypeFor[requestFields]().FieldByName("Page")
	var extractor QueryValueExtractor[int]
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	if err := extractor.FromRequestField(request, field); !errors.Is(err, ErrEmptyValue) {
		t.Fatalf("expected ErrEmptyValue, got %v", err)
	}
}

func TestFormValueExtractorFromRequestIntEmpty(t *testing.T) {
	type requestFields struct {
		Age FormValueExtractor[int] `hx:"age"`
	}

	field, _ := reflect.TypeFor[requestFields]().FieldByName("Age")
	request := httptest.NewRequest(http.MethodPost, "/", nil)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var extractor FormValueExtractor[int]

	if err := extractor.FromRequestField(request, field); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := extractor.Value(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestFormValueExtractorFromRequestIntRequiredEmpty(t *testing.T) {
	type requestFields struct {
		Age FormValueExtractor[int] `hx:"age" required:"true"`
	}

	field, _ := reflect.TypeFor[requestFields]().FieldByName("Age")
	request := httptest.NewRequest(http.MethodPost, "/", nil)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var extractor FormValueExtractor[int]

	if err := extractor.FromRequestField(request, field); !errors.Is(err, ErrEmptyValue) {
		t.Fatalf("expected ErrEmptyValue, got %v", err)
	}
}

func TestHeaderValueExtractorFromRequestIntEmpty(t *testing.T) {
	type requestFields struct {
		Count HeaderValueExtractor[int] `hx:"X-Count"`
	}

	field, _ := reflect.TypeFor[requestFields]().FieldByName("Count")
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	var extractor HeaderValueExtractor[int]

	if err := extractor.FromRequestField(request, field); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := extractor.Value(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestHeaderValueExtractorFromRequestIntRequiredEmpty(t *testing.T) {
	type requestFields struct {
		Count HeaderValueExtractor[int] `hx:"X-Count" required:"true"`
	}

	field, _ := reflect.TypeFor[requestFields]().FieldByName("Count")
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	var extractor HeaderValueExtractor[int]

	if err := extractor.FromRequestField(request, field); !errors.Is(err, ErrEmptyValue) {
		t.Fatalf("expected ErrEmptyValue, got %v", err)
	}
}

func TestCookieValueExtractorFromRequestIntEmpty(t *testing.T) {
	type requestFields struct {
		ID CookieValueExtractor[int] `hx:"id"`
	}

	field, _ := reflect.TypeFor[requestFields]().FieldByName("ID")
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	var extractor CookieValueExtractor[int]

	if err := extractor.FromRequestField(request, field); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := extractor.Value(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestCookieValueExtractorFromRequestIntRequiredEmpty(t *testing.T) {
	type requestFields struct {
		ID CookieValueExtractor[int] `hx:"id" required:"true"`
	}

	field, _ := reflect.TypeFor[requestFields]().FieldByName("ID")
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	var extractor CookieValueExtractor[int]

	if err := extractor.FromRequestField(request, field); err == nil {
		t.Fatal("expected error for missing required cookie")
	}
}

func TestPathValueExtractorFromRequestIntEmpty(t *testing.T) {
	type requestFields struct {
		ID PathValueExtractor[int] `hx:"id"`
	}

	field, _ := reflect.TypeFor[requestFields]().FieldByName("ID")
	request := httptest.NewRequest(http.MethodGet, "/users/", nil)
	var extractor PathValueExtractor[int]

	if err := extractor.FromRequestField(request, field); !errors.Is(err, ErrEmptyValue) {
		t.Fatalf("expected ErrEmptyValue, got %v", err)
	}
}

func TestQueryValueExtractorFromRequestIntInvalid(t *testing.T) {
	type requestFields struct {
		Page QueryValueExtractor[int] `hx:"page"`
	}

	field, _ := reflect.TypeFor[requestFields]().FieldByName("Page")
	request := httptest.NewRequest(http.MethodGet, "/?page=not-a-number", nil)
	var extractor QueryValueExtractor[int]

	if err := extractor.FromRequestField(request, field); err == nil {
		t.Fatal("expected conversion error")
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

func TestValueExtractorsIgnoreXML(t *testing.T) {
	type requestBody struct {
		Header HeaderValueExtractor[string] `xml:"header,attr"`
		Query  QueryValueExtractor[string]  `xml:"query"`
		Name   string                       `xml:"name"`
	}

	var got requestBody
	got.Header.raw = "original-header"
	got.Query.raw = "original-query"

	data := []byte(`<request header="ignored"><query><nested>ignored</nested></query><name>hello</name></request>`)
	if err := xml.Unmarshal(data, &got); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Header.String() != "original-header" {
		t.Fatalf("expected header extractor to be ignored, got %q", got.Header.String())
	}
	if got.Query.String() != "original-query" {
		t.Fatalf("expected query extractor to be ignored, got %q", got.Query.String())
	}
	if got.Name != "hello" {
		t.Fatalf("expected regular XML field %q, got %q", "hello", got.Name)
	}
}
