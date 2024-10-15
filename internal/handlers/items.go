package handlers

import (
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
	UserID   string  `json:"userID" binding:"required"`
	Price    float32 `json:"price" binding:"required"` // What if it is unable to convert?
	Quantity int     `json:"quantity" binding:"required"`
}

// Struct for seller
type Sell struct {
	ItemName string  `json:"itemName" binding:"required"`
	UserID   string  `json:"userID" binding:"required"`
	Price    float32 `json:"price" binding:"required"` // What if it is unable to convert?
	Quantity int     `json:"quantity" binding:"required"`
}

func BuyHandler(c *gin.Context) {
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

func SellHandler(c *gin.Context) {
	var sellTransaction Sell
	// Bind the JSON body to the newItem struct
	if err := c.ShouldBindJSON(&sellTransaction); err != nil {
		// If there's an error in binding, return 400 with the error message.
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // TODO incorporate missing struct field for error
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"Transaction ConfirmationID": 2,
			"itenName": sellTransaction.ItemName,
			"userID":   sellTransaction.UserID,
			"price":    sellTransaction.Price,
			"quantity": sellTransaction.Quantity})
		return
	}
}
