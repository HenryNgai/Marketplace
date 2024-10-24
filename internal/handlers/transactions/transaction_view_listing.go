package transactions

import (
	"database/sql"
	"fmt"
	"net/http"

	"Marketplace/internal/db/models"

	"github.com/gin-gonic/gin"
)

func GetListingHandler(c *gin.Context, database *sql.DB) {
	itemName := c.Query("itemName")
	if itemName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "itemName query parameter is required"})
		return
	}

	// Call GetListing to retrieve the listings from the database
	listings, err := GetListing(database, itemName)
	if err != nil {
		// If there was an error querying the database, return an internal server error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Query successful but no listings
	if len(listings) == 0 {
		c.JSON(http.StatusOK, gin.H{"listings": []models.Listing{}})
		return
	}

	// Success NOTE: No need for return since c.JSON automatically returns.
	c.JSON(http.StatusOK, gin.H{fmt.Sprintf("listings for %s", itemName): listings})
}

// GetListing returns a list of all the listings that have the specified itemName
func GetListing(db *sql.DB, itemName string) ([]models.Listing, error) {
	query := `SELECT * FROM listings WHERE item_name = $1`
	rows, err := db.Query(query, itemName)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Close after function execution

	var listings []models.Listing // List of listings
	for rows.Next() {
		var listing models.Listing
		err := rows.Scan(&listing.ListingID, &listing.ItemName, &listing.UserID, &listing.Price, &listing.Quantity)
		if err != nil {
			return nil, fmt.Errorf("could not load query results to listing struct %w", err)
		}
		listings = append(listings, listing) // Expensive? Doubles when out of space. Kinda like list in python.
	}
	return listings, nil
}
