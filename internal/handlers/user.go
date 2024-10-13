package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func RegisterUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "You hit the user registration portion!"})
}
