// Package httpx provides HTTP response handling utilities and interfaces.
package httpx

import (
	"cmp"
	"encoding/json"
	"encoding/xml"
	"html/template"
	"io"
	"net/http"
)

// ResponseRender defines the interface for types that can render themselves as HTTP responses.
// Implementations should handle setting appropriate headers and writing response data.
type ResponseRender interface {
	IntoResponse(http.ResponseWriter) error
}

// JSONResponse represents a JSON response with data and status code.
// It automatically sets the Content-Type header to application/json.
type JSONResponse struct {
	Data       any // Data to be encoded as JSON
	StatusCode int // HTTP status code (defaults to 200 OK if not set)
}

// IntoResponse implements ResponseRender for JSON responses.
// It sets the appropriate content type, status code, and encodes the data as JSON.
func (j JSONResponse) IntoResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(cmp.Or(j.StatusCode, http.StatusOK))
	return json.NewEncoder(w).Encode(j.Data)
}

// XMLResponse represents an XML response with data and status code.
// It automatically sets the Content-Type header to application/xml.
type XMLResponse struct {
	Data       any // Data to be encoded as XML
	StatusCode int // HTTP status code (defaults to 200 OK if not set)
}

// IntoResponse implements ResponseRender for XML responses.
// It sets the appropriate content type, status code, and encodes the data as XML.
func (x XMLResponse) IntoResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(cmp.Or(x.StatusCode, http.StatusOK))
	return xml.NewEncoder(w).Encode(x.Data)
}

// StringResponse represents a plain text response with string data and status code.
// It automatically sets the Content-Type header to text/plain.
type StringResponse struct {
	Data       string // String data to be sent in response
	StatusCode int    // HTTP status code (defaults to 200 OK if not set)
}

// IntoResponse implements ResponseRender for string responses.
// It sets the appropriate content type, status code, and writes the string data.
func (s StringResponse) IntoResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(cmp.Or(s.StatusCode, http.StatusOK))
	_, err := io.WriteString(w, s.Data)
	return err
}

// HTMLResponse represents an HTML response with template, data, and status code.
// It automatically sets the Content-Type header to text/html.
type HTMLResponse struct {
	Data       any                // Data to be passed to the template
	StatusCode int                // HTTP status code (defaults to 200 OK if not set)
	Template   *template.Template // Template to be executed
}

// IntoResponse implements ResponseRender for HTML responses.
// It sets the appropriate content type, status code, and executes the template with the provided data.
func (h HTMLResponse) IntoResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(cmp.Or(h.StatusCode, http.StatusOK))
	return h.Template.Execute(w, h.Data)
}
