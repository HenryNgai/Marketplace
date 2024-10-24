package main

import (
	"Marketplace/internal/db"                    // db package
	"Marketplace/internal/handlers/transactions" // Handlers package
	"database/sql"
	"log"      // Logging
	"net/http" // HTTP
	"os"

	"github.com/gin-gonic/gin" // Gin framework
)

func main() {
	// Connect to DB and verify connection (ping)
	database, err := db.ConnectToPostgres()
	if err != nil {
		log.Fatalf("Please try connecting to database again %v", err)
	}
	log.Println("Successfuly connected and printed to database")
	defer database.Close() // Delay closing until main program ends

	router := gin.Default() // Creates Gin router. Provides logging and recovery as well.

	registerRoutes(router, database) // Register routes

	// Run the server
	port := os.Getenv("PORT") // Use environment variable for the port
	if port == "" {
		port = "8080" // Default to port 8080
	}
	log.Fatal(router.Run(":" + port))

}

// Register API routes
func registerRoutes(router *gin.Engine, database *sql.DB) {
	router.GET("/ping", PingHandler)
	router.POST("/sellItem", func(c *gin.Context) { // Anonymous function (also known as closure). Pass database to handler.
		transactions.SellHandler(c, database)
	})
	router.POST("/buyItem", func(c *gin.Context) {
		transactions.BuyHandler(c, database)
	})
	router.GET("/getListing", func(c *gin.Context) {
		transactions.GetListingHandler(c, database)
	})
	router.POST("/buyItem", func(c *gin.Context) {
		transactions.RemoveListingHandler(c, database)
	})
}

// Check server status
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
