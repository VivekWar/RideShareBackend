package handlers

import (
    "net/http"
    "strconv"
    "time"
    "rideshare-backend/internal/models"
    "rideshare-backend/internal/services"
    "rideshare-backend/internal/utils"
    
    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    "database/sql"
)

type TripHandler struct {
    db          *sql.DB
    tripService *services.TripService
    validator   *validator.Validate
}

func NewTripHandler(db *sql.DB) *TripHandler {
    return &TripHandler{
        db:          db,
        tripService: services.NewTripService(db),
        validator:   validator.New(),
    }
}

type CreateTripRequest struct {
    From            string  `json:"from" validate:"required"`
    To              string  `json:"to" validate:"required"`
    DepartureTime   string  `json:"departureTime" validate:"required"`
    MaxPassengers   int     `json:"maxPassengers" validate:"required,min=1,max=8"`
    PricePerPerson  float64 `json:"pricePerPerson" validate:"required,min=0"`
    Description     string  `json:"description"`
}

type SearchTripsRequest struct {
    From          string `json:"from"`
    To            string `json:"to"`
    DepartureDate string `json:"departureDate"`
    MaxPrice      float64 `json:"maxPrice"`
}

func (h *TripHandler) CreateTrip(c *gin.Context) {
    var req CreateTripRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }
    
    if err := h.validator.Struct(req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationErrors(err)})
        return
    }
    
    userID, _ := c.Get("userID")
    
    // Parse departure time
    departureTime, err := time.Parse(time.RFC3339, req.DepartureTime)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid departure time format"})
        return
    }

    // Convert description to pointer
    var description *string
    if req.Description != "" {
        description = &req.Description
    }
    
    trip := &models.Trip{
        DriverID:        userID.(int),
        FromLocation:    req.From,
        ToLocation:      req.To,
        DepartureTime:   departureTime,
        MaxPassengers:   req.MaxPassengers,
        PricePerPerson:  req.PricePerPerson,
        Description:     description,
        Status:          "active",
    }
    
    if err := h.tripService.CreateTrip(trip); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create trip"})
        return
    }
    
    c.JSON(http.StatusCreated, trip)
}

func (h *TripHandler) GetTrips(c *gin.Context) {
    userID, _ := c.Get("userID")
    
    trips, err := h.tripService.GetUserTrips(userID.(int))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trips"})
        return
    }
    
    c.JSON(http.StatusOK, trips)
}

func (h *TripHandler) GetTrip(c *gin.Context) {
    tripID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trip ID"})
        return
    }
    
    trip, err := h.tripService.GetTripByID(tripID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
        return
    }
    
    c.JSON(http.StatusOK, trip)
}

func (h *TripHandler) UpdateTrip(c *gin.Context) {
    tripID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trip ID"})
        return
    }
    
    var req CreateTripRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }
    
    if err := h.validator.Struct(req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": utils.FormatValidationErrors(err)})
        return
    }
    
    userID, _ := c.Get("userID")
    
    // Parse departure time
    departureTime, err := time.Parse(time.RFC3339, req.DepartureTime)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid departure time format"})
        return
    }

    // Convert description to pointer
    var description *string
    if req.Description != "" {
        description = &req.Description
    }
    
    trip := &models.Trip{
        ID:              tripID,
        DriverID:        userID.(int),
        FromLocation:    req.From,
        ToLocation:      req.To,
        DepartureTime:   departureTime,
        MaxPassengers:   req.MaxPassengers,
        PricePerPerson:  req.PricePerPerson,
        Description:     description,
    }
    
    if err := h.tripService.UpdateTrip(trip); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update trip"})
        return
    }
    
    c.JSON(http.StatusOK, trip)
}

func (h *TripHandler) DeleteTrip(c *gin.Context) {
    tripID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trip ID"})
        return
    }
    
    userID, _ := c.Get("userID")
    
    if err := h.tripService.DeleteTrip(tripID, userID.(int)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete trip"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Trip deleted successfully"})
}

func (h *TripHandler) SearchTrips(c *gin.Context) {
    var req SearchTripsRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }
    
    trips, err := h.tripService.SearchTrips(req.From, req.To, req.DepartureDate, req.MaxPrice)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search trips"})
        return
    }
    
    c.JSON(http.StatusOK, trips)
}

func (h *TripHandler) JoinTrip(c *gin.Context) {
    tripID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trip ID"})
        return
    }
    
    userID, _ := c.Get("userID")
    
    if err := h.tripService.JoinTrip(tripID, userID.(int)); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Successfully joined trip"})
}