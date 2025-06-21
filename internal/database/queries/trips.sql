-- name: GetTripByID :one
SELECT t.id, t.driver_id, t.from_location, t.to_location, t.departure_time,
       t.max_passengers, t.current_passengers, t.price_per_person, t.description,
       t.status, t.created_at, t.updated_at,
       u.id, u.name, u.email, u.phone
FROM trips t
LEFT JOIN users u ON t.driver_id = u.id
WHERE t.id = $1;

-- name: GetUserTrips :many
SELECT t.id, t.driver_id, t.from_location, t.to_location, t.departure_time,
       t.max_passengers, t.current_passengers, t.price_per_person, t.description,
       t.status, t.created_at, t.updated_at
FROM trips t
WHERE t.driver_id = $1
ORDER BY t.departure_time DESC;

-- name: CreateTrip :one
INSERT INTO trips (driver_id, from_location, to_location, departure_time, max_passengers, price_per_person, description, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, 'active', NOW(), NOW())
RETURNING id, created_at, updated_at;

-- name: UpdateTrip :one
UPDATE trips
SET from_location = $1, to_location = $2, departure_time = $3,
    max_passengers = $4, price_per_person = $5, description = $6, updated_at = NOW()
WHERE id = $7 AND driver_id = $8
RETURNING updated_at;

-- name: DeleteTrip :exec
UPDATE trips
SET status = 'cancelled', updated_at = NOW()
WHERE id = $1 AND driver_id = $2;

-- name: SearchTrips :many
SELECT t.id, t.driver_id, t.from_location, t.to_location, t.departure_time,
       t.max_passengers, t.current_passengers, t.price_per_person, t.description,
       t.status, t.created_at, t.updated_at,
       u.id, u.name, u.email, u.phone
FROM trips t
LEFT JOIN users u ON t.driver_id = u.id
WHERE t.status = 'active' 
AND t.current_passengers < t.max_passengers
AND t.departure_time > NOW()
ORDER BY t.departure_time ASC;

-- name: JoinTrip :exec
INSERT INTO trip_passengers (trip_id, passenger_id, status, joined_at)
VALUES ($1, $2, 'confirmed', NOW());

-- name: UpdatePassengerCount :exec
UPDATE trips
SET current_passengers = current_passengers + 1, updated_at = NOW()
WHERE id = $1;
