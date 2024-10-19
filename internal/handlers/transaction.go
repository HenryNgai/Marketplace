package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"Marketplace/internal/db/models"

	"github.com/gin-gonic/gin"
)

// Struct for buyer
type Buy struct {
	ItemName string  `json:"itemName" binding:"required"`
	UserID   int     `json:"userID" binding:"required"`
	Price    float32 `json:"price" binding:"required"` // What if it is unable to convert?
	Quantity int     `json:"quantity" binding:"required"`
}

func BuyHandler(c *gin.Context, database *sql.DB) {
	var buyTransaction Buy
	// Bind the JSON body to the struct
	err := c.ShouldBindJSON(&buyTransaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to parse body of post request"})
	} else {
		// Validate the input data
		if buyTransaction.UserID <= 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid seller ID"})
			return
		}
		if buyTransaction.ItemName == "" || buyTransaction.Price <= 0 || buyTransaction.Quantity <= 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid item details"})
			return
		}

	}
}

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
func InsertSellListing(db *sql.DB, listing models.Listing) (int64, error) {
	// SQL query to insert a new listing
	query := `
		INSERT INTO listings (item_name, user_id, price, quantity)
		VALUES ($1, $2, $3, $4)
		RETURNING id `

	// Execute the query with the struct values
	var listingID int64
	err := db.QueryRow(query, listing.ItemName, listing.UserID, listing.Price, listing.Quantity).Scan(&listingID)
	if err != nil {
		return 0, fmt.Errorf("error inserting listing: %v", err)
	}

	return listingID, nil
}

// func GetListing(db *sql.DB, item_ID int) {
// 	var listing models.Listing
// 	query := `SELECT id, item_name, user_id, price, quantity FROM listings WHERE id = $1`
// 	err := db.QueryRow(query, item_ID).Scan(&listing.ID, &listing.ItemName, &listing.UserID, &listing.Price, &listing.Quantity)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, fmt.Errorf("listing not found")
// 		}
// 		return nil, err
// 	}
// 	return &listing, nil
// }
