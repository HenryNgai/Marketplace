package transactions

import (
	"Marketplace/internal/db/models"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SellHandler(c *gin.Context, database *sql.DB) {
	var listing models.Listing
	// Bind the JSON body to the struct
	err := c.ShouldBindJSON(&listing)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to parse body of post request"})
	} else {
		// Validate the input data
		if listing.UserID <= 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid seller ID"})
			return
		}
		if listing.ItemName == "" || listing.Price <= 0 || listing.Quantity <= 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid item details"})
			return
		}

		// Execute query
		listingID, err := InsertSellListing(database, listing)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Item successfully listed for sale", "ListingID": listingID})
		return
	}
}

// InsertSellListing inserts a new sell listing into the listings table
func InsertSellListing(db *sql.DB, listing models.Listing) (int, error) {
	// SQL query to insert a new listing
	query := `
		INSERT INTO listings (item_name, user_id, price, quantity)
		VALUES ($1, $2, $3, $4)
		RETURNING listing_id`

	// Execute the query with the struct values
	var listingID int
	err := db.QueryRow(query, listing.ItemName, listing.UserID, listing.Price, listing.Quantity).Scan(&listingID)
	if err != nil {
		return 0, fmt.Errorf("error inserting listing: %v", err)
	}

	return listingID, nil
}
