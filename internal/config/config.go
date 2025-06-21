package config

import (
    "os"

    
    "github.com/joho/godotenv"
)

type Config struct {
    DatabaseURL    string
    JWTSecret      string
    Port           string
    EmailHost      string
    EmailPort      int
    EmailUser      string
    EmailPassword  string
    FrontendURL    string
    Environment    string
}

func Load() *Config {
    // Load .env file in development
    godotenv.Load()
    
    return &Config{
        DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost:5432/rideshare_db"),
        JWTSecret:   getEnv("JWT_SECRET", "skfnlsfnlvksfnsfnlfsknslfknlfksnfslknfslvknspfa[ir[ijifnfpnf;vnlsfnkbnklfbn;sfnf;amokknfnflbnlflnbfbnslncnlnfnflnfnanlfsnflnlkfbnljbfnlkfnlnclnlvnlblnlfnnlslkflksldjfblbjlbjlbjgdlb;;dgsbjbsjbs;j"),
        Port:        getEnv("PORT", "8080"),
    }
}
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
