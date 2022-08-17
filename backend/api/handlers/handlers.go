package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mvlberry/backend/api/database"
)

// function to create account when api request is made
func Signup(c *gin.Context) {
	var newUser database.User

	//retrieve account details from POST request, which is sent as JSON
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Could not parse JSON",
		})
	}

	//checks if account could be successfully added to database
	switch res := database.Adduser(newUser); res {
	case 0:
		c.JSON(http.StatusOK, gin.H{
			"msg": "Account successfully created. Redirecting to login page",
		})
	case 1:
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Username already exists",
		})
	case 2:
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Could not create account",
		})

	}
}
