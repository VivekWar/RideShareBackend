package models

import (
    "time"
)

type User struct {
    ID           int       `json:"id" db:"id"`
    Name         string    `json:"name" db:"name"`
    Email        string    `json:"email" db:"email"`
    Password     string    `json:"password,omitempty" db:"password"`
    Phone        string    `json:"phone" db:"phone"`
    ProfileImage *string   `json:"profileImage,omitempty" db:"profile_image"`
    IsVerified   bool      `json:"isVerified" db:"is_verified"`
    CreatedAt    time.Time `json:"createdAt" db:"created_at"`
    UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

type UserStats struct {
    TotalTrips     int     `json:"totalTrips"`
    CompletedTrips int     `json:"completedTrips"`
    CancelledTrips int     `json:"cancelledTrips"`
    Rating         float64 `json:"rating"`
    ReviewCount    int     `json:"reviewCount"`
}

type UserProfile struct {
    User
    Stats *UserStats `json:"stats,omitempty"`
}
