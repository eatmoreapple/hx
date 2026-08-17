package binding

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/eatmoreapple/hx/httpx"
)

type TestExtractor string

func (t TestExtractor) ValueName() string {
	return "test"
}

type EmptyNameExtractor string

func (EmptyNameExtractor) ValueName() string {
	return ""
}

type FormValue string

func (v *FormValue) UnmarshalForm(values []string) error {
	*v = FormValue(strings.Join(values, ","))
	return nil
}

func valueNameOf[T httpx.NamedValue](value T) string {
	return value.ValueName()
}

type TestStruct struct {
	Name httpx.FromQuery[TestExtractor] `json:"name" hx:"ignored"`
}

func TestDefault(t *testing.T) {
	tests := []struct {
		method      string
		contentType string
		expected    Binder
	}{
		{http.MethodGet, "application/json", queryBinder},
		{http.MethodPost, "application/json", jsonBinder},
		{http.MethodPost, "application/xml", xmlBinder},
		{http.MethodPost, "application/x-www-form-urlencoded", formBinder},
		{http.MethodPost, "multipart/form-data", formBinder},
		{http.MethodPost, "text/plain", queryBinder},
		{http.MethodPost, "invalid", queryBinder},
	}

	for _, tt := range tests {
		binder := Default(tt.method, tt.contentType)
		if binder != tt.expected {
			t.Errorf("expected binder %T, got %T", tt.expected, binder)
		}
	}
}

func TestFormUnmarshaler(t *testing.T) {
	type formValues struct {
		Custom  FormValue `form:"custom"`
		Default string    `form:"default"`
	}

	var got formValues
	err := bindValues(map[string][]string{
		"custom":  {"first", "second"},
		"default": {"value"},
	}, &got)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Custom != "first,second" {
		t.Fatalf("expected custom value %q, got %q", "first,second", got.Custom)
	}
	if got.Default != "value" {
		t.Fatalf("expected default value %q, got %q", "value", got.Default)
	}
}

func TestFormBinderIgnoresValueExtractor(t *testing.T) {
	type requestValues struct {
		Name  string                  `form:"name"`
		Query httpx.FromQuery[string] `form:"query" hx:"query"`
	}

	request := httptest.NewRequest(
		http.MethodPost,
		"/?query=from-query",
		strings.NewReader("name=hello&query=from-form"),
	)
	request.Header.Set("Content-Type", MIMEPOSTForm)

	var got requestValues
	if err := formBinder.Bind(request, &got); err != nil {
		t.Fatalf("unexpected form binding error: %v", err)
	}
	if got.Query.String() != "" {
		t.Fatalf("expected form binder to ignore extractor, got %q", got.Query.String())
	}
	if err := Generic().Bind(request, &got); err != nil {
		t.Fatalf("unexpected generic binding error: %v", err)
	}
	if got.Name != "hello" {
		t.Fatalf("expected form value %q, got %q", "hello", got.Name)
	}
	if got.Query.String() != "from-query" {
		t.Fatalf("expected query extractor value %q, got %q", "from-query", got.Query.String())
	}
}

func TestGenericBinder(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?test=hello&ignored=wrong", nil)
	var ts TestStruct

	binder := Generic()
	if err := binder.Bind(req, &ts); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if ts.Name.String() != "hello" {
		t.Errorf("expected name %s, got %s", "hello", ts.Name.String())
	}
}

func TestNamedValueConstraint(t *testing.T) {
	if got := valueNameOf(TestExtractor("value")); got != "test" {
		t.Fatalf("expected value name %q, got %q", "test", got)
	}
}

func TestGenericBinderPointer(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?test=hello", nil)

	type TestStructPtr struct {
		Name *httpx.FromQuery[string] `hx:"test"`
	}
	var ts TestStructPtr

	binder := Generic()
	if err := binder.Bind(req, &ts); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if ts.Name == nil {
		t.Fatal("expected name to be not nil")
	}

	if ts.Name.String() != "hello" {
		t.Errorf("expected name %s, got %s", "hello", ts.Name.String())
	}
}

func TestGenericBinderValueNameTag(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?query=from-query&form=from-form", nil)
	req.SetPathValue("id", "42")
	req.Header.Set("X-Token", "from-header")
	req.AddCookie(&http.Cookie{Name: "session", Value: "from-cookie"})

	type queryValue string
	type taggedExtractors struct {
		Path   httpx.FromPath[string]      `hx:"id"`
		Query  httpx.FromQuery[queryValue] `hx:"query"`
		Header httpx.FromHeader[string]    `hx:"X-Token"`
		Form   httpx.FromForm[string]      `hx:"form"`
		Cookie httpx.FromCookie[string]    `hx:"session"`
	}

	var got taggedExtractors
	if err := Generic().Bind(req, &got); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Path.String() != "42" {
		t.Errorf("expected path value %q, got %q", "42", got.Path.String())
	}
	if got.Query.String() != "from-query" {
		t.Errorf("expected query value %q, got %q", "from-query", got.Query.String())
	}
	if got.Header.String() != "from-header" {
		t.Errorf("expected header value %q, got %q", "from-header", got.Header.String())
	}
	if got.Form.String() != "from-form" {
		t.Errorf("expected form value %q, got %q", "from-form", got.Form.String())
	}
	if got.Cookie.String() != "from-cookie" {
		t.Errorf("expected cookie value %q, got %q", "from-cookie", got.Cookie.String())
	}
}

func TestGenericBinderValueNameRequired(t *testing.T) {
	type missingValueName struct {
		Query httpx.FromQuery[string]
	}

	var target missingValueName
	err := Generic().Bind(httptest.NewRequest(http.MethodGet, "/", nil), &target)
	if !errors.Is(err, httpx.ErrValueNameRequired) {
		t.Fatalf("expected ErrValueNameRequired, got %v", err)
	}
	if !strings.Contains(err.Error(), `field "Query"`) {
		t.Fatalf("expected error to identify Query field, got %v", err)
	}
}

func TestGenericBinderEmptyValueName(t *testing.T) {
	type emptyValueName struct {
		Query httpx.FromQuery[EmptyNameExtractor] `hx:"fallback"`
	}

	var target emptyValueName
	err := Generic().Bind(httptest.NewRequest(http.MethodGet, "/?fallback=value", nil), &target)
	if !errors.Is(err, httpx.ErrValueNameRequired) {
		t.Fatalf("expected ErrValueNameRequired, got %v", err)
	}
}

func TestJSONBinder(t *testing.T) {
	body := `{"name": "hello"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	type Data struct {
		Name string `json:"name"`
	}
	var data Data

	if err := jsonBinder.Bind(req, &data); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if data.Name != "hello" {
		t.Errorf("expected name %s, got %s", "hello", data.Name)
	}
}
