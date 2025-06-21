package models

import (
    "time"
)

type Trip struct {
    ID                int       `json:"id" db:"id"`
    DriverID          int       `json:"driverId" db:"driver_id"`
    FromLocation      string    `json:"from" db:"from_location"`
    ToLocation        string    `json:"to" db:"to_location"`
    DepartureTime     time.Time `json:"departureTime" db:"departure_time"`
    MaxPassengers     int       `json:"maxPassengers" db:"max_passengers"`
    CurrentPassengers int       `json:"currentPassengers" db:"current_passengers"`
    PricePerPerson    float64   `json:"pricePerPerson" db:"price_per_person"`
    Description       *string   `json:"description,omitempty" db:"description"`
    Status            string    `json:"status" db:"status"`
    CreatedAt         time.Time `json:"createdAt" db:"created_at"`
    UpdatedAt         time.Time `json:"updatedAt" db:"updated_at"`
    
    // Relationships
    Driver     *User  `json:"driver,omitempty"`
    Passengers []User `json:"passengers,omitempty"`
    UserRole  string `json:"userRole,omitempty"`
}

type TripSearchCriteria struct {
    From          string  `json:"from"`
    To            string  `json:"to"`
    DepartureDate string  `json:"departureDate"`
    MaxPrice      float64 `json:"maxPrice"`
    Limit         int     `json:"limit"`
    Offset        int     `json:"offset"`
}
