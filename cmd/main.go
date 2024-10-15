package main

import (
	"Marketplace/internal/handlers" // Handlers package
	"log"                           // Logging
	"net/http"                      // HTTP
	"os"

	"github.com/gin-gonic/gin" // Gin framework
)

func main() {
	router := gin.Default() // Creates Gin router. Provides logging and recovery as well.

	registerRoutes(router) // Register routes

	// Run the server
	port := os.Getenv("PORT") // Use environment variable for the port
	if port == "" {
		port = "8080" // Default to port 8080
	}
	log.Fatal(router.Run(":" + port))

}

// Register API routes
func registerRoutes(router *gin.Engine) {
	router.GET("/ping", PingHandler)
	router.POST("/sellItem", handlers.SellHandler)
	router.POST("/buyItem", handlers.BuyHandler)
}

// Check server status
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
