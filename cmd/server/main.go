package main

import (
	"database/sql"
	"log"
	"os"
	"time"
	"rideshare-backend/internal/config"
	"rideshare-backend/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg := config.Load()
	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Port // fallback to config
	}
	
	// Connect to database
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

	// Configure CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://ridesharefrontend.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Add explicit OPTIONS handler
	r.OPTIONS("/*any", func(c *gin.Context) {
		c.AbortWithStatus(204) // No content for OPTIONS requests
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
	
	log.Printf("Server starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}
