-- name: GetUserByID :one
SELECT id, name, email, phone, profile_image, is_verified, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, name, email, password, phone, profile_image, is_verified, created_at, updated_at
FROM users
WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (name, email, password, phone, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
RETURNING id, created_at, updated_at;

-- name: UpdateUser :one
UPDATE users
SET name = $1, phone = $2, profile_image = $3, updated_at = NOW()
WHERE id = $4
RETURNING updated_at;

-- name: UpdateUserPassword :exec
UPDATE users
SET password = $1, updated_at = NOW()
WHERE id = $2;

-- name: VerifyUser :exec
UPDATE users
SET is_verified = true, updated_at = NOW()
WHERE id = $1;

-- name: GetUserStats :one
SELECT 
    COUNT(*) as total_trips,
    COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_trips,
    COUNT(CASE WHEN status = 'cancelled' THEN 1 END) as cancelled_trips
FROM trips 
WHERE driver_id = $1;
