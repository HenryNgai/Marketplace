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
	c.JSON(http.StatusOK, gin.H{"Transaction Confirmation ID": "99817289"})
}
