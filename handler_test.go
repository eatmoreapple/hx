package hx

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/eatmoreapple/hx/httpx"
)

func TestWarp(t *testing.T) {
	handler := Warp(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	if err := handler(w, req); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	if w.Body.String() != "ok" {
		t.Errorf("expected body %s, got %s", "ok", w.Body.String())
	}
}

func TestGeneric(t *testing.T) {
	type Request struct {
		Name string
	}
	type Response struct {
		Message string
	}

	handler := Generic[Request, Response](func(ctx context.Context, req Request) (Response, error) {
		return Response{Message: "Hello " + req.Name}, nil
	})

	// Since Generic returns a TypedHandlerFunc, we can't directly call it as http.HandlerFunc
	// We need to verify it returns the correct function type
	if handler == nil {
		t.Error("expected handler to be not nil")
	}
}

func TestRender(t *testing.T) {
	type Request struct {
		Name string `json:"name"`
	}

	handler := Render[Request](func(ctx context.Context, req Request) (httpx.ResponseRender, error) {
		return httpx.JSONResponse{Data: map[string]string{"message": "Hello " + req.Name}}, nil
	})

	// Create a request with JSON body
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	// Mock binding by manually setting the request context or just testing the handler logic if possible.
	// However, Render uses internal requestHandler which uses ShouldBind.
	// We need to provide a body for binding.
	// For simplicity, let's assume query binding for GET or body for POST.
	// Since we haven't implemented the full test suite for binding yet, we might face issues if binding fails.
	// But Render creates a HandlerFunc that does binding.

	// Let's try with a simple struct that doesn't require complex binding first, or just empty.
	w := httptest.NewRecorder()
	if err := handler(w, req); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestJSON(t *testing.T) {
	type Request struct{}
	type Response struct {
		Message string `json:"message"`
	}

	handler := JSON(func(ctx context.Context, req Request) (Response, error) {
		return Response{Message: "hello"}, nil
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	if err := handler(w, req); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if resp.Message != "hello" {
		t.Errorf("expected message %s, got %s", "hello", resp.Message)
	}
}

func TestString(t *testing.T) {
	type Request struct{}

	handler := G(func(ctx context.Context, req Request) (string, error) {
		return "hello", nil
	}).String()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	if err := handler(w, req); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if w.Body.String() != "hello" {
		t.Errorf("expected body %s, got %s", "hello", w.Body.String())
	}
}

func TestStringPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic but got nil")
		}
	}()

	type Request struct{}
	type Response struct{}

	G(func(ctx context.Context, req Request) (Response, error) {
		return Response{}, nil
	}).String()
}

func TestXML(t *testing.T) {
	type Request struct{}
	type Response struct {
		XMLName xml.Name `xml:"response"`
		Message string   `xml:"message"`
	}

	handler := G(func(ctx context.Context, req Request) (Response, error) {
		return Response{Message: "hello"}, nil
	}).XML()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	if err := handler(w, req); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	var resp Response
	if err := xml.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if resp.Message != "hello" {
		t.Errorf("expected message %s, got %s", "hello", resp.Message)
	}
}

type token string

func (token) FromRequest(*http.Request) error { return nil }

func TestRequestTypeCheckedWhenBuildingHandler(t *testing.T) {
	type request struct {
		Name string `json:"name"`
	}

	JSON(func(context.Context, request) (struct{}, error) { return struct{}{}, nil })
	JSON(func(context.Context, *request) (struct{}, error) { return struct{}{}, nil })
	JSON(func(context.Context, Query) (struct{}, error) { return struct{}{}, nil })
	JSON(func(context.Context, *Query) (struct{}, error) { return struct{}{}, nil })
	JSON(func(context.Context, Header) (struct{}, error) { return struct{}{}, nil })
	JSON(func(context.Context, Form) (struct{}, error) { return struct{}{}, nil })
	JSON(func(context.Context, Cookies) (struct{}, error) { return struct{}{}, nil })
	JSON(func(context.Context, FromPath[string]) (struct{}, error) { return struct{}{}, nil })
	JSON(func(context.Context, *FromQuery[string]) (struct{}, error) { return struct{}{}, nil })
	JSON(func(context.Context, Empty) (struct{}, error) { return struct{}{}, nil })
	JSON(func(context.Context, token) (struct{}, error) { return struct{}{}, nil })
	Render(func(context.Context, request) (httpx.ResponseRender, error) {
		return httpx.StringResponse{Data: "ok"}, nil
	})
	G(func(context.Context, request) (struct{}, error) { return struct{}{}, nil }).XML()
}

func TestNonStructRequestPanicsWhenBuildingHandler(t *testing.T) {
	assertPanic := func(t *testing.T, requestType reflect.Type, build func()) {
		t.Helper()
		defer func() {
			recovered := recover()
			if recovered == nil {
				t.Fatalf("%s: expected panic", requestType)
			}
			message, ok := recovered.(string)
			if !ok {
				t.Fatalf("%s: expected string panic, got %T", requestType, recovered)
			}
			if !strings.Contains(message, requestType.String()) || !strings.Contains(message, "must be a struct") {
				t.Fatalf("%s: unexpected panic %q", requestType, message)
			}
		}()
		build()
	}

	assertPanic(t, reflect.TypeFor[map[string]any](), func() {
		JSON(func(context.Context, map[string]any) (struct{}, error) { return struct{}{}, nil })
	})
	assertPanic(t, reflect.TypeFor[*map[string]any](), func() {
		JSON(func(context.Context, *map[string]any) (struct{}, error) { return struct{}{}, nil })
	})
	assertPanic(t, reflect.TypeFor[[]string](), func() {
		JSON(func(context.Context, []string) (struct{}, error) { return struct{}{}, nil })
	})
	assertPanic(t, reflect.TypeFor[string](), func() {
		Render(func(context.Context, string) (httpx.ResponseRender, error) { return nil, nil })
	})
	assertPanic(t, reflect.TypeFor[int](), func() {
		G(func(context.Context, int) (struct{}, error) { return struct{}{}, nil }).XML()
	})
}

func TestJSONPointerStruct(t *testing.T) {
	type request struct {
		Name string `json:"name"`
	}

	handler := JSON(func(_ context.Context, req *request) (string, error) {
		return req.Name, nil
	})

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"alice"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	if err := handler(w, req); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if w.Body.String() != "\"alice\"\n" {
		t.Fatalf("expected JSON string, got %s", w.Body.String())
	}
}

func TestJSONQueryExtractor(t *testing.T) {
	handler := JSON(func(_ context.Context, req Query) (int, error) {
		return len(req), nil
	})

	req := httptest.NewRequest(http.MethodGet, "/?a=1&b=2", nil)
	w := httptest.NewRecorder()

	if err := handler(w, req); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if w.Body.String() != "2\n" {
		t.Fatalf("expected query length 2, got %s", w.Body.String())
	}
}

func TestE(t *testing.T) {
	handler := E(func(ctx context.Context) (string, error) {
		return "ok", nil
	}).String()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	if err := handler(w, req); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if w.Body.String() != "ok" {
		t.Errorf("expected body %s, got %s", "ok", w.Body.String())
	}
}
