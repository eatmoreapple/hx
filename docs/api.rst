API Reference
=============

This section provides detailed documentation for all HX components.

Core Types
----------

Router
~~~~~~

.. code-block:: go

   type Router struct {
       ErrHandler ErrorHandler
       // ... other fields
   }

The main router that handles HTTP request routing and middleware.

**Methods:**

* ``New(options ...RouterOption) *Router`` - Creates a new router
* ``Group(prefix string) *Router`` - Creates a route group with path prefix
* ``Use(middleware ...Middleware)`` - Adds middleware to the router
* ``Handle(method, path string, handler HandlerFunc)`` - Registers a route
* ``GET/POST/PUT/DELETE/PATCH/OPTIONS/HEAD(path string, handler HandlerFunc)`` - HTTP method shortcuts

**Example:**

.. code-block:: go

   router := hx.New()
   router.GET("/users", handler)
   router.POST("/users", handler)
   
   // Route groups
   api := router.Group("/api/v1")
   api.GET("/users", handler)  // Maps to /api/v1/users

HandlerFunc
~~~~~~~~~~~

.. code-block:: go

   type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

Standard handler function type that returns an error instead of void.

TypedHandlerFunc
~~~~~~~~~~~~~~~~

.. code-block:: go

   type TypedHandlerFunc[Request, Response any] func(context.Context, Request) (Response, error)

Generic handler function with type-safe request and response handling.

**Methods:**

* ``JSON() HandlerFunc`` - Converts to JSON response handler
* ``String() HandlerFunc`` - Converts to string response handler (Response must be string)
* ``XML() HandlerFunc`` - Converts to XML response handler

Handler Creation Functions
--------------------------

Generic / G
~~~~~~~~~~~

.. code-block:: go

   func Generic[Request, Response any](h TypedHandlerFunc[Request, Response]) TypedHandlerFunc[Request, Response]
   func G[Request, Response any](h TypedHandlerFunc[Request, Response]) TypedHandlerFunc[Request, Response]

Creates a type-safe handler with specified Request and Response types.

**Example:**

.. code-block:: go

   func userHandler(ctx context.Context, req UserRequest) (UserResponse, error) {
       // implementation
   }
   
   router.GET("/user/{id}", hx.G(userHandler).JSON())

Render / R
~~~~~~~~~~

.. code-block:: go

   func Render[Request any](h TypedHandlerFunc[Request, httpx.ResponseRender]) HandlerFunc
   func R[Request any](h TypedHandlerFunc[Request, httpx.ResponseRender]) HandlerFunc

Creates a handler that returns a ResponseRender for custom response handling.

E
~

.. code-block:: go

   func E[Response any](h func(ctx context.Context) (Response, error)) TypedHandlerFunc[httpx.Empty, Response]

Convenience function for handlers that don't require request data.

**Example:**

.. code-block:: go

   func healthCheck(ctx context.Context) (string, error) {
       return "OK", nil
   }
   
   router.GET("/health", hx.E(healthCheck).String())

Request Extraction
------------------

The ``httpx`` package provides types for extracting data from different parts of HTTP requests.

FromPath
~~~~~~~~

.. code-block:: go

   type FromPath[T ValueNamer] T

Extracts values from URL path parameters.

**Example:**

.. code-block:: go

   type UserID string
   func (u UserID) ValueName() string { return "id" }
   
   type Request struct {
       ID FromPath[UserID] `json:"id"`
   }
   
   // For route "/user/{id}", extracts the {id} value

FromQuery
~~~~~~~~~

.. code-block:: go

   type FromQuery[T ValueNamer] T

Extracts values from URL query parameters.

FromHeader
~~~~~~~~~~

.. code-block:: go

   type FromHeader[T ValueNamer] T

Extracts values from HTTP headers.

FromForm
~~~~~~~~

.. code-block:: go

   type FromForm[T ValueNamer] T

Extracts values from form data.

FromCookie
~~~~~~~~~~

.. code-block:: go

   type FromCookie[T ValueNamer] T

Extracts values from HTTP cookies.

ValueNamer Interface
~~~~~~~~~~~~~~~~~~~~

.. code-block:: go

   type ValueNamer interface {
       ValueName() string
   }

Interface that extraction types must implement to specify the field name.

Response Types
--------------

ResponseRender Interface
~~~~~~~~~~~~~~~~~~~~~~~~

.. code-block:: go

   type ResponseRender interface {
       IntoResponse(w http.ResponseWriter) error
   }

Interface for custom response rendering.

JSONResponse
~~~~~~~~~~~~

.. code-block:: go

   type JSONResponse struct {
       Data any
   }

Renders response as JSON.

StringResponse
~~~~~~~~~~~~~~

.. code-block:: go

   type StringResponse struct {
       Data string
   }

Renders response as plain text.

XMLResponse
~~~~~~~~~~~

.. code-block:: go

   type XMLResponse struct {
       Data any
   }

Renders response as XML.

Middleware
----------

Middleware Type
~~~~~~~~~~~~~~~

.. code-block:: go

   type Middleware func(HandlerFunc) HandlerFunc

Function type for middleware that wraps handlers.

Chain
~~~~~

.. code-block:: go

   func Chain(middleware ...Middleware) Middleware

Chains multiple middleware functions together.

**Example:**

.. code-block:: go

   func loggingMiddleware(next hx.HandlerFunc) hx.HandlerFunc {
       return func(w http.ResponseWriter, r *http.Request) error {
           log.Printf("%s %s", r.Method, r.URL.Path)
           return next(w, r)
       }
   }
   
   router.Use(loggingMiddleware)

Router Options
--------------

WithErrorHandler
~~~~~~~~~~~~~~~~

.. code-block:: go

   func WithErrorHandler(handler ErrorHandler) RouterOption

Sets a custom error handler for the router.

WithMiddleware
~~~~~~~~~~~~~~

.. code-block:: go

   func WithMiddleware(middleware ...Middleware) RouterOption

Adds middleware to the router during creation.

**Example:**

.. code-block:: go

   router := hx.New(
       hx.WithErrorHandler(customErrorHandler),
       hx.WithMiddleware(loggingMiddleware, authMiddleware),
   )

Binding
-------

ShouldBind
~~~~~~~~~~

.. code-block:: go

   func ShouldBind(r *http.Request, e any) error

Binds request data to the given interface using appropriate binders based on Content-Type.

Binder Interface
~~~~~~~~~~~~~~~~

.. code-block:: go

   type Binder interface {
       Bind(*http.Request, any) error
   }

Interface for request data binding implementations.

Available binders:

* ``JSONBinder`` - Binds JSON request bodies
* ``XMLBinder`` - Binds XML request bodies  
* ``FormBinder`` - Binds form data (multipart and URL-encoded)
* ``QueryBinder`` - Binds URL query parameters

Error Handling
--------------

ErrorHandler
~~~~~~~~~~~~

.. code-block:: go

   type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

Function type for handling errors returned by handlers.

**Default Error Handler:**

.. code-block:: go

   func defaultErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
       http.Error(w, err.Error(), http.StatusInternalServerError)
   }

Utilities
---------

Warp
~~~~

.. code-block:: go

   func Warp(h http.HandlerFunc) HandlerFunc

Wraps a standard ``http.HandlerFunc`` into HX's ``HandlerFunc``.

**Example:**

.. code-block:: go

   standardHandler := func(w http.ResponseWriter, r *http.Request) {
       w.Write([]byte("Hello"))
   }
   
   router.GET("/hello", hx.Warp(standardHandler))