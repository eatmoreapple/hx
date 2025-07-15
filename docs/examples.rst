Examples
========

This section provides practical examples of using HX in various scenarios.

REST API Example
----------------

Here's a complete example of a REST API for managing users:

.. code-block:: go

   package main

   import (
       "context"
       "fmt"
       "net/http"
       "strconv"
       "sync"

       "github.com/eatmoreapple/hx"
       . "github.com/eatmoreapple/hx/httpx"
   )

   // User model
   type User struct {
       ID    int    `json:"id"`
       Name  string `json:"name"`
       Email string `json:"email"`
   }

   // In-memory storage
   var (
       users   = make(map[int]User)
       usersMu = sync.RWMutex{}
       nextID  = 1
   )

   // Extractors
   type UserID string
   func (u UserID) ValueName() string { return "id" }

   // Request types
   type GetUserRequest struct {
       ID FromPath[UserID] `json:"id"`
   }

   type CreateUserRequest struct {
       Name  string `json:"name" form:"name"`
       Email string `json:"email" form:"email"`
   }

   type UpdateUserRequest struct {
       ID    FromPath[UserID] `json:"id"`
       Name  string           `json:"name" form:"name"`
       Email string           `json:"email" form:"email"`
   }

   type DeleteUserRequest struct {
       ID FromPath[UserID] `json:"id"`
   }

   // Handlers
   func listUsers(ctx context.Context, req Empty) ([]User, error) {
       usersMu.RLock()
       defer usersMu.RUnlock()

       userList := make([]User, 0, len(users))
       for _, user := range users {
           userList = append(userList, user)
       }
       return userList, nil
   }

   func getUser(ctx context.Context, req GetUserRequest) (User, error) {
       id, err := strconv.Atoi(string(req.ID))
       if err != nil {
           return User{}, fmt.Errorf("invalid user ID: %v", err)
       }

       usersMu.RLock()
       defer usersMu.RUnlock()

       user, exists := users[id]
       if !exists {
           return User{}, fmt.Errorf("user not found")
       }
       return user, nil
   }

   func createUser(ctx context.Context, req CreateUserRequest) (User, error) {
       if req.Name == "" || req.Email == "" {
           return User{}, fmt.Errorf("name and email are required")
       }

       usersMu.Lock()
       defer usersMu.Unlock()

       user := User{
           ID:    nextID,
           Name:  req.Name,
           Email: req.Email,
       }
       users[nextID] = user
       nextID++

       return user, nil
   }

   func updateUser(ctx context.Context, req UpdateUserRequest) (User, error) {
       id, err := strconv.Atoi(string(req.ID))
       if err != nil {
           return User{}, fmt.Errorf("invalid user ID: %v", err)
       }

       usersMu.Lock()
       defer usersMu.Unlock()

       user, exists := users[id]
       if !exists {
           return User{}, fmt.Errorf("user not found")
       }

       if req.Name != "" {
           user.Name = req.Name
       }
       if req.Email != "" {
           user.Email = req.Email
       }

       users[id] = user
       return user, nil
   }

   func deleteUser(ctx context.Context, req DeleteUserRequest) (map[string]string, error) {
       id, err := strconv.Atoi(string(req.ID))
       if err != nil {
           return nil, fmt.Errorf("invalid user ID: %v", err)
       }

       usersMu.Lock()
       defer usersMu.Unlock()

       if _, exists := users[id]; !exists {
           return nil, fmt.Errorf("user not found")
       }

       delete(users, id)
       return map[string]string{"message": "user deleted successfully"}, nil
   }

   func main() {
       router := hx.New()

       // Routes
       router.GET("/users", hx.E(listUsers).JSON())
       router.GET("/users/{id}", hx.G(getUser).JSON())
       router.POST("/users", hx.G(createUser).JSON())
       router.PUT("/users/{id}", hx.G(updateUser).JSON())
       router.DELETE("/users/{id}", hx.G(deleteUser).JSON())

       fmt.Println("Server starting on :8080")
       http.ListenAndServe(":8080", router)
   }

Test the API:

.. code-block:: bash

   # List users
   curl http://localhost:8080/users

   # Create user
   curl -X POST http://localhost:8080/users \
        -H "Content-Type: application/json" \
        -d '{"name":"John Doe","email":"john@example.com"}'

   # Get user
   curl http://localhost:8080/users/1

   # Update user
   curl -X PUT http://localhost:8080/users/1 \
        -H "Content-Type: application/json" \
        -d '{"name":"Jane Doe"}'

   # Delete user
   curl -X DELETE http://localhost:8080/users/1

File Upload Example
-------------------

.. code-block:: go

   package main

   import (
       "context"
       "fmt"
       "io"
       "net/http"
       "os"
       "path/filepath"

       "github.com/eatmoreapple/hx"
       . "github.com/eatmoreapple/hx/httpx"
   )

   type FileUploadRequest struct {
       Description string `form:"description"`
       // File will be extracted from multipart form
   }

   type FileUploadResponse struct {
       Filename    string `json:"filename"`
       Size        int64  `json:"size"`
       Description string `json:"description"`
       Message     string `json:"message"`
   }

   func uploadFile(ctx context.Context, req FileUploadRequest) (FileUploadResponse, error) {
       // Get the HTTP request from context (you'll need to pass it)
       r := ctx.Value("http_request").(*http.Request)
       
       file, header, err := r.FormFile("file")
       if err != nil {
           return FileUploadResponse{}, fmt.Errorf("failed to get file: %v", err)
       }
       defer file.Close()

       // Create uploads directory
       uploadDir := "./uploads"
       os.MkdirAll(uploadDir, 0755)

       // Create destination file
       dst, err := os.Create(filepath.Join(uploadDir, header.Filename))
       if err != nil {
           return FileUploadResponse{}, fmt.Errorf("failed to create file: %v", err)
       }
       defer dst.Close()

       // Copy file content
       size, err := io.Copy(dst, file)
       if err != nil {
           return FileUploadResponse{}, fmt.Errorf("failed to save file: %v", err)
       }

       return FileUploadResponse{
           Filename:    header.Filename,
           Size:        size,
           Description: req.Description,
           Message:     "File uploaded successfully",
       }, nil
   }

   func main() {
       router := hx.New()
       
       // Serve upload form
       router.GET("/upload", hx.Warp(func(w http.ResponseWriter, r *http.Request) {
           html := `
           <html>
           <body>
               <form action="/upload" method="post" enctype="multipart/form-data">
                   <input type="file" name="file" required><br><br>
                   <input type="text" name="description" placeholder="Description"><br><br>
                   <input type="submit" value="Upload">
               </form>
           </body>
           </html>`
           w.Header().Set("Content-Type", "text/html")
           w.Write([]byte(html))
       }))

       router.POST("/upload", hx.G(uploadFile).JSON())

       fmt.Println("Server starting on :8080")
       fmt.Println("Visit http://localhost:8080/upload to upload files")
       http.ListenAndServe(":8080", router)
   }

Middleware Example
------------------

.. code-block:: go

   package main

   import (
       "context"
       "fmt"
       "log"
       "net/http"
       "time"

       "github.com/eatmoreapple/hx"
       . "github.com/eatmoreapple/hx/httpx"
   )

   // Logging middleware
   func loggingMiddleware(next hx.HandlerFunc) hx.HandlerFunc {
       return func(w http.ResponseWriter, r *http.Request) error {
           start := time.Now()
           log.Printf("Started %s %s", r.Method, r.URL.Path)
           
           err := next(w, r)
           
           duration := time.Since(start)
           log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, duration)
           
           return err
       }
   }

   // CORS middleware
   func corsMiddleware(next hx.HandlerFunc) hx.HandlerFunc {
       return func(w http.ResponseWriter, r *http.Request) error {
           w.Header().Set("Access-Control-Allow-Origin", "*")
           w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
           w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
           
           if r.Method == http.MethodOptions {
               w.WriteHeader(http.StatusOK)
               return nil
           }
           
           return next(w, r)
       }
   }

   // Authentication middleware
   func authMiddleware(next hx.HandlerFunc) hx.HandlerFunc {
       return func(w http.ResponseWriter, r *http.Request) error {
           token := r.Header.Get("Authorization")
           if token == "" {
               return fmt.Errorf("authorization header required")
           }
           
           // Validate token (simplified)
           if token != "Bearer valid-token" {
               return fmt.Errorf("invalid token")
           }
           
           return next(w, r)
       }
   }

   func publicHandler(ctx context.Context, req Empty) (string, error) {
       return "This is a public endpoint", nil
   }

   func protectedHandler(ctx context.Context, req Empty) (string, error) {
       return "This is a protected endpoint", nil
   }

   func main() {
       router := hx.New()

       // Global middleware
       router.Use(loggingMiddleware, corsMiddleware)

       // Public routes
       router.GET("/public", hx.E(publicHandler).String())

       // Protected routes group
       protected := router.Group("/api")
       protected.Use(authMiddleware)
       protected.GET("/protected", hx.E(protectedHandler).String())

       fmt.Println("Server starting on :8080")
       http.ListenAndServe(":8080", router)
   }

Test the middleware:

.. code-block:: bash

   # Public endpoint (works)
   curl http://localhost:8080/public

   # Protected endpoint without auth (fails)
   curl http://localhost:8080/api/protected

   # Protected endpoint with auth (works)
   curl -H "Authorization: Bearer valid-token" http://localhost:8080/api/protected

Custom Error Handling
----------------------

.. code-block:: go

   package main

   import (
       "context"
       "encoding/json"
       "fmt"
       "net/http"

       "github.com/eatmoreapple/hx"
       . "github.com/eatmoreapple/hx/httpx"
   )

   // Custom error types
   type APIError struct {
       Code    int    `json:"code"`
       Message string `json:"message"`
       Details string `json:"details,omitempty"`
   }

   func (e APIError) Error() string {
       return e.Message
   }

   type ValidationError struct {
       Field   string `json:"field"`
       Message string `json:"message"`
   }

   func (e ValidationError) Error() string {
       return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
   }

   // Custom error handler
   func customErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
       w.Header().Set("Content-Type", "application/json")

       var statusCode int
       var response interface{}

       switch e := err.(type) {
       case APIError:
           statusCode = e.Code
           response = e
       case ValidationError:
           statusCode = http.StatusBadRequest
           response = map[string]interface{}{
               "error": "validation_error",
               "field": e.Field,
               "message": e.Message,
           }
       default:
           statusCode = http.StatusInternalServerError
           response = map[string]interface{}{
               "error": "internal_error",
               "message": "An internal error occurred",
           }
       }

       w.WriteHeader(statusCode)
       json.NewEncoder(w).Encode(response)
   }

   func successHandler(ctx context.Context, req Empty) (string, error) {
       return "Success!", nil
   }

   func apiErrorHandler(ctx context.Context, req Empty) (string, error) {
       return "", APIError{
           Code:    http.StatusNotFound,
           Message: "Resource not found",
           Details: "The requested resource could not be found",
       }
   }

   func validationErrorHandler(ctx context.Context, req Empty) (string, error) {
       return "", ValidationError{
           Field:   "email",
           Message: "Invalid email format",
       }
   }

   func panicHandler(ctx context.Context, req Empty) (string, error) {
       return "", fmt.Errorf("something went wrong")
   }

   func main() {
       router := hx.New(hx.WithErrorHandler(customErrorHandler))

       router.GET("/success", hx.E(successHandler).String())
       router.GET("/api-error", hx.E(apiErrorHandler).String())
       router.GET("/validation-error", hx.E(validationErrorHandler).String())
       router.GET("/panic", hx.E(panicHandler).String())

       fmt.Println("Server starting on :8080")
       http.ListenAndServe(":8080", router)
   }

JWT Authentication Example
--------------------------

.. code-block:: go

   package main

   import (
       "context"
       "encoding/json"
       "fmt"
       "net/http"
       "strings"
       "time"

       "github.com/eatmoreapple/hx"
       . "github.com/eatmoreapple/hx/httpx"
   )

   // Simple JWT-like implementation (use a real JWT library in production)
   type Claims struct {
       UserID   string    `json:"user_id"`
       Username string    `json:"username"`
       IssuedAt time.Time `json:"issued_at"`
   }

   func generateToken(userID, username string) string {
       claims := Claims{
           UserID:   userID,
           Username: username,
           IssuedAt: time.Now(),
       }
       data, _ := json.Marshal(claims)
       return string(data) // In production, use proper JWT signing
   }

   func parseToken(token string) (*Claims, error) {
       var claims Claims
       err := json.Unmarshal([]byte(token), &claims)
       if err != nil {
           return nil, fmt.Errorf("invalid token")
       }
       
       // Check token age (24 hours)
       if time.Since(claims.IssuedAt) > 24*time.Hour {
           return nil, fmt.Errorf("token expired")
       }
       
       return &claims, nil
   }

   // Request types
   type LoginRequest struct {
       Username string `json:"username" form:"username"`
       Password string `json:"password" form:"password"`
   }

   type LoginResponse struct {
       Token    string `json:"token"`
       Username string `json:"username"`
   }

   // Context key for user claims
   type contextKey string
   const userClaimsKey contextKey = "user_claims"

   // JWT middleware
   func jwtMiddleware(next hx.HandlerFunc) hx.HandlerFunc {
       return func(w http.ResponseWriter, r *http.Request) error {
           authHeader := r.Header.Get("Authorization")
           if authHeader == "" {
               return fmt.Errorf("authorization header required")
           }

           parts := strings.Split(authHeader, " ")
           if len(parts) != 2 || parts[0] != "Bearer" {
               return fmt.Errorf("invalid authorization header format")
           }

           claims, err := parseToken(parts[1])
           if err != nil {
               return fmt.Errorf("invalid token: %v", err)
           }

           // Add claims to context
           ctx := context.WithValue(r.Context(), userClaimsKey, claims)
           r = r.WithContext(ctx)

           return next(w, r)
       }
   }

   func login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
       // Simple authentication (use proper password hashing in production)
       if req.Username == "" || req.Password == "" {
           return LoginResponse{}, fmt.Errorf("username and password required")
       }
       
       if req.Username != "admin" || req.Password != "password" {
           return LoginResponse{}, fmt.Errorf("invalid credentials")
       }

       token := generateToken("1", req.Username)
       return LoginResponse{
           Token:    token,
           Username: req.Username,
       }, nil
   }

   func profile(ctx context.Context, req Empty) (map[string]interface{}, error) {
       claims := ctx.Value(userClaimsKey).(*Claims)
       return map[string]interface{}{
           "user_id":  claims.UserID,
           "username": claims.Username,
           "message":  "This is your profile",
       }, nil
   }

   func main() {
       router := hx.New()

       // Public route
       router.POST("/login", hx.G(login).JSON())

       // Protected routes
       protected := router.Group("/api")
       protected.Use(jwtMiddleware)
       protected.GET("/profile", hx.E(profile).JSON())

       fmt.Println("Server starting on :8080")
       fmt.Println("Login with: curl -X POST http://localhost:8080/login -d '{\"username\":\"admin\",\"password\":\"password\"}' -H 'Content-Type: application/json'")
       http.ListenAndServe(":8080", router)
   }