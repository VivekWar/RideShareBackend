package handlers

import (
    "net/http"
    "rideshare-backend/internal/services"
    "rideshare-backend/internal/models"
    
    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    "database/sql"
)

type UserHandler struct {
    db          *sql.DB
    authService *services.AuthService
    validator   *validator.Validate
    tripService *services.TripService 
}
// ✅ ADD THIS METHOD TO YOUR EXISTING UserHandler
func (h *UserHandler) GetUserTrips(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    
    trips, err := h.tripService.GetUserTrips(userID.(int))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user trips"})
        return
    }
    
    c.JSON(http.StatusOK, trips)
}

func NewUserHandler(db *sql.DB) *UserHandler {
    return &UserHandler{
        db:          db,
        authService: services.NewAuthService(db, nil),
        validator:   validator.New(),
    }
}

type UpdateProfileRequest struct {
    Name         string  `json:"name" validate:"required,min=2,max=100"`
    Phone        string  `json:"phone" validate:"required"`
    ProfileImage *string `json:"profileImage"`
}

func (h *UserHandler) GetProfile(c *gin.Context) {
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

func (h *UserHandler) UpdateProfile(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    
    var req UpdateProfileRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }
    
    if err := h.validator.Struct(req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed"})
        return
    }
    
    user := &models.User{
        ID:           userID.(int),
        Name:         req.Name,
        Phone:        req.Phone,
        ProfileImage: req.ProfileImage,
    }
    
    if err := h.authService.UpdateUser(user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
        return
    }
    
    // Get updated user data
    updatedUser, err := h.authService.GetUserByID(userID.(int))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated profile"})
        return
    }
    
    // Remove password from response
    updatedUser.Password = ""
    
    c.JSON(http.StatusOK, updatedUser)
}

func (h *UserHandler) GetUserStats(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    
    // Get user statistics
    var stats models.UserStats
    query := `
        SELECT 
            COUNT(*) as total_trips,
            COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_trips,
            COUNT(CASE WHEN status = 'cancelled' THEN 1 END) as cancelled_trips,
            COALESCE(AVG(CASE WHEN status = 'completed' THEN 5.0 END), 0) as rating,
            COUNT(CASE WHEN status = 'completed' THEN 1 END) as review_count
        FROM trips 
        WHERE driver_id = $1
    `
    
    err := h.db.QueryRow(query, userID.(int)).Scan(
        &stats.TotalTrips,
        &stats.CompletedTrips,
        &stats.CancelledTrips,
        &stats.Rating,
        &stats.ReviewCount,
    )
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user stats"})
        return
    }
    
    c.JSON(http.StatusOK, stats)
}
