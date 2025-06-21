package models

import (
    "time"
)

type TripPassenger struct {
    ID          int       `json:"id" db:"id"`
    TripID      int       `json:"tripId" db:"trip_id"`
    PassengerID int       `json:"passengerId" db:"passenger_id"`
    Status      string    `json:"status" db:"status"`
    JoinedAt    time.Time `json:"joinedAt" db:"joined_at"`
    
    // Relationships
    Trip      *Trip `json:"trip,omitempty"`
    Passenger *User `json:"passenger,omitempty"`
}

type TripMatch struct {
    TripID      int     `json:"tripId"`
    PassengerID int     `json:"passengerId"`
    Similarity  float64 `json:"similarity"`
    Distance    float64 `json:"distance"`
}
