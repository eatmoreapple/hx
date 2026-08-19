package extractor

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestFormExtractorFromRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("name=hx"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var extractor FormExtractor
	if err := extractor.FromRequest(request); err != nil {
		t.Fatalf("FromRequest returned an unexpected error: %v", err)
	}
	if got := url.Values(extractor).Get("name"); got != "hx" {
		t.Fatalf("form value = %q, want %q", got, "hx")
	}
}

func TestFormExtractorFromMultipartRequest(t *testing.T) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("name", "hx"); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodPost, "/", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	var extractor FormExtractor
	if err := extractor.FromRequest(request); err != nil {
		t.Fatalf("FromRequest returned an unexpected error: %v", err)
	}
	if got := url.Values(extractor).Get("name"); got != "hx" {
		t.Fatalf("form value = %q, want %q", got, "hx")
	}
}

func TestFormExtractorReturnsParseError(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("name=%zz"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var extractor FormExtractor
	if err := extractor.FromRequest(request); err == nil {
		t.Fatal("FromRequest returned nil, want a form parsing error")
	}
}
