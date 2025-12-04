package hx

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	r := New()

	r.GET("/hello", Warp(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	if w.Body.String() != "hello" {
		t.Errorf("expected body %s, got %s", "hello", w.Body.String())
	}
}

func TestRouterGroup(t *testing.T) {
	r := New()
	g := r.Group("/api")

	g.GET("/users", Warp(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("users"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	if w.Body.String() != "users" {
		t.Errorf("expected body %s, got %s", "users", w.Body.String())
	}
}

func TestRouterMiddleware(t *testing.T) {
	r := New()

	r.Use(func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			w.Header().Set("X-Test", "true")
			return next(w, r)
		}
	})

	r.GET("/", Warp(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Header().Get("X-Test") != "true" {
		t.Error("expected middleware to set header")
	}
}

func TestRouterErrorHandler(t *testing.T) {
	expectedErr := errors.New("oops")

	r := New(WithErrorHandler(func(w http.ResponseWriter, r *http.Request, err error) {
		if err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("custom error"))
	}))

	r.GET("/", func(w http.ResponseWriter, r *http.Request) error {
		return expectedErr
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	if w.Body.String() != "custom error" {
		t.Errorf("expected body %s, got %s", "custom error", w.Body.String())
	}
}

func TestRouterMethods(t *testing.T) {
	r := New()
	handler := Warp(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.POST("/post", handler)
	r.PUT("/put", handler)
	r.DELETE("/delete", handler)
	r.PATCH("/patch", handler)
	r.OPTIONS("/options", handler)
	r.HEAD("/head", handler)

	tests := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/post"},
		{http.MethodPut, "/put"},
		{http.MethodDelete, "/delete"},
		{http.MethodPatch, "/patch"},
		{http.MethodOptions, "/options"},
		{http.MethodHead, "/head"},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected status code %d for %s %s, got %d", http.StatusOK, tt.method, tt.path, w.Code)
		}
	}
}

func TestJoinPath(t *testing.T) {
	tests := []struct {
		a, b, expected string
	}{
		{"/", "path", "/path"},
		{"/api", "users", "/api/users"},
		{"/api/", "users", "/api/users"},
		{"/api", "/users", "/api/users"},
		{"/api/", "/users", "/api/users"},
	}

	for _, tt := range tests {
		if got := joinPath(tt.a, tt.b); got != tt.expected {
			t.Errorf("joinPath(%q, %q) = %q, expected %q", tt.a, tt.b, got, tt.expected)
		}
	}
}
