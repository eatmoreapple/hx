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

	// Create a handler using Pipe with slice of middlewares
	handler1 := hx.Pipe(
		[]hx.MiddlewareFunc[UserRequest]{authMiddleware, validationMiddleware},
		finalHandler,
	)

	// Create a handler using PipeFunc with variadic middlewares
	handler2 := hx.PipeFunc(
		finalHandler,
		authMiddleware,
		validationMiddleware,
	)

	// Both handlers work the same way
	_ = handler1
	_ = handler2

	fmt.Println("Example compiled successfully")
}