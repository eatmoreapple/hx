# HX - Lightweight HTTP Framework for Go

HX is a lightweight, flexible HTTP framework for Go that simplifies request handling and data extraction. It provides a clean, type-safe way to handle HTTP requests with minimal boilerplate.

## Features

- 🚀 Lightweight and fast
- 💪 Type-safe request data extraction
- 🔄 Automatic request binding
- 🛠 Extensible design

## Installation

require go version 1.24+

```bash
go get github.com/eatmoreapple/hx
```


## Quick Start

Here's a simple example that demonstrates the core features of HX:

```go
package main

import (
	"context"
	"net/http"

	"github.com/eatmoreapple/hx"
	"github.com/eatmoreapple/hx/httpx"
)

type Router string

func (r Router) ValueName() string {
	return "id"
}

type Ua string

func (u Ua) ValueName() string {
	return "user-agent"
}

type User struct {
	Name string                 `json:"name" form:"name"` // extract from request query
	Id   httpx.FromPath[Router] `json:"id"`               // extract from request path
	Ua   httpx.FromHeader[Ua]   `json:"ua"`               // extract from request header
}

func app(ctx context.Context, extractor User) (any, error) {
	return extractor, nil
}

func main() {
	router := hx.New()
	router.GET("/{id}", hx.G(app).JSON())

	http.ListenAndServe(":9999", router)
}
```

Open your browser and navigate to [http://localhost:9999/1?name=eatmoreapple](http://localhost:9999/1?name=eatmoreapple)

You should see a response like:

```json
{
  "name": "eatmoreapple",
  "id": "1",
  "ua": "Mozilla/5.0 ..."
}
```


## License

HX is released under the [Apache License 2.0](https://github.com/eatmoreapple/hx/blob/main/LICENSE).