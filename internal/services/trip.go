package services

import (
    "database/sql"
    "fmt"
    "rideshare-backend/internal/models"
    "strings"
   
)

type TripService struct {
    db *sql.DB
}

func NewTripService(db *sql.DB) *TripService {
    return &TripService{db: db}
}

func (s *TripService) CreateTrip(trip *models.Trip) error {
    query := `
        INSERT INTO trips (driver_id, from_location, to_location, departure_time, max_passengers, price_per_person, description, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `
    
    err := s.db.QueryRow(
        query,
        trip.DriverID,
        trip.FromLocation,
        trip.ToLocation,
        trip.DepartureTime,
        trip.MaxPassengers,
        trip.PricePerPerson,
        trip.Description,
        "active",
    ).Scan(&trip.ID, &trip.CreatedAt, &trip.UpdatedAt)
    
    if err != nil {
        return fmt.Errorf("failed to create trip: %w", err)
    }
    
    return nil
}
func (s *TripService) GetAllUserTrips(userID int) ([]models.Trip, error) {
    query := `
        SELECT DISTINCT
            t.id, t.driver_id, t.from_location, t.to_location, 
            t.departure_time, t.max_passengers, t.current_passengers,
            t.price_per_person, t.description, t.status,
            t.created_at, t.updated_at,
            u.name, u.email, u.phone,
            CASE 
                WHEN t.driver_id = $1 THEN 'driver'
                WHEN tp.passenger_id = $1 THEN 'passenger'
                ELSE 'none'
            END as user_role
        FROM trips t
        LEFT JOIN users u ON t.driver_id = u.id
        LEFT JOIN trip_passengers tp ON t.id = tp.trip_id
        WHERE t.driver_id = $1 OR tp.passenger_id = $1
        ORDER BY t.departure_time DESC
    `
    
    rows, err := s.db.Query(query, userID)
    if err != nil {
        fmt.Printf("Database query error: %v", err)
        return nil, fmt.Errorf("failed to get user trips: %w", err)
    }
    defer rows.Close()
    
    var trips []models.Trip
    for rows.Next() {
        var trip models.Trip
        var driverName, driverEmail, driverPhone sql.NullString
        var userRole string
        
        err := rows.Scan(
            &trip.ID, &trip.DriverID, &trip.FromLocation, &trip.ToLocation,
            &trip.DepartureTime, &trip.MaxPassengers, &trip.CurrentPassengers,
            &trip.PricePerPerson, &trip.Description, &trip.Status,
            &trip.CreatedAt, &trip.UpdatedAt,
            &driverName, &driverEmail, &driverPhone, &userRole,
        )
        if err != nil {
            fmt.Printf("Row scan error: %v", err)
            continue
        }
        
        // Create driver info if available
        if driverName.Valid {
            trip.Driver = &models.User{
                ID:    trip.DriverID,
                Name:  driverName.String,
                Email: driverEmail.String,
                Phone: driverPhone.String,
            }
        }
        
        trip.UserRole = userRole
        trips = append(trips, trip)
    }
    
    return trips, nil
}
func (s *TripService) GetUserTrips(userID int) ([]models.Trip, error) {
    query := `
        SELECT DISTINCT
            t.id, t.driver_id, t.from_location, t.to_location, 
            t.departure_time, t.max_passengers, t.current_passengers,
            t.price_per_person, t.description, t.status,
            t.created_at, t.updated_at,
            u.name, u.email, u.phone,
            CASE 
                WHEN t.driver_id = $1 THEN 'driver'
                WHEN tp.passenger_id = $1 THEN 'passenger'
                ELSE 'none'
            END as user_role
        FROM trips t
        LEFT JOIN users u ON t.driver_id = u.id
        LEFT JOIN trip_passengers tp ON t.id = tp.trip_id
        WHERE (t.driver_id = $1 OR tp.passenger_id = $1)
        AND t.status != 'cancelled'
        ORDER BY t.departure_time DESC
    `
    
    rows, err := s.db.Query(query, userID)
    if err != nil {
        fmt.Printf("Database query error: %v", err)
        return nil, fmt.Errorf("failed to get user trips: %w", err)
    }
    defer rows.Close()
    
    var trips []models.Trip
    for rows.Next() {
        var trip models.Trip
        var driverName, driverEmail, driverPhone sql.NullString
        var userRole string
        
        err := rows.Scan(
            &trip.ID,
            &trip.DriverID,
            &trip.FromLocation,
            &trip.ToLocation,
            &trip.DepartureTime,
            &trip.MaxPassengers,
            &trip.CurrentPassengers,
            &trip.PricePerPerson,
            &trip.Description,
            &trip.Status,
            &trip.CreatedAt,
            &trip.UpdatedAt,
            &driverName,
            &driverEmail,
            &driverPhone,
            &userRole,
        )
        if err != nil {
            fmt.Printf("Row scan error: %v", err)
            return nil, fmt.Errorf("failed to scan trip: %w", err)
        }
        
        // Create driver info if available
        if driverName.Valid {
            trip.Driver = &models.User{
                ID:    trip.DriverID,
                Name:  driverName.String,
                Email: driverEmail.String,
                Phone: driverPhone.String,
            }
        }
        
        // Add user role information
        trip.UserRole = userRole
        
        trips = append(trips, trip)
    }
    
    return trips, nil
}

func (s *TripService) GetTripByID(id int) (*models.Trip, error) {
    trip := &models.Trip{}
    driver := &models.User{}
    
    query := `
        SELECT t.id, t.driver_id, t.from_location, t.to_location, t.departure_time,
               t.max_passengers, t.current_passengers, t.price_per_person, t.description,
               t.status, t.created_at, t.updated_at,
               u.id, u.name, u.email, u.phone
        FROM trips t
        LEFT JOIN users u ON t.driver_id = u.id
        WHERE t.id = $1
    `
    
    err := s.db.QueryRow(query, id).Scan(
        &trip.ID,
        &trip.DriverID,
        &trip.FromLocation,
        &trip.ToLocation,
        &trip.DepartureTime,
        &trip.MaxPassengers,
        &trip.CurrentPassengers,
        &trip.PricePerPerson,
        &trip.Description,
        &trip.Status,
        &trip.CreatedAt,
        &trip.UpdatedAt,
        &driver.ID,
        &driver.Name,
        &driver.Email,
        &driver.Phone,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("trip not found")
        }
        return nil, fmt.Errorf("failed to get trip: %w", err)
    }
    
    trip.Driver = driver
    return trip, nil
}

func (s *TripService) SearchTrips(from, to, departureDate string, maxPrice float64) ([]models.Trip, error) {
    var conditions []string
    var args []interface{}
    argIndex := 1
    
    baseQuery := `
        SELECT t.id, t.driver_id, t.from_location, t.to_location, t.departure_time,
               t.max_passengers, t.current_passengers, t.price_per_person, t.description,
               t.status, t.created_at, t.updated_at,
               u.id, u.name, u.email, u.phone
        FROM trips t
        LEFT JOIN users u ON t.driver_id = u.id
        WHERE t.status = 'active' AND t.current_passengers < t.max_passengers
    `
    
    if from != "" {
        conditions = append(conditions, fmt.Sprintf("LOWER(t.from_location) LIKE LOWER($%d)", argIndex))
        args = append(args, "%"+from+"%")
        argIndex++
    }
    
    if to != "" {
        conditions = append(conditions, fmt.Sprintf("LOWER(t.to_location) LIKE LOWER($%d)", argIndex))
        args = append(args, "%"+to+"%")
        argIndex++
    }
    
    if departureDate != "" {
        conditions = append(conditions, fmt.Sprintf("DATE(t.departure_time) = $%d", argIndex))
        args = append(args, departureDate)
        argIndex++
    }
    
    if maxPrice > 0 {
        conditions = append(conditions, fmt.Sprintf("t.price_per_person <= $%d", argIndex))
        args = append(args, maxPrice)
        argIndex++
    }
    
    if len(conditions) > 0 {
        baseQuery += " AND " + strings.Join(conditions, " AND ")
    }
    
    baseQuery += " ORDER BY t.departure_time ASC"
    
    rows, err := s.db.Query(baseQuery, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to search trips: %w", err)
    }
    defer rows.Close()
    
    var trips []models.Trip
    for rows.Next() {
        var trip models.Trip
        var driver models.User
        
        err := rows.Scan(
            &trip.ID,
            &trip.DriverID,
            &trip.FromLocation,
            &trip.ToLocation,
            &trip.DepartureTime,
            &trip.MaxPassengers,
            &trip.CurrentPassengers,
            &trip.PricePerPerson,
            &trip.Description,
            &trip.Status,
            &trip.CreatedAt,
            &trip.UpdatedAt,
            &driver.ID,
            &driver.Name,
            &driver.Email,
            &driver.Phone,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan trip: %w", err)
        }
        
        trip.Driver = &driver
        trips = append(trips, trip)
    }
    
    return trips, nil
}

func (s *TripService) JoinTrip(tripID, passengerID int) error {
    tx, err := s.db.Begin()
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback()
    
    // Check if trip exists and has available seats
    var currentPassengers, maxPassengers int
    err = tx.QueryRow(
        "SELECT current_passengers, max_passengers FROM trips WHERE id = $1 AND status = 'active'",
        tripID,
    ).Scan(&currentPassengers, &maxPassengers)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("trip not found or not active")
        }
        return fmt.Errorf("failed to check trip availability: %w", err)
    }
    
    if currentPassengers >= maxPassengers {
        return fmt.Errorf("trip is full")
    }
    
    // Check if user already joined this trip
    var existingID int
    err = tx.QueryRow(
        "SELECT id FROM trip_passengers WHERE trip_id = $1 AND passenger_id = $2",
        tripID, passengerID,
    ).Scan(&existingID)
    
    if err == nil {
        return fmt.Errorf("you have already joined this trip")
    } else if err != sql.ErrNoRows {
        return fmt.Errorf("failed to check existing participation: %w", err)
    }
    
    // Add passenger to trip
    _, err = tx.Exec(
        "INSERT INTO trip_passengers (trip_id, passenger_id, status, joined_at) VALUES ($1, $2, 'confirmed', NOW())",
        tripID, passengerID,
    )
    if err != nil {
        return fmt.Errorf("failed to join trip: %w", err)
    }
    
    // Update current passenger count
    _, err = tx.Exec(
        "UPDATE trips SET current_passengers = current_passengers + 1, updated_at = NOW() WHERE id = $1",
        tripID,
    )
    if err != nil {
        return fmt.Errorf("failed to update passenger count: %w", err)
    }
    
    return tx.Commit()
}

func (s *TripService) UpdateTrip(trip *models.Trip) error {
    query := `
        UPDATE trips
        SET from_location = $1, to_location = $2, departure_time = $3,
            max_passengers = $4, price_per_person = $5, description = $6, updated_at = NOW()
        WHERE id = $7 AND driver_id = $8
        RETURNING updated_at
    `
    
    err := s.db.QueryRow(
        query,
        trip.FromLocation,
        trip.ToLocation,
        trip.DepartureTime,
        trip.MaxPassengers,
        trip.PricePerPerson,
        trip.Description,
        trip.ID,
        trip.DriverID,
    ).Scan(&trip.UpdatedAt)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("trip not found or unauthorized")
        }
        return fmt.Errorf("failed to update trip: %w", err)
    }
    
    return nil
}

func (s *TripService) DeleteTrip(tripID, driverID int) error {
    result, err := s.db.Exec(
        "UPDATE trips SET status = 'cancelled', updated_at = NOW() WHERE id = $1 AND driver_id = $2",
        tripID, driverID,
    )
    if err != nil {
        return fmt.Errorf("failed to delete trip: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get affected rows: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("trip not found or unauthorized")
    }
    
    return nil
}