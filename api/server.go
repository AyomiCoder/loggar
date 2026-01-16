package api

import (
	"database/sql"
	"log"

	"github.com/AyomiCoder/loggar/api/handlers"
	"github.com/AyomiCoder/loggar/api/middleware"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

// InitDB initializes the database connection
func InitDB(databaseURL string) error {
	var err error
	db, err = sql.Open("postgres", databaseURL)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	// Set the database for handlers
	handlers.SetDB(db)

	log.Println("Database connected successfully")
	return nil
}

// NewServer creates and configures the Gin server
func NewServer() *gin.Engine {
	router := gin.Default()

	// Health check endpoint
	// Health check endpoint
	healthHandler := func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	}
	router.GET("/health", healthHandler)
	router.HEAD("/health", healthHandler)

	// Authentication routes
	auth := router.Group("/auth")
	{
		auth.GET("/github", handlers.AuthGitHubHandler)
		auth.GET("/github/callback", handlers.AuthGitHubCallbackHandler)
		auth.GET("/google", handlers.AuthGoogleHandler)
		auth.GET("/google/callback", handlers.AuthGoogleCallbackHandler)
	}

	// Protected routes (require JWT)
	apiRoutes := router.Group("/api")
	apiRoutes.Use(middleware.AuthMiddleware())
	{
		apiRoutes.POST("/analyze", handlers.AnalyzeHandler)
	}

	return router
}

// Run starts the server on the specified port
func Run(port string) error {
	router := NewServer()
	return router.Run(":" + port)
}
