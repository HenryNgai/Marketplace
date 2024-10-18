package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Struct for item
type Item struct {
	ItemName    string `json:"itemName" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// Struct for buyer
type Buy struct {
	ItemName string  `json:"itemName" binding:"required"`
	UserID   int     `json:"userID" binding:"required"`
	Price    float32 `json:"price" binding:"required"` // What if it is unable to convert?
	Quantity int     `json:"quantity" binding:"required"`
}

// Struct for seller
type Sell struct {
	ItemName string  `json:"itemName" binding:"required"`
	UserID   int     `json:"userID" binding:"required"`
	Price    float32 `json:"price" binding:"required"` // What if it is unable to convert?
	Quantity int     `json:"quantity" binding:"required"`
}

func BuyHandler(c *gin.Context, database *sql.DB) {
	var buyTransaction Buy
	// Bind the JSON body to the newItem struct
	if err := c.ShouldBindJSON(&buyTransaction); err != nil {
		// If there's an error in binding, return 400 with the error message.
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // TODO incorporate missing struct field for error
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"Transaction ConfirmationID": 1,
			"itenName": buyTransaction.ItemName,
			"userID":   buyTransaction.UserID,
			"price":    buyTransaction.Price,
			"quantity": buyTransaction.Quantity})
		return
	}
}

func SellHandler(c *gin.Context, database *sql.DB) {
	var sellTransaction Sell
	// Bind the JSON body to the newItem struct
	err := c.ShouldBindJSON(&sellTransaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to parse body of post request"})
	} else {
		// Validate the input data
		if sellTransaction.UserID <= 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid seller ID"})
			return
		}
		if sellTransaction.ItemName == "" || sellTransaction.Price <= 0 || sellTransaction.Quantity <= 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid item details"})
			return
		}

		// Execute query
		listingID, err := InsertSellListing(database, sellTransaction)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Item successfully listed for sale", "ListingID": listingID})
		return
	}
}

// InsertSellListing inserts a new sell listing into the listings table
func InsertSellListing(db *sql.DB, sell Sell) (int64, error) {
	// SQL query to insert a new listing
	query := `
		INSERT INTO listings (item_name, user_id, price, quantity)
		VALUES ($1, $2, $3, $4)
		RETURNING id `

	// Execute the query with the struct values
	var listingID int64
	err := db.QueryRow(query, sell.ItemName, sell.UserID, sell.Price, sell.Quantity).Scan(&listingID)
	if err != nil {
		return 0, fmt.Errorf("error inserting listing: %v", err)
	}

	return listingID, nil
}
