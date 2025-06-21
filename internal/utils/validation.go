package utils

import (
    "fmt"
    "strings"
    
    "github.com/go-playground/validator/v10"
)

func FormatValidationErrors(err error) string {
    var errors []string
    
    if validationErrors, ok := err.(validator.ValidationErrors); ok {
        for _, e := range validationErrors {
            switch e.Tag() {
            case "required":
                errors = append(errors, fmt.Sprintf("%s is required", e.Field()))
            case "email":
                errors = append(errors, fmt.Sprintf("%s must be a valid email", e.Field()))
            case "min":
                errors = append(errors, fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param()))
            case "max":
                errors = append(errors, fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param()))
            default:
                errors = append(errors, fmt.Sprintf("%s is invalid", e.Field()))
            }
        }
    }
    
    return strings.Join(errors, ", ")
}

func ValidateEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func ValidatePassword(password string) error {
    if len(password) < 6 {
        return fmt.Errorf("password must be at least 6 characters long")
    }
    return nil
}
