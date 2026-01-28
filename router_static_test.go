package hx

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestRouterStatic(t *testing.T) {
	// Create a temporary directory for static files
	tmpDir := t.TempDir()

	// Create a test file
	testFile := "test.txt"
	content := []byte("hello static world")
	if err := os.WriteFile(filepath.Join(tmpDir, testFile), content, 0644); err != nil {
		t.Fatal(err)
	}

	r := New()
	// Serve files from tmpDir under /static path using os.DirFS
	r.Static("/static", os.DirFS(tmpDir))

	// Test request
	req := httptest.NewRequest(http.MethodGet, "/static/"+testFile, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	if w.Body.String() != string(content) {
		t.Errorf("expected body %s, got %s", string(content), w.Body.String())
	}

	// Test nested static route
	g := r.Group("/api")
	g.Static("/assets", os.DirFS(tmpDir))

	req = httptest.NewRequest(http.MethodGet, "/api/assets/"+testFile, nil)
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	if w.Body.String() != string(content) {
		t.Errorf("expected body %s, got %s", string(content), w.Body.String())
	}
}
