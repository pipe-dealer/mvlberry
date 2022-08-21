package handlers

import (
	"fmt"
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
		fmt.Printf("User %s was successfully created\n", newUser.Username)
	case 1:
		c.JSON(http.StatusOK, gin.H{
			"msg": "Username already exists",
		})
		fmt.Println("Could not create account: username already exists")
	case 2:
		c.JSON(http.StatusOK, gin.H{
			"msg": "Could not create account",
		})
		fmt.Println("Could not create account")

	}
}

func Login(c *gin.Context) {
	var user database.User
	//get login details
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Could not parse JSON",
		})
	}
	//get all users
	allusers := database.Getusers()

	for _, v := range allusers {
		c_username, c_password := v.Username, v.Password

		// if username was found
		if user.Username == c_username {
			// if password is correct
			if user.Password == c_password {
				c.JSON(http.StatusOK, gin.H{
					"msg": "Login successful. Redirecting to home page",
				})
				fmt.Printf("User %s has logged in\n", user.Username)
				return
			} else {
				c.JSON(http.StatusOK, gin.H{
					"msg": "Username or password incorrect",
				})
				fmt.Printf("User %s has attempted to log in\n", user.Username)
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Username or password incorrect",
	})

	fmt.Println("A user has tried to login in with an unknown account")

}

func Getfriends(c *gin.Context) {
	var friends []string
	fmt.Println("Getting friends")

	friendsDetails := database.Getusers()
	for _, v := range friendsDetails {
		friends = append(friends, v.Username)
	}
	c.JSON(http.StatusOK, friends)
}
