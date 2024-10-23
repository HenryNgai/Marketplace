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
		listings = append(listings, listing) // Expensive? Doubles when our of space. Kinda like list in python.
	}
	return listings, nil
}
