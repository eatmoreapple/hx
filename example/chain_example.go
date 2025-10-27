package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/eatmoreapple/hx"
)

// UserRequest represents the request structure
type UserRequest struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// UserResponse represents the response structure
type UserResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	// Define middleware functions that only return errors
	authMiddleware := func(ctx context.Context, req UserRequest) error {
		// Simulate authentication check
		if req.ID <= 0 {
			return errors.New("unauthorized: invalid user ID")
		}
		fmt.Println("Authentication passed")
		return nil
	}

	validationMiddleware := func(ctx context.Context, req UserRequest) error {
		// Simulate validation check
		if req.Name == "" {
			return errors.New("validation failed: name is required")
		}
		fmt.Println("Validation passed")
		return nil
	}

	// Define the final handler that returns a response
	finalHandler := func(ctx context.Context, req UserRequest) (UserResponse, error) {
		fmt.Println("Processing request in final handler")
		// Simulate processing the request
		response := UserResponse{
			ID:        req.ID,
			Name:      req.Name,
			CreatedAt: time.Now(),
		}
		return response, nil
	}

	// Create a chain with middlewares and final handler
	handler := hx.Chain([]hx.MiddlewareFunc[UserRequest]{
		authMiddleware,
		validationMiddleware,
	}, finalHandler)

	// Convert to JSON HTTP handler
	httpHandler := handler.JSON()

	// Register the handler
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		// In a real application, you would bind the request to UserRequest
		// For this example, we'll create a sample request
		req := UserRequest{ID: 1, Name: "John Doe"}
		
		// Create a context
		ctx := r.Context()
		
		// Call the handler directly for demonstration
		resp, err := handler(ctx, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		
		// In a real application, you would use the hx framework's request handling
		// For now, we'll manually serialize to JSON
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"id":%d,"name":"%s","created_at":"%s"}`, resp.ID, resp.Name, resp.CreatedAt.Format(time.RFC3339))
	})

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}