Quick Start
===========

This guide will help you get started with HX quickly. You'll learn the basic concepts and create your first HX application.

Basic Concepts
--------------

HX is built around several key concepts:

* **Type-safe handlers** - Functions that accept typed request data and return typed responses
* **Request extraction** - Automatic binding of request data to Go structs
* **Response rendering** - Automatic serialization of responses to various formats
* **Middleware** - Composable request processing pipeline

Your First HX Application
--------------------------

Let's create a simple REST API that demonstrates HX's core features:

.. code-block:: go

   package main

   import (
       "context"
       "net/http"

       "github.com/eatmoreapple/hx"
       . "github.com/eatmoreapple/hx/httpx"
   )

   // Define custom extractors for path and header values
   type UserID string

   func (u UserID) ValueName() string {
       return "id"
   }

   type UserAgent string

   func (u UserAgent) ValueName() string {
       return "user-agent"
   }

   // Request structure with automatic data extraction
   type UserRequest struct {
       Name string                `json:"name" form:"name"`     // from query/form
       ID   FromPath[UserID]      `json:"id"`                   // from URL path
       UA   FromHeader[UserAgent] `json:"user_agent"`           // from headers
   }

   // Response structure
   type UserResponse struct {
       ID        string `json:"id"`
       Name      string `json:"name"`
       UserAgent string `json:"user_agent"`
       Message   string `json:"message"`
   }

   // Handler function with type safety
   func getUserInfo(ctx context.Context, req UserRequest) (UserResponse, error) {
       return UserResponse{
           ID:        string(req.ID),
           Name:      req.Name,
           UserAgent: string(req.UA),
           Message:   "Hello from HX!",
       }, nil
   }

   func main() {
       router := hx.New()
       
       // Register a JSON endpoint
       router.GET("/user/{id}", hx.G(getUserInfo).JSON())
       
       http.ListenAndServe(":8080", router)
   }

Test the application by visiting:
http://localhost:8080/user/123?name=john

You should see a JSON response like:

.. code-block:: json

   {
     "id": "123",
     "name": "john",
     "user_agent": "Mozilla/5.0...",
     "message": "Hello from HX!"
   }

Request Data Extraction
-----------------------

HX provides several ways to extract data from HTTP requests:

Query Parameters
~~~~~~~~~~~~~~~~

.. code-block:: go

   type QueryRequest struct {
       Page  int    `form:"page"`
       Limit int    `form:"limit"`
       Query string `form:"q"`
   }

Form Data
~~~~~~~~~

.. code-block:: go

   type FormRequest struct {
       Username string `form:"username"`
       Password string `form:"password"`
   }

JSON Body
~~~~~~~~~

.. code-block:: go

   type JSONRequest struct {
       Name  string `json:"name"`
       Email string `json:"email"`
   }

Path Parameters
~~~~~~~~~~~~~~~

.. code-block:: go

   type IDExtractor string
   func (i IDExtractor) ValueName() string { return "id" }

   type PathRequest struct {
       ID FromPath[IDExtractor] `json:"id"`
   }

Headers
~~~~~~~

.. code-block:: go

   type AuthExtractor string
   func (a AuthExtractor) ValueName() string { return "authorization" }

   type HeaderRequest struct {
       Auth FromHeader[AuthExtractor] `json:"auth"`
   }

Response Formats
----------------

HX supports multiple response formats:

JSON Response
~~~~~~~~~~~~~

.. code-block:: go

   func jsonHandler(ctx context.Context, req Empty) (map[string]interface{}, error) {
       return map[string]interface{}{
           "message": "Hello, World!",
           "status":  "success",
       }, nil
   }

   router.GET("/json", hx.G(jsonHandler).JSON())

String Response
~~~~~~~~~~~~~~~

.. code-block:: go

   func stringHandler(ctx context.Context, req Empty) (string, error) {
       return "Hello, World!", nil
   }

   router.GET("/text", hx.G(stringHandler).String())

XML Response
~~~~~~~~~~~~

.. code-block:: go

   type XMLResponse struct {
       Message string `xml:"message"`
       Status  string `xml:"status"`
   }

   func xmlHandler(ctx context.Context, req Empty) (XMLResponse, error) {
       return XMLResponse{
           Message: "Hello, World!",
           Status:  "success",
       }, nil
   }

   router.GET("/xml", hx.G(xmlHandler).XML())

Error Handling
--------------

HX provides built-in error handling. Simply return an error from your handler:

.. code-block:: go

   func errorHandler(ctx context.Context, req Empty) (string, error) {
       return "", fmt.Errorf("something went wrong")
   }

   router.GET("/error", hx.G(errorHandler).String())

You can also customize error handling:

.. code-block:: go

   func customErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
       http.Error(w, "Custom error: "+err.Error(), http.StatusInternalServerError)
   }

   router := hx.New(hx.WithErrorHandler(customErrorHandler))

Next Steps
----------

Now that you understand the basics, explore these topics:

* :doc:`api` - Detailed API reference
* :doc:`examples` - More examples and use cases
* :doc:`advanced` - Advanced features like middleware and custom extractors