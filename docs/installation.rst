Installation
============

Requirements
------------

HX requires Go 1.24 or higher.

Installing HX
-------------

To install HX, use ``go get``:

.. code-block:: bash

   go get github.com/eatmoreapple/hx

Import HX in your Go code:

.. code-block:: go

   import "github.com/eatmoreapple/hx"

Verifying Installation
----------------------

Create a simple test file to verify your installation:

.. code-block:: go

   package main

   import (
       "context"
       "fmt"
       "net/http"

       "github.com/eatmoreapple/hx"
       . "github.com/eatmoreapple/hx/httpx"
   )

   func hello(ctx context.Context, req Empty) (string, error) {
       return "Hello, HX!", nil
   }

   func main() {
       router := hx.New()
       router.GET("/hello", hx.G(hello).String())
       
       fmt.Println("Server starting on :8080")
       http.ListenAndServe(":8080", router)
   }

Run the file and visit http://localhost:8080/hello to see if HX is working correctly.