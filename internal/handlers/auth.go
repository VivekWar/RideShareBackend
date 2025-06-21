package handlers

import (
    "net/http"
    "rideshare-backend/internal/config"
    "rideshare-backend/internal/models"
    "rideshare-backend/internal/services"
    "rideshare-backend/internal/utils"
    
    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    "database/sql"
    "log"
)

type AuthHandler struct {
    db          *sql.DB
    config      *config.Config
    authService *services.AuthService
    validator   *validator.Validate
}

func NewAuthHandler(db *sql.DB, cfg *config.Config) *AuthHandler {
    return &AuthHandler{
        db:          db,
        config:      cfg,
        authService: services.NewAuthService(db, cfg),
        validator:   validator.New(),
    }
}

type RegisterRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=100"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
    Phone    string `json:"phone" validate:"required"`
}

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
    Token string      `json:"token"`
    User  models.User `json:"user"`
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        log.Printf("Invalid request body: %v", err) // Add logging
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }
    
    if err := h.validator.Struct(req); err != nil {
        log.Printf("Validation error: %v", err) // Add logging
        c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationErrors(err)})
        return
    }
    
    // Check if user already exists
    existingUser, err := h.authService.GetUserByEmail(req.Email)
    if err != nil && err.Error() != "user not found" {
        // Only log if it's an unexpected error, not just "user not found"
        log.Printf("Error checking existing user: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user existence"})
        return
    }
    
    if existingUser != nil {
        c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
        return
    }
    
    // Hash password
    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        log.Printf("Password hashing error: %v", err) // Add logging
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
        return
    }
    
    // Create user
    user := &models.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: hashedPassword,
        Phone:    req.Phone,
    }
    
    if err := h.authService.CreateUser(user); err != nil {
        log.Printf("User creation error: %v", err) // Add logging
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }
    
    /// Generate JWT token
    token, err := utils.GenerateJWT(user.ID, user.Email, user.Name, h.config.JWTSecret)
    if err != nil {
        log.Printf("JWT generation error: %v", err) // Add logging
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }
    
    // Send welcome email
    go func() {
        if err := h.authService.SendWelcomeEmail(user.Email, user.Name); err != nil {
            log.Printf("Failed to send welcome email: %v", err)
        }
    }()
    
    // Remove password from response
    user.Password = ""
    
    c.JSON(http.StatusCreated, AuthResponse{
        Token: token,
        User:  *user,
    })
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }
    
    if err := h.validator.Struct(req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationErrors(err)})
        return
    }
    
    // Get user by email
    user, err := h.authService.GetUserByEmail(req.Email)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }
    
    // Check password
    if !utils.CheckPassword(req.Password, user.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }
    
    // Generate JWT token
    token, err := utils.GenerateJWT(user.ID, user.Email, user.Name, h.config.JWTSecret)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    // Remove password from response
    user.Password = ""
    
    c.JSON(http.StatusOK, AuthResponse{
        Token: token,
        User:  *user,
    })
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    
    user, err := h.authService.GetUserByID(userID.(int))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    
    // Remove password from response
    user.Password = ""
    
    c.JSON(http.StatusOK, user)
}
