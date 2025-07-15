Advanced Features
=================

This section covers advanced HX features for building sophisticated applications.

Custom Request Extractors
--------------------------

You can create custom extractors by implementing the ``RequestExtractor`` interface:

.. code-block:: go

   package main

   import (
       "context"
       "fmt"
       "net/http"
       "strconv"
       "strings"

       "github.com/eatmoreapple/hx"
       . "github.com/eatmoreapple/hx/httpx"
   )

   // Custom extractor for pagination
   type PaginationExtractor struct {
       Page  int `json:"page"`
       Limit int `json:"limit"`
       Total int `json:"total"`
   }

   func (p *PaginationExtractor) FromRequest(r *http.Request) error {
       // Set defaults
       p.Page = 1
       p.Limit = 10

       // Extract page
       if pageStr := r.URL.Query().Get("page"); pageStr != "" {
           if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
               p.Page = page
           }
       }

       // Extract limit
       if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
           if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
               p.Limit = limit
           }
       }

       return nil
   }

   // Custom extractor for search filters
   type SearchFilters struct {
       Query    string   `json:"query"`
       Tags     []string `json:"tags"`
       Category string   `json:"category"`
       SortBy   string   `json:"sort_by"`
       SortDesc bool     `json:"sort_desc"`
   }

   func (s *SearchFilters) FromRequest(r *http.Request) error {
       query := r.URL.Query()

       s.Query = query.Get("q")
       s.Category = query.Get("category")
       s.SortBy = query.Get("sort_by")
       
       if s.SortBy == "" {
           s.SortBy = "created_at"
       }

       s.SortDesc = query.Get("order") == "desc"

       // Parse tags from comma-separated values
       if tagsStr := query.Get("tags"); tagsStr != "" {
           s.Tags = strings.Split(tagsStr, ",")
           // Trim whitespace
           for i, tag := range s.Tags {
               s.Tags[i] = strings.TrimSpace(tag)
           }
       }

       return nil
   }

   // Request type using custom extractors
   type SearchRequest struct {
       Pagination PaginationExtractor `json:"pagination"`
       Filters    SearchFilters       `json:"filters"`
   }

   func searchHandler(ctx context.Context, req SearchRequest) (map[string]interface{}, error) {
       return map[string]interface{}{
           "pagination": req.Pagination,
           "filters":    req.Filters,
           "results":    []string{"item1", "item2", "item3"}, // Mock results
       }, nil
   }

   func main() {
       router := hx.New()
       router.GET("/search", hx.G(searchHandler).JSON())

       fmt.Println("Server starting on :8080")
       fmt.Println("Try: http://localhost:8080/search?q=golang&tags=web,api&category=tutorial&page=2&limit=5&sort_by=title&order=desc")
       http.ListenAndServe(":8080", router)
   }

Custom Response Types
---------------------

Implement the ``ResponseRender`` interface for custom response handling:

.. code-block:: go

   package main

   import (
       "context"
       "encoding/csv"
       "fmt"
       "net/http"
       "strconv"

       "github.com/eatmoreapple/hx"
       . "github.com/eatmoreapple/hx/httpx"
   )

   // CSV Response
   type CSVResponse struct {
       Headers []string
       Rows    [][]string
   }

   func (c CSVResponse) IntoResponse(w http.ResponseWriter) error {
       w.Header().Set("Content-Type", "text/csv")
       w.Header().Set("Content-Disposition", "attachment; filename=data.csv")

       writer := csv.NewWriter(w)
       defer writer.Flush()

       // Write headers
       if err := writer.Write(c.Headers); err != nil {
           return err
       }

       // Write rows
       for _, row := range c.Rows {
           if err := writer.Write(row); err != nil {
               return err
           }
       }

       return nil
   }

   // PDF Response (simplified)
   type PDFResponse struct {
       Content []byte
       Filename string
   }

   func (p PDFResponse) IntoResponse(w http.ResponseWriter) error {
       w.Header().Set("Content-Type", "application/pdf")
       if p.Filename != "" {
           w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", p.Filename))
       }
       _, err := w.Write(p.Content)
       return err
   }

   // Template Response
   type TemplateResponse struct {
       Template string
       Data     interface{}
   }

   func (t TemplateResponse) IntoResponse(w http.ResponseWriter) error {
       w.Header().Set("Content-Type", "text/html")
       
       // Simple template rendering (use a real template engine in production)
       html := fmt.Sprintf(`
       <!DOCTYPE html>
       <html>
       <head><title>%s</title></head>
       <body>
           <h1>%s</h1>
           <pre>%+v</pre>
       </body>
       </html>`, t.Template, t.Template, t.Data)
       
       _, err := w.Write([]byte(html))
       return err
   }

   func csvHandler(ctx context.Context, req Empty) (CSVResponse, error) {
       return CSVResponse{
           Headers: []string{"ID", "Name", "Email"},
           Rows: [][]string{
               {"1", "John Doe", "john@example.com"},
               {"2", "Jane Smith", "jane@example.com"},
               {"3", "Bob Johnson", "bob@example.com"},
           },
       }, nil
   }

   func pdfHandler(ctx context.Context, req Empty) (PDFResponse, error) {
       // Mock PDF content
       content := []byte("%PDF-1.4\n1 0 obj\n<<\n/Type /Catalog\n/Pages 2 0 R\n>>\nendobj\n...")
       return PDFResponse{
           Content:  content,
           Filename: "report.pdf",
       }, nil
   }

   func templateHandler(ctx context.Context, req Empty) (TemplateResponse, error) {
       return TemplateResponse{
           Template: "User Dashboard",
           Data: map[string]interface{}{
               "User": "John Doe",
               "Time": "2025-01-01 12:00:00",
           },
       }, nil
   }

   func main() {
       router := hx.New()

       router.GET("/export/csv", hx.R(func(ctx context.Context, req Empty) (ResponseRender, error) {
           return csvHandler(ctx, req)
       }))

       router.GET("/export/pdf", hx.R(func(ctx context.Context, req Empty) (ResponseRender, error) {
           return pdfHandler(ctx, req)
       }))

       router.GET("/dashboard", hx.R(func(ctx context.Context, req Empty) (ResponseRender, error) {
           return templateHandler(ctx, req)
       }))

       fmt.Println("Server starting on :8080")
       http.ListenAndServe(":8080", router)
   }

Advanced Middleware Patterns
-----------------------------

Rate Limiting Middleware
~~~~~~~~~~~~~~~~~~~~~~~~

.. code-block:: go

   package main

   import (
       "fmt"
       "net/http"
       "sync"
       "time"

       "github.com/eatmoreapple/hx"
   )

   type RateLimiter struct {
       requests map[string][]time.Time
       mutex    sync.Mutex
       limit    int
       window   time.Duration
   }

   func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
       return &RateLimiter{
           requests: make(map[string][]time.Time),
           limit:    limit,
           window:   window,
       }
   }

   func (rl *RateLimiter) Middleware() hx.Middleware {
       return func(next hx.HandlerFunc) hx.HandlerFunc {
           return func(w http.ResponseWriter, r *http.Request) error {
               clientIP := r.RemoteAddr
               
               rl.mutex.Lock()
               defer rl.mutex.Unlock()

               now := time.Now()
               
               // Clean old requests
               if requests, exists := rl.requests[clientIP]; exists {
                   filtered := requests[:0]
                   for _, reqTime := range requests {
                       if now.Sub(reqTime) < rl.window {
                           filtered = append(filtered, reqTime)
                       }
                   }
                   rl.requests[clientIP] = filtered
               }

               // Check limit
               if len(rl.requests[clientIP]) >= rl.limit {
                   return fmt.Errorf("rate limit exceeded")
               }

               // Add current request
               rl.requests[clientIP] = append(rl.requests[clientIP], now)

               return next(w, r)
           }
       }
   }

Circuit Breaker Middleware
~~~~~~~~~~~~~~~~~~~~~~~~~~

.. code-block:: go

   package main

   import (
       "fmt"
       "net/http"
       "sync"
       "time"

       "github.com/eatmoreapple/hx"
   )

   type CircuitState int

   const (
       StateClosed CircuitState = iota
       StateOpen
       StateHalfOpen
   )

   type CircuitBreaker struct {
       maxFailures  int
       resetTimeout time.Duration
       state        CircuitState
       failures     int
       lastFailTime time.Time
       mutex        sync.Mutex
   }

   func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
       return &CircuitBreaker{
           maxFailures:  maxFailures,
           resetTimeout: resetTimeout,
           state:        StateClosed,
       }
   }

   func (cb *CircuitBreaker) Middleware() hx.Middleware {
       return func(next hx.HandlerFunc) hx.HandlerFunc {
           return func(w http.ResponseWriter, r *http.Request) error {
               cb.mutex.Lock()
               
               // Check if we should reset
               if cb.state == StateOpen && time.Since(cb.lastFailTime) > cb.resetTimeout {
                   cb.state = StateHalfOpen
                   cb.failures = 0
               }

               // Reject if circuit is open
               if cb.state == StateOpen {
                   cb.mutex.Unlock()
                   return fmt.Errorf("circuit breaker is open")
               }

               cb.mutex.Unlock()

               // Execute request
               err := next(w, r)

               cb.mutex.Lock()
               defer cb.mutex.Unlock()

               if err != nil {
                   cb.failures++
                   cb.lastFailTime = time.Now()
                   
                   if cb.failures >= cb.maxFailures {
                       cb.state = StateOpen
                   }
               } else if cb.state == StateHalfOpen {
                   cb.state = StateClosed
                   cb.failures = 0
               }

               return err
           }
       }
   }

Request Context Enhancement
---------------------------

.. code-block:: go

   package main

   import (
       "context"
       "fmt"
       "net/http"
       "time"

       "github.com/eatmoreapple/hx"
       . "github.com/eatmoreapple/hx/httpx"
   )

   // Context keys
   type contextKey string

   const (
       requestIDKey contextKey = "request_id"
       userIDKey    contextKey = "user_id"
       traceIDKey   contextKey = "trace_id"
   )

   // Request ID middleware
   func requestIDMiddleware(next hx.HandlerFunc) hx.HandlerFunc {
       return func(w http.ResponseWriter, r *http.Request) error {
           requestID := r.Header.Get("X-Request-ID")
           if requestID == "" {
               requestID = fmt.Sprintf("%d", time.Now().UnixNano())
           }

           ctx := context.WithValue(r.Context(), requestIDKey, requestID)
           r = r.WithContext(ctx)

           w.Header().Set("X-Request-ID", requestID)

           return next(w, r)
       }
   }

   // User context middleware
   func userContextMiddleware(next hx.HandlerFunc) hx.HandlerFunc {
       return func(w http.ResponseWriter, r *http.Request) error {
           userID := r.Header.Get("X-User-ID")
           if userID != "" {
               ctx := context.WithValue(r.Context(), userIDKey, userID)
               r = r.WithContext(ctx)
           }

           return next(w, r)
       }
   }

   func contextHandler(ctx context.Context, req Empty) (map[string]interface{}, error) {
       response := make(map[string]interface{})

       if requestID := ctx.Value(requestIDKey); requestID != nil {
           response["request_id"] = requestID
       }

       if userID := ctx.Value(userIDKey); userID != nil {
           response["user_id"] = userID
       }

       response["message"] = "Context data extracted successfully"

       return response, nil
   }

   func main() {
       router := hx.New()

       router.Use(requestIDMiddleware, userContextMiddleware)
       router.GET("/context", hx.E(contextHandler).JSON())

       fmt.Println("Server starting on :8080")
       http.ListenAndServe(":8080", router)
   }

Request Validation
------------------

.. code-block:: go

   package main

   import (
       "context"
       "fmt"
       "net/http"
       "regexp"
       "strings"

       "github.com/eatmoreapple/hx"
       . "github.com/eatmoreapple/hx/httpx"
   )

   // Validation interface
   type Validator interface {
       Validate() error
   }

   // Validation middleware
   func validationMiddleware(next hx.HandlerFunc) hx.HandlerFunc {
       return func(w http.ResponseWriter, r *http.Request) error {
           // This middleware would need to be applied at the handler level
           // for access to the typed request
           return next(w, r)
       }
   }

   // User creation request with validation
   type CreateUserRequest struct {
       Name     string `json:"name" form:"name"`
       Email    string `json:"email" form:"email"`
       Password string `json:"password" form:"password"`
       Age      int    `json:"age" form:"age"`
   }

   func (r CreateUserRequest) Validate() error {
       var errors []string

       // Name validation
       if r.Name == "" {
           errors = append(errors, "name is required")
       } else if len(r.Name) < 2 {
           errors = append(errors, "name must be at least 2 characters")
       }

       // Email validation
       if r.Email == "" {
           errors = append(errors, "email is required")
       } else {
           emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
           if !emailRegex.MatchString(r.Email) {
               errors = append(errors, "invalid email format")
           }
       }

       // Password validation
       if r.Password == "" {
           errors = append(errors, "password is required")
       } else if len(r.Password) < 8 {
           errors = append(errors, "password must be at least 8 characters")
       }

       // Age validation
       if r.Age < 0 || r.Age > 150 {
           errors = append(errors, "age must be between 0 and 150")
       }

       if len(errors) > 0 {
           return fmt.Errorf("validation errors: %s", strings.Join(errors, ", "))
       }

       return nil
   }

   func createUserHandler(ctx context.Context, req CreateUserRequest) (map[string]interface{}, error) {
       // Validate request
       if err := req.Validate(); err != nil {
           return nil, err
       }

       // Process valid request
       return map[string]interface{}{
           "message": "User created successfully",
           "user": map[string]interface{}{
               "name":  req.Name,
               "email": req.Email,
               "age":   req.Age,
           },
       }, nil
   }

   func main() {
       router := hx.New()

       router.POST("/users", hx.G(createUserHandler).JSON())

       fmt.Println("Server starting on :8080")
       fmt.Println("Test with: curl -X POST http://localhost:8080/users -H 'Content-Type: application/json' -d '{\"name\":\"John\",\"email\":\"john@example.com\",\"password\":\"password123\",\"age\":25}'")
       http.ListenAndServe(":8080", router)
   }

Database Integration
--------------------

.. code-block:: go

   package main

   import (
       "context"
       "database/sql"
       "fmt"
       "net/http"

       "github.com/eatmoreapple/hx"
       . "github.com/eatmoreapple/hx/httpx"
       _ "github.com/mattn/go-sqlite3" // SQLite driver
   )

   type User struct {
       ID    int    `json:"id" db:"id"`
       Name  string `json:"name" db:"name"`
       Email string `json:"email" db:"email"`
   }

   type UserService struct {
       db *sql.DB
   }

   func NewUserService(db *sql.DB) *UserService {
       return &UserService{db: db}
   }

   func (s *UserService) GetUser(id int) (*User, error) {
       user := &User{}
       err := s.db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id).
           Scan(&user.ID, &user.Name, &user.Email)
       if err != nil {
           return nil, err
       }
       return user, nil
   }

   func (s *UserService) CreateUser(name, email string) (*User, error) {
       result, err := s.db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", name, email)
       if err != nil {
           return nil, err
       }

       id, err := result.LastInsertId()
       if err != nil {
           return nil, err
       }

       return &User{
           ID:    int(id),
           Name:  name,
           Email: email,
       }, nil
   }

   // Dependency injection middleware
   func serviceMiddleware(userService *UserService) hx.Middleware {
       return func(next hx.HandlerFunc) hx.HandlerFunc {
           return func(w http.ResponseWriter, r *http.Request) error {
               ctx := context.WithValue(r.Context(), "userService", userService)
               r = r.WithContext(ctx)
               return next(w, r)
           }
       }
   }

   type UserIDExtractor string
   func (u UserIDExtractor) ValueName() string { return "id" }

   type GetUserRequest struct {
       ID FromPath[UserIDExtractor] `json:"id"`
   }

   type CreateUserRequest struct {
       Name  string `json:"name" form:"name"`
       Email string `json:"email" form:"email"`
   }

   func getUserHandler(ctx context.Context, req GetUserRequest) (*User, error) {
       userService := ctx.Value("userService").(*UserService)
       
       id := 0 // Convert string to int (simplified)
       fmt.Sscanf(string(req.ID), "%d", &id)
       
       return userService.GetUser(id)
   }

   func createUserHandler(ctx context.Context, req CreateUserRequest) (*User, error) {
       userService := ctx.Value("userService").(*UserService)
       return userService.CreateUser(req.Name, req.Email)
   }

   func main() {
       // Initialize database
       db, err := sql.Open("sqlite3", ":memory:")
       if err != nil {
           panic(err)
       }
       defer db.Close()

       // Create table
       _, err = db.Exec(`
           CREATE TABLE users (
               id INTEGER PRIMARY KEY AUTOINCREMENT,
               name TEXT NOT NULL,
               email TEXT NOT NULL UNIQUE
           )
       `)
       if err != nil {
           panic(err)
       }

       userService := NewUserService(db)

       router := hx.New()
       router.Use(serviceMiddleware(userService))

       router.GET("/users/{id}", hx.G(getUserHandler).JSON())
       router.POST("/users", hx.G(createUserHandler).JSON())

       fmt.Println("Server starting on :8080")
       http.ListenAndServe(":8080", router)
   }