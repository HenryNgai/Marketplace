package models

// Struct for listing from seller
type Listing struct {
	ListingID int     `json:"listingID" binding:"required"` // Primary Key, Auto Generated - DB
	ItemName  string  `json:"itemName" binding:"required"`
	UserID    int     `json:"userID" binding:"required"`
	Price     float32 `json:"price" binding:"required"` // What if it is unable to convert?
	Quantity  int     `json:"quantity" binding:"required"`
}

// Struct for purchase from buyer
type Purchase struct {
	ListingID int `json:"listingID" binding:"required"`
	Quantity  int `json:"quantity" binding:"required"`
}

// Struct for items
type Item struct {
	ItemName    string `json:"itemName" binding:"required"`
	ItemID      int    `json:"itemID" binding:"required"`
	Description string `json:"description" binding:"required"`
}
