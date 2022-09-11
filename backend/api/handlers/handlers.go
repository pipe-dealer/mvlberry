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
	case 0: //signup successful
		c.JSON(http.StatusOK, gin.H{
			"msg": "Account successfully created. Redirecting to login page",
		})
		fmt.Printf("User %s was successfully created\n", newUser.Username)
	case 1: //username exists
		c.JSON(http.StatusOK, gin.H{
			"msg": "Username already exists",
		})
	case 2: //overall error
		c.JSON(http.StatusOK, gin.H{
			"msg": "Could not create account",
		})

	}
}

// login handler
func Login(c *gin.Context) {
	var user database.User
	//get login details
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Could not parse JSON",
		})
	}
	//get all users and checks if entered details exist
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
			} else { //incorrect password
				c.JSON(http.StatusOK, gin.H{
					"msg": "Username or password incorrect",
				})
				fmt.Printf("User %s has attempted to log in\n", user.Username)
				return
			}
		}
	}
	//incorrect username
	c.JSON(http.StatusOK, gin.H{
		"msg": "Username or password incorrect",
	})
}

// get friends of user
func GetFriends(c *gin.Context) {
	var allfriends []string
	username := c.Query("current_user")
	friends := database.Getfriends(username)

	//get friend's username add it to array
	for _, v := range friends {
		f_user, _ := database.GetuserbyID(v.F_id)
		f_username := f_user.Username
		allfriends = append(allfriends, f_username)
	}
	//send friends' usernames back to client
	c.JSON(http.StatusOK, allfriends)
}

// upgrade connection to websocket and start chat
func StartWS(c *gin.Context) {
	var chatId int
	// client starts ws connection with another friend, sends query with their username and friend's username
	user1 := c.Query("user1")
	user2 := c.Query("user2")

	//get users' id
	user1_details, _ := database.Getuserbyusername(user1)
	user2_details, _ := database.Getuserbyusername(user2)
	user1_id := user1_details.Id
	user2_id := user2_details.Id

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

// send friend request to requested user
func SendRequest(c *gin.Context) {
	//get current user and requested user
	var newRequest database.ReceivedRequest
	if err := c.ShouldBindJSON(&newRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Could not parse JSON",
		})

	}
	c_user, _ := database.Getuserbyusername(newRequest.C_username)
	c_id := c_user.Id

	//check if requested user exists
	req_user, err := database.Getuserbyusername(newRequest.Req_username)
	if err == 1 {
		c.JSON(http.StatusOK, gin.H{
			"msg": fmt.Sprintf("%v is not a valid username.", newRequest.Req_username),
		})
		return
	}
	req_id := req_user.Id

	//check if user is already friend of current user
	for _, v := range database.Getfriends(newRequest.C_username) {
		if v.F_id == req_id {
			c.JSON(http.StatusOK, gin.H{
				"msg": fmt.Sprintf("%v is already a friend", newRequest.Req_username),
			})
			return
		}
	}

	if res := database.Addrequest(c_id, req_id); res == 1 {
		c.JSON(http.StatusOK, gin.H{ //general error
			"msg": "Could not send request",
		})
	} else if res == 0 { //successful request
		c.JSON(http.StatusOK, gin.H{
			"msg": fmt.Sprintf("Request sent to %v", newRequest.Req_username),
		})
		fmt.Printf("%v sent a request to %v \n", newRequest.C_username, newRequest.Req_username)
	} else if res == 2 { //request already sent or received from same user
		c.JSON(http.StatusOK, gin.H{
			"msg": fmt.Sprintf("You have already sent a request to %v or %v has sent a request to you", newRequest.Req_username, newRequest.Req_username),
		})
	}
}

// get friend requests sent to user
func GetIncomingRequests(c *gin.Context) {
	c_username := c.Query("current_user")
	c_user, _ := database.Getuserbyusername(c_username)
	c_id := c_user.Id
	allRequests := database.Getincomingrequests(c_id)

	c.JSON(http.StatusOK, allRequests)
}

// accept or delete request
func HandleRequest(c *gin.Context) {
	//parse JSON data
	var requestform database.RequestForm
	if err := c.ShouldBindJSON(&requestform); err != nil {
		fmt.Println("Could not parse request form")
	}

	//if user decided to accept request, add
	if requestform.Accepted {
		request := database.GetrequestbyID(requestform.Req_id) //get request record
		switch err := database.Addfriend(request.Req_id, request.Id, request.R_id); err {
		case 0: //friend request was successfully accepted
			friend_username, _ := database.GetuserbyID(request.R_id)
			c.JSON(http.StatusOK, gin.H{
				"msg": fmt.Sprintf("%v has been added as a friend", friend_username.Username),
			})
		case 1: //error
			c.JSON(http.StatusOK, gin.H{
				"msg": "Could not add friend",
			})

		}
		//if user rejected request
	} else if !requestform.Accepted {
		//delete request record
		if err := database.Deleterequest(requestform.Req_id); err == 1 {
			c.JSON(http.StatusOK, gin.H{
				"msg": "Could not delete request",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"msg": "Request deleted",
			})

		}
	}
}
