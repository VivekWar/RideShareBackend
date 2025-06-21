package services

import (
    "database/sql"
    "fmt"
    "rideshare-backend/internal/config"
    "rideshare-backend/internal/models"
)

type AuthService struct {
    db     *sql.DB
    config *config.Config
}

func NewAuthService(db *sql.DB, cfg *config.Config) *AuthService {
    return &AuthService{
        db:     db,
        config: cfg,
    }
}

func (s *AuthService) CreateUser(user *models.User) error {
    query := `
        INSERT INTO users (name, email, password, phone, created_at, updated_at)
        VALUES ($1, $2, $3, $4, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `
    
    err := s.db.QueryRow(
        query,
        user.Name,
        user.Email,
        user.Password,
        user.Phone,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
    
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    
    return nil
}

func (s *AuthService) GetUserByEmail(email string) (*models.User, error) {
    user := &models.User{}
    query := `
        SELECT id, name, email, password, phone, profile_image, is_verified, created_at, updated_at
        FROM users
        WHERE email = $1
    `
    
    err := s.db.QueryRow(query, email).Scan(
        &user.ID,
        &user.Name,
        &user.Email,
        &user.Password,
        &user.Phone,
        &user.ProfileImage,
        &user.IsVerified,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    return user, nil
}

func (s *AuthService) GetUserByID(id int) (*models.User, error) {
    user := &models.User{}
    query := `
        SELECT id, name, email, phone, profile_image, is_verified, created_at, updated_at
        FROM users
        WHERE id = $1
    `
    
    err := s.db.QueryRow(query, id).Scan(
        &user.ID,
        &user.Name,
        &user.Email,
        &user.Phone,
        &user.ProfileImage,
        &user.IsVerified,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    return user, nil
}

func (s *AuthService) UpdateUser(user *models.User) error {
    query := `
        UPDATE users
        SET name = $1, phone = $2, profile_image = $3, updated_at = NOW()
        WHERE id = $4
        RETURNING updated_at
    `
    
    err := s.db.QueryRow(
        query,
        user.Name,
        user.Phone,
        user.ProfileImage,
        user.ID,
    ).Scan(&user.UpdatedAt)
    
    if err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }
    
    return nil
}

func (s *AuthService) SendWelcomeEmail(email, name string) error {
    // Implement email sending logic here
    // This is a placeholder implementation
    fmt.Printf("Sending welcome email to %s (%s)\n", name, email)
    return nil
}
