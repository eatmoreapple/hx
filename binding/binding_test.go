package binding

import (
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

type TestStruct struct {
	Name httpx.FromQuery[TestExtractor] `json:"name"`
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

func TestGenericBinder(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?test=hello", nil)
	var ts TestStruct

	binder := Generic()
	if err := binder.Bind(req, &ts); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if ts.Name.String() != "hello" {
		t.Errorf("expected name %s, got %s", "hello", ts.Name.String())
	}
}

func TestGenericBinderPointer(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?test=hello", nil)

	type TestStructPtr struct {
		Name *httpx.FromQuery[TestExtractor]
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
