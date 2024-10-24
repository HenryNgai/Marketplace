package transactions

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RemoveListingHandler(c *gin.Context, database *sql.DB) {
	listingIDstr := c.Query("listingID")
	listingID, err := strconv.Atoi(listingIDstr)
	if int(listingID) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "itemName query parameter is required"})
		return
	}

	// Call GetListing to retrieve the listings from the database
	err = RemoveListing(database, listingID)
	if err != nil {
		if err.Error() == "listing not found" {
			// If no listing was found, return 404
			c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found"})
		} else {
			// For other errors, return 500
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Successfully deleted the listing
	c.JSON(http.StatusOK, gin.H{"message": "Listing removed successfully"})
}

func RemoveListing(db *sql.DB, listingID int) error {
	// Execute the DELETE SQL query
	result, err := db.Exec("DELETE FROM listings WHERE id = $1", listingID)
	if err != nil {
		return err
	}

	// Check how many rows were affected (should be 1 if the listing was deleted)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If no rows were affected, return a "listing not found" error
	if rowsAffected == 0 {
		return errors.New("listing not found")
	}

	return nil // No errors, listing successfully deleted
}
