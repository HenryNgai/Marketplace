package models

type Listing struct {
	// Struct for listing from seller
	ItemName string  `json:"itemName" binding:"required"`
	UserID   int     `json:"userID" binding:"required"`
	Price    float32 `json:"price" binding:"required"` // What if it is unable to convert?
	Quantity int     `json:"quantity" binding:"required"`
}
