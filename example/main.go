package main

import (
	"context"
	"log"
	"net/http"

	"github.com/eatmoreapple/hx"
)

// UserRequest represents the request structure
type UserRequest struct {
	ID   int    `json:"id" path:"id"`
	Name string `json:"name"`
}

// UserResponse represents the response structure
type UserResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// StringResponse represents a string response
type StringResponse string

// XMLResponse represents an XML response structure
type XMLResponse struct {
	Message string `xml:"message"`
}

func main() {
	// Create a new router
	r := hx.New()

	// JSON example
	jsonHandler := hx.G(func(ctx context.Context, req UserRequest) (UserResponse, error) {
		return UserResponse{
			ID:   req.ID,
			Name: req.Name,
		}, nil
	}).JSON()

	// String example
	stringHandler := hx.G(func(ctx context.Context, req UserRequest) (string, error) {
		return "Hello, " + req.Name, nil
	}).String()

	// XML example
	xmlHandler := hx.G(func(ctx context.Context, req UserRequest) (XMLResponse, error) {
		return XMLResponse{
			Message: "Hello, " + req.Name,
		}, nil
	}).XML()

	// Register routes
	r.GET("/user/{id}", jsonHandler)
	r.GET("/greet/{id}", stringHandler)
	r.GET("/xml/{id}", xmlHandler)

	// Start the server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}