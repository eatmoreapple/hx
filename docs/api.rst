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
* ``Static(pathPrefix string, root fs.FS)`` - Serves files from ``root`` under ``pathPrefix``

**Example:**

.. code-block:: go

   router := hx.New()
   router.GET("/users", handler)
   router.POST("/users", handler)

   // Route groups
   api := router.Group("/api/v1")
   api.GET("/users", handler)  // Maps to /api/v1/users

   // /assets/js/main.js serves ./public/assets/js/main.js
   router.Static("/assets", os.DirFS("./public/assets"))

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

``JSON()``, ``String()``, ``XML()``, and ``Render`` check the request type
while the handler is built. It must be a struct, a pointer to a struct, or a
type that implements ``FromRequest``. Any other type panics immediately:

.. code-block:: text

   hx: request type ... must be a struct, a pointer to a struct, or implement FromRequest

The panic happens when the handler is built, usually while the route is
registered, rather than on the first request. ``G`` itself only names the
handler; the check runs when that handler is turned into a ``HandlerFunc``.

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

JSON
~~~~

.. code-block:: go

   func JSON[Request, Response any](h TypedHandlerFunc[Request, Response]) HandlerFunc

Converts a typed handler into a handler that serializes its response as JSON.
Request and response types are inferred from the handler.

**Example:**

.. code-block:: go

   router.GET("/user/{id}", hx.JSON(userHandler))

Render / R
~~~~~~~~~~

.. code-block:: go

   func Render[Request any](h TypedHandlerFunc[Request, httpx.ResponseRender]) HandlerFunc
   func R[Request any](h TypedHandlerFunc[Request, httpx.ResponseRender]) HandlerFunc

Creates a handler that returns a ResponseRender for custom response handling.

E
~

.. code-block:: go

   func E[Response any](h func(ctx context.Context) (Response, error)) TypedHandlerFunc[hx.Empty, Response]

Convenience function for handlers that don't require request data.

**Example:**

.. code-block:: go

   func healthCheck(ctx context.Context) (string, error) {
       return "OK", nil
   }
   
   router.GET("/health", hx.E(healthCheck).String())

Request Extraction
------------------

The root ``hx`` package provides shortcut aliases for the request extractor
types in ``httpx``. Applications can use ``hx.FromPath``, ``hx.FromQuery``,
``hx.FromHeader``, ``hx.FromForm``, and ``hx.FromCookie`` directly.

FromPath
~~~~~~~~

.. code-block:: go

   type FromPath[T Value]

Extracts values from URL path parameters.

**Example:**

.. code-block:: go

   type Request struct {
       ID FromPath[int] `hx:"id" json:"id"`
   }
   
   // For route "/users/{id}", extracts the {id} value as an int.
   // Use Value() for the parsed result and String() for the raw path segment.

FromQuery
~~~~~~~~~

.. code-block:: go

   type FromQuery[T Value]

Extracts values from URL query parameters.

**Example:**

.. code-block:: go

   type Request struct {
       Page FromQuery[int] `hx:"page"`
   }

FromHeader
~~~~~~~~~~

.. code-block:: go

   type FromHeader[T Value]

Extracts values from HTTP headers.

FromForm
~~~~~~~~

.. code-block:: go

   type FromForm[T Value]

Extracts values from form data.

FromCookie
~~~~~~~~~~

.. code-block:: go

   type FromCookie[T Value]

Extracts values from HTTP cookies.

Empty and required values
~~~~~~~~~~~~~~~~~~~~~~~~~

As a struct field, ``FromQuery``, ``FromForm``, ``FromHeader``, and
``FromCookie`` treat a missing or empty value as the zero value of ``T``.
Tag the field ``required:"true"`` to return ``extractor.ErrEmptyValue``
instead.

``FromPath`` is always required. The ``required`` tag does not change path
parameters.

Calling ``FromRequest`` directly, outside struct-field binding, also rejects
an empty value.

.. code-block:: go

   type Request struct {
       Page FromQuery[int]    `hx:"page"`              // missing page is 0
       Q    FromQuery[string] `hx:"q" required:"true"` // missing q is ErrEmptyValue
       ID   FromPath[int]     `hx:"id"`                // path parameters are always required
   }

Value Interface
~~~~~~~~~~~~~~~

.. code-block:: go

   type Value interface {
       ~int | ~int8 | ~int16 | ~int32 | ~int64 |
           ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
           ~float32 | ~float64 |
           ~string |
           ~bool |
           uuid.UUID
   }

``T`` is parsed from the raw request string. Predeclared types such as
``int`` and ``string`` are converted directly; defined types
(``type UserID int``) use the same underlying kind so they can also
implement ``ValueNamer``.

When an extractor is a struct field, its request value name is resolved from
the first non-empty source in this order: ``hx`` tag, extractor-specific tag
(``path``, ``query``, ``header``, ``form``, or ``cookie``), then the Go field
name:

.. code-block:: go

   type Request struct {
       Search FromQuery[string] `query:"q"`
       Token  FromHeader[string] `hx:"X-Token" header:"ignored"`
       Page   FromQuery[int] // uses "Page"
   }

ValueNamer Interface
~~~~~~~~~~~~~~~~~~~~

.. code-block:: go

   type ValueNamer interface {
       ValueName() string
   }

Optional interface for reusable named value types. Only defined types can
implement it (predeclared ``int`` and ``string`` cannot have methods).
A non-empty ``ValueName`` takes precedence over all struct-field metadata.
Direct extraction outside a struct requires ``ValueName`` and otherwise
returns ``ErrValueNameRequired``.

.. code-block:: go

   type UserID int

   func (UserID) ValueName() string { return "id" }

   type Request struct {
       ID FromPath[UserID]
   }

``NamedValue`` combines ``Value`` and ``ValueNamer`` for generic code that
requires a value to provide its own name.

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

HX has two middleware layers.

Router middleware wraps a ``HandlerFunc``. It runs before the request is bound
and can see ``http.ResponseWriter`` and ``*http.Request``. Register it with
``Router.Use`` or ``Chain``.

Handler middleware wraps a typed handler. It runs after binding and can see
the typed request and response. Build it with ``Pipe`` and finish it with
``Handle``.

Middleware Type
~~~~~~~~~~~~~~~

.. code-block:: go

   type Middleware func(HandlerFunc) HandlerFunc

Function type for router middleware. It wraps a ``HandlerFunc`` and runs
before the request is bound.

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

HandlerMiddleware
~~~~~~~~~~~~~~~~~

.. code-block:: go

   type HandlerMiddleware[Request, Response any] func(ctx context.Context, req Request, next TypedHandlerFunc[Request, Response]) (Response, error)

Middleware for a typed handler. It runs after the request is bound.

Call ``next`` with the context and request that the rest of the chain should
see. Skip ``next`` to stop the chain. The inner middleware and the handler do
not run, and the returned response becomes the result. Middleware outside this
one still continues after ``next`` returns.

``JSON()``, ``String()``, and ``XML()`` write the response only when the error
is nil. A middleware that returns both a response and an error drops the
response, and the router handles the error.

HandlerChain
~~~~~~~~~~~~

.. code-block:: go

   type HandlerChain[Request, Response any] []HandlerMiddleware[Request, Response]

An ordered list of handler middleware. The first middleware runs first, and
its code after ``next`` runs last.

Pipe
~~~~

.. code-block:: go

   func Pipe[Request, Response any](middlewares ...HandlerMiddleware[Request, Response]) HandlerChain[Request, Response]

Returns a chain that runs ``middlewares`` from left to right. The returned
chain is a copy, so later changes to the caller's slice do not affect it.
An empty call returns a nil chain.

**Example:**

.. code-block:: go

   func trimName(ctx context.Context, req CreateUser, next hx.TypedHandlerFunc[CreateUser, User]) (User, error) {
       req.Name = strings.TrimSpace(req.Name)
       return next(ctx, req)
   }

   func audit(ctx context.Context, req CreateUser, next hx.TypedHandlerFunc[CreateUser, User]) (User, error) {
       user, err := next(ctx, req)
       log.Printf("created %s: %v", req.Name, err)
       return user, err
   }

   router.POST("/users", hx.Pipe(trimName, audit).Handle(createUser).JSON())

``trimName`` runs first, then ``audit``, then ``createUser``. The same chain
can wrap more than one handler:

.. code-block:: go

   users := hx.Pipe(trimName, audit)
   router.POST("/users", users.Handle(createUser).JSON())
   router.PUT("/users/{id}", users.Handle(updateUser).JSON())

Handle
~~~~~~

.. code-block:: go

   func (c HandlerChain[Request, Response]) Handle(h TypedHandlerFunc[Request, Response]) TypedHandlerFunc[Request, Response]

Ends the chain with ``h``. Middleware runs first, then ``h``. An empty chain
returns ``h`` unchanged. Call ``JSON()``, ``String()``, or ``XML()`` on the
result to register it as a route.

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

Form Value Decoding
~~~~~~~~~~~~~~~~~~~

The ``encoding/form`` package provides the form and query value decoder used
by ``FormBinder`` and ``QueryBinder``:

.. code-block:: go

   func Unmarshal(values url.Values, dest any) error

Multipart file headers can be decoded separately:

.. code-block:: go

   func UnmarshalFiles(files map[string][]*multipart.FileHeader, dest any) error

Fields can override the default scalar or slice conversion by implementing:

.. code-block:: go

   type Unmarshaler interface {
       UnmarshalForm([]string) error
   }

Error Handling
--------------

ErrorHandler
~~~~~~~~~~~~

.. code-block:: go

   type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

Function type for handling errors returned by handlers.

ExtractError
~~~~~~~~~~~~

.. code-block:: go

   type ExtractError struct {
       Field string
       Err   error
   }

Failures from ``FromRequest`` and ``FromRequestField`` are returned as
``*httpx.ExtractError``. ``Field`` is the struct field name when a field is
bound. When the request type itself implements ``FromRequest``, ``Field`` is
that type's name.

``Unwrap`` returns the original error, so ``errors.Is`` and ``errors.As``
still match ``extractor.ErrEmptyValue`` and ``httpx.ErrValueNameRequired``.
An error that already unwraps to ``ExtractError`` is not wrapped again.

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
