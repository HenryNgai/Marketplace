package transactions

import (
	"Marketplace/internal/db/models"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BuyHandler(c *gin.Context, database *sql.DB) {
	var purchase models.Purchase
	// Bind the JSON body to the struct
	err := c.ShouldBindJSON(&purchase)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to parse body of post request"})
	} else {
		// Validate the input data
		if purchase.ListingID <= 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid listingID"})
			return
		}
		if purchase.Quantity <= 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid quantity"})
			return
		}
		// Execute query
		listingID, err := PurchaseListing(database, purchase)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Item Successfully Purchased", "ListingID Purchased": listingID})
		return

	}
}

// PurchaseListing executes a query to update the quantity of the listing in the DB based on the purchase.
func PurchaseListing(db *sql.DB, purchase models.Purchase) (int, error) {
	// SQL query to update the quantity of the listing based on the ListingID and ensure sufficient stock
	query := `
	UPDATE listings
	SET quantity = quantity - $1
	WHERE listing_id = $2
        AND quantity >= $1;
	`
	// Execute the query with the provided quantity and ListingID from the purchase
	result, err := db.Exec(query, purchase.Quantity, purchase.ListingID)

	// Handle query execution failure
	if err != nil {
		return 0, fmt.Errorf("could not execute purchase query: %w", err)
	}

	// Check the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("could not fetch rows affected: %w", err)
	}

	// Check if the query affected any rows (indicating a successful purchase)
	if rowsAffected == 0 {
		return 0, fmt.Errorf("purchase failed: insufficient quantity or invalid ListingID")
	}

	// Return the ListingID and a nil error to indicate a successful purchase
	return purchase.ListingID, nil
}
