package main

import (
	"database/sql"
	"log"
	"os"
	"rideshare-backend/internal/config"
	"rideshare-backend/internal/handlers"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg := config.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Port
	}
	
	// Database connection
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	
	// Initialize Gin router
	r := gin.Default()

	// Configure CORS with rs/cors
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://ridesharefrontend.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Origin"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           int(12 * time.Hour.Seconds()),
		Debug:            true, // Enable for troubleshooting
	})

	// Apply CORS middleware
	r.Use(func(ctx *gin.Context) {
		corsMiddleware.HandlerFunc(ctx.Writer, ctx.Request)
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}
		ctx.Next()
	})
    
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, cfg)
	tripHandler := handlers.NewTripHandler(db)
	userHandler := handlers.NewUserHandler(db)
	
	// Setup routes
	api := r.Group("/")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", handlers.AuthMiddleware(cfg.JWTSecret), authHandler.GetCurrentUser)
		}
		
		protected := api.Group("/")
		protected.Use(handlers.AuthMiddleware(cfg.JWTSecret))
		{
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
