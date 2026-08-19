package binding

import (
	"bytes"
	"errors"
	"mime/multipart"
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

func TestFormBinderMultipart(t *testing.T) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("name", "hello"); err != nil {
		t.Fatal(err)
	}
	file, err := writer.CreateFormFile("avatar", "avatar.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.Write([]byte("content")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	type upload struct {
		Name   string                `form:"name"`
		Avatar *multipart.FileHeader `form:"avatar"`
	}
	req := httptest.NewRequest(http.MethodPost, "/?name=query", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var got upload
	if err := formBinder.Bind(req, &got); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "hello" {
		t.Fatalf("expected form value %q, got %q", "hello", got.Name)
	}
	if got.Avatar == nil || got.Avatar.Filename != "avatar.txt" {
		t.Fatalf("unexpected uploaded file: %#v", got.Avatar)
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

func TestGenericBinderAnonymousStruct(t *testing.T) {
	type Embedded struct {
		Query httpx.FromQuery[string]
	}
	type Request struct {
		Embedded
	}

	req := httptest.NewRequest(http.MethodGet, "/?Query=hello", nil)
	var got Request

	if err := Generic().Bind(req, &got); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Embedded.Query.String() != "hello" {
		t.Fatalf("expected embedded query %q, got %q", "hello", got.Embedded.Query.String())
	}
}

func TestGenericBinderRequiresStructPointer(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?Query=hello", nil)
	var value struct {
		Query httpx.FromQuery[string]
	}
	var nilValue *struct {
		Query httpx.FromQuery[string]
	}

	for name, target := range map[string]any{
		"value":              value,
		"nil pointer":        nilValue,
		"non-struct pointer": new(string),
	} {
		t.Run(name, func(t *testing.T) {
			if err := Generic().Bind(req, target); !errors.Is(err, ErrGenericBinderTarget) {
				t.Fatalf("expected ErrGenericBinderTarget, got %v", err)
			}
		})
	}
}

func TestGenericBinderSkipsUnexportedFields(t *testing.T) {
	type embedded struct {
		Query httpx.FromQuery[string]
	}
	type request struct {
		embedded
	}

	var got request
	if err := Generic().Bind(httptest.NewRequest(http.MethodGet, "/?Query=hello", nil), &got); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.embedded.Query.String() != "" {
		t.Fatalf("expected unexported embedded field to be skipped, got %q", got.embedded.Query.String())
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

func TestGenericBinderPreservesInitializedPointer(t *testing.T) {
	type TestStructPtr struct {
		Name *httpx.FromQuery[string]
	}
	existing := &httpx.FromQuery[string]{}
	got := TestStructPtr{Name: existing}

	if err := Generic().Bind(httptest.NewRequest(http.MethodGet, "/?Name=hello", nil), &got); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != existing {
		t.Fatal("expected initialized extractor pointer to be preserved")
	}
	if got.Name.String() != "hello" {
		t.Fatalf("expected name %q, got %q", "hello", got.Name.String())
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

func TestGenericBinderExtractorSpecificTags(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"/?query-name=from-query",
		strings.NewReader("form-name=from-form"),
	)
	req.Header.Set("Content-Type", MIMEPOSTForm)
	req.Header.Set("X-Token", "from-header")
	req.SetPathValue("path-id", "42")
	req.AddCookie(&http.Cookie{Name: "session-id", Value: "from-cookie"})

	type taggedExtractors struct {
		Path   httpx.FromPath[string]   `path:"path-id"`
		Query  httpx.FromQuery[string]  `query:"query-name"`
		Header httpx.FromHeader[string] `header:"X-Token"`
		Form   httpx.FromForm[string]   `form:"form-name"`
		Cookie httpx.FromCookie[string] `cookie:"session-id"`
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

func TestGenericBinderValueNameTagPrecedence(t *testing.T) {
	type taggedExtractor struct {
		Query httpx.FromQuery[string] `hx:"preferred" query:"fallback"`
	}

	req := httptest.NewRequest(http.MethodGet, "/?preferred=from-hx&fallback=from-query-tag", nil)
	var got taggedExtractor
	if err := Generic().Bind(req, &got); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Query.String() != "from-hx" {
		t.Fatalf("expected hx value %q, got %q", "from-hx", got.Query.String())
	}
}

func TestGenericBinderFieldNameFallback(t *testing.T) {
	type fieldNameValue struct {
		Query httpx.FromQuery[string]
	}

	var got fieldNameValue
	err := Generic().Bind(httptest.NewRequest(http.MethodGet, "/?Query=value", nil), &got)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Query.String() != "value" {
		t.Fatalf("expected field-name value %q, got %q", "value", got.Query.String())
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
