package main

import (
    "database/sql"
    "log"
    "rideshare-backend/internal/config"
    "rideshare-backend/internal/handlers"
    "os"
    "github.com/gin-gonic/gin"
    "github.com/rs/cors"
    _ "github.com/lib/pq"
)

func main() {
    // Load configuration
    cfg := config.Load()
    // Get port from environment (Railway provides this)
    port := os.Getenv("PORT")
    if port == "" {
        port = cfg.Port // fallback to config
    }
    // Connect to database (simple connection)
    db, err := sql.Open("postgres", cfg.DatabaseURL)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()
    
    // Test database connection
    if err := db.Ping(); err != nil {
        log.Fatal("Failed to ping database:", err)
    }
    
    // Initialize Gin router
    r := gin.Default()
    
    // Setup CORS
    c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:3000","https://ridesharefrontend.vercel.app"},
        
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{   "Origin",
        "Content-Type", 
        "Accept",
        "Authorization",
        "X-Requested-With",},
        AllowCredentials: true,
    })
    r.Use(func(ctx *gin.Context) {
        c.HandlerFunc(ctx.Writer, ctx.Request)
        ctx.Next()
    })
    
    // Initialize handlers
    authHandler := handlers.NewAuthHandler(db, cfg)
    tripHandler := handlers.NewTripHandler(db)
    userHandler := handlers.NewUserHandler(db)
    
    // Setup routes
    api := r.Group("/api/v1")
    {
        // Auth routes
        auth := api.Group("/auth")
        {
            auth.POST("/register", authHandler.Register)
            auth.POST("/login", authHandler.Login)
            auth.GET("/me", handlers.AuthMiddleware(cfg.JWTSecret), authHandler.GetCurrentUser)
        }
        
        // Protected routes
        protected := api.Group("/")
        protected.Use(handlers.AuthMiddleware(cfg.JWTSecret))
        {
            // Trip routes
            trips := protected.Group("/trips")
            {
                trips.GET("", tripHandler.GetTrips)
                trips.POST("", tripHandler.CreateTrip)
                trips.GET("/:id", tripHandler.GetTrip)
                trips.PUT("/:id", tripHandler.UpdateTrip)
                trips.DELETE("/:id", tripHandler.DeleteTrip)
                trips.POST("/:id/join", tripHandler.JoinTrip)
                trips.POST("/search", tripHandler.SearchTrips)
            }
            
            // User routes
            users := protected.Group("/users")
            {
                users.GET("/trips", userHandler.GetUserTrips)
                users.GET("/profile", userHandler.GetProfile)
                users.PUT("/profile", userHandler.UpdateProfile)
            }
        }
    }
    
    log.Printf("Server starting on port %s", cfg.Port)
    log.Fatal(r.Run(":" + cfg.Port))
}
