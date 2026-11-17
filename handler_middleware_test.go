package hx

import (
	"context"
	"errors"
	"testing"
)

type chainRequest struct {
	Name string
}

type chainKey struct{}

func TestHandlerChainOrderAndRequest(t *testing.T) {
	var steps []string

	first := func(ctx context.Context, req chainRequest, next TypedHandlerFunc[chainRequest, string]) (string, error) {
		steps = append(steps, "first:"+req.Name)
		req.Name = req.Name + "-1"
		ctx = context.WithValue(ctx, chainKey{}, "from-first")
		resp, err := next(ctx, req)
		steps = append(steps, "first-after:"+resp)
		return resp + "/first", err
	}
	second := func(ctx context.Context, req chainRequest, next TypedHandlerFunc[chainRequest, string]) (string, error) {
		steps = append(steps, "second:"+req.Name+":"+ctx.Value(chainKey{}).(string))
		req.Name = req.Name + "-2"
		resp, err := next(ctx, req)
		steps = append(steps, "second-after:"+resp)
		return resp + "/second", err
	}

	handler := Pipe(first, second).Handle(func(ctx context.Context, req chainRequest) (string, error) {
		steps = append(steps, "handler:"+req.Name+":"+ctx.Value(chainKey{}).(string))
		return req.Name, nil
	})

	got, err := handler(context.Background(), chainRequest{Name: "ada"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got != "ada-1-2/second/first" {
		t.Fatalf("expected response %q, got %q", "ada-1-2/second/first", got)
	}

	expected := []string{
		"first:ada",
		"second:ada-1:from-first",
		"handler:ada-1-2:from-first",
		"second-after:ada-1-2",
		"first-after:ada-1-2/second",
	}
	if len(steps) != len(expected) {
		t.Fatalf("expected steps %v, got %v", expected, steps)
	}
	for i, step := range expected {
		if steps[i] != step {
			t.Fatalf("expected steps %v, got %v", expected, steps)
		}
	}
}

func TestHandlerChainStopsWhenNextIsSkipped(t *testing.T) {
	expectedErr := errors.New("stop")
	called := false

	handler := Pipe(
		func(ctx context.Context, req chainRequest, next TypedHandlerFunc[chainRequest, string]) (string, error) {
			return "", expectedErr
		},
	).Handle(func(ctx context.Context, req chainRequest) (string, error) {
		called = true
		return req.Name, nil
	})

	_, err := handler(context.Background(), chainRequest{Name: "ada"})
	if err != expectedErr {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
	if called {
		t.Fatal("expected handler to be skipped")
	}
}

func TestHandlerChainEmpty(t *testing.T) {
	handler := Pipe[chainRequest, string]().Handle(func(ctx context.Context, req chainRequest) (string, error) {
		return req.Name, nil
	})

	got, err := handler(context.Background(), chainRequest{Name: "ada"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got != "ada" {
		t.Fatalf("expected response %q, got %q", "ada", got)
	}
}

func TestPipeCopiesMiddlewares(t *testing.T) {
	called := false
	middlewares := []HandlerMiddleware[chainRequest, string]{
		func(ctx context.Context, req chainRequest, next TypedHandlerFunc[chainRequest, string]) (string, error) {
			called = true
			return next(ctx, req)
		},
	}
	chain := Pipe(middlewares...)
	middlewares[0] = func(ctx context.Context, req chainRequest, next TypedHandlerFunc[chainRequest, string]) (string, error) {
		t.Fatal("replacement middleware should not run")
		return next(ctx, req)
	}

	_, err := chain.Handle(func(ctx context.Context, req chainRequest) (string, error) {
		return req.Name, nil
	})(context.Background(), chainRequest{Name: "ada"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !called {
		t.Fatal("expected original middleware to run")
	}
}
