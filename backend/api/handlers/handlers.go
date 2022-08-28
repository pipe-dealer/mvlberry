package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mvlberry/backend/api/database"
	"github.com/mvlberry/backend/api/ws"
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
	var allfriends []string
	username := c.Query("user")
	friends := database.Getfriends(username)
	fmt.Println("Getting friends of user ", username)
	//get friend's username add it to array
	for _, v := range friends {
		f_user := database.GetuserByID(v.F_id).Username
		allfriends = append(allfriends, f_user)
	}
	//send friends' usernames back to client
	c.JSON(http.StatusOK, allfriends)

}

func Startws(c *gin.Context) {
	// client starts ws connection with another friend, sends query with their username and friend's username
	var chatId int

	user1 := c.Query("user1")
	user2 := c.Query("user2")

	user1_id := database.GetuserByUsername(user1).Id
	user2_id := database.GetuserByUsername(user2).Id

	user1_friends := database.Getfriends(user1)

	// get friendship id that corresponds between these two friends
	//assign friendship id to chat id
	for _, v := range user1_friends {
		if v.F_id == user2_id {
			chatId = v.Fs_id
		}
	}
	//initiate websocket connection
	ws.Startws(user1, user1_id, chatId, c)
}
