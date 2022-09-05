package database

import "fmt"

//User datatype
type User struct {
	Id       int
	Username string `json:"username"`
	Password string `json:"password"`
}

type Friends struct {
	Id    int //user id
	F_id  int //friend id
	Fs_id int // friendship id
}

type Request struct {
	C_username   string `json:"c_username"`
	Req_username string `json:"req_username"`
}

//checks if the new username already exists
func checkDuplicates(username string) int {
	sqlSelect := "SELECT username FROM users"
	rows, err := Db.Query(sqlSelect)
	if err != nil {
		return 2
	}

	//iterates through every record and checks if username is the same as new username
	for rows.Next() {
		var currentusername string
		if err := rows.Scan(&currentusername); err != nil {
			panic(err)
		}
		if currentusername == username {
			return 1
		}
	}

	return 0
}

//Adds new user to database, and returns a response code, which will produce a different response message depending on the code
func Adduser(newuser User) int {
	username := newuser.Username
	password := newuser.Password

	//if username does not exist, add new user to database
	if checkDuplicates(username) == 0 {
		sqlInsert := "INSERT INTO users (username,password) VALUES ($1,$2)"

		if _, err := Db.Exec(sqlInsert, username, password); err != nil {
			panic(err)
		}
		return 0 //user successfully created
	} else if checkDuplicates(username) == 1 {
		return 1 //username already exists
	} else {
		return 2 //could not create user
	}

}

//gets all account details
func Getusers() []User {
	var allUsers []User //create array of User struct

	sqlGet := "SELECT * FROM users ORDER BY id ASC;" //select all records in users table

	rows, err := Db.Query(sqlGet)

	if err != nil {
		panic(err)
	}

	//gets id, username and password from each record and stores it into allUsers
	for rows.Next() {
		var id int
		var username string
		var password string

		if err := rows.Scan(&id, &username, &password); err != nil {
			panic(err)
		}

		allUsers = append(allUsers, User{
			Id:       id,
			Username: username,
			Password: password,
		})
	}

	return allUsers
}

//gets details of specific user
func GetuserbyID(u_id int) (User, int) {
	sqlGet := "SELECT * FROM users WHERE id = $1"
	row := Db.QueryRow(sqlGet, u_id)

	var id int
	var username string
	var password string

	//assigns query data to each variable initialised above
	if err := row.Scan(&id, &username, &password); err != nil {
		return User{}, 1
	}
	user := User{
		Id:       id,
		Username: username,
		Password: password,
	}

	return user, 0

}

//Gets user details from username
func Getuserbyusername(u_username string) (User, int) {
	sqlGet := "SELECT * FROM users WHERE username = $1"
	row := Db.QueryRow(sqlGet, u_username)

	var id int
	var username string
	var password string

	//assigns query data to each variable initialised above
	if err := row.Scan(&id, &username, &password); err != nil {
		return User{}, 1 //if username does not exist return error code 1
	}
	user := User{
		Id:       id,
		Username: username,
		Password: password,
	}

	return user, 0

}

//gets friends of user
func Getfriends(u_username string) []Friends {
	var allFriends []Friends

	sqlGet := "SELECT * FROM friendships WHERE id in (SELECT users.id FROM users WHERE users.username = $1)"

	rows, err := Db.Query(sqlGet, u_username)
	if err != nil {
		fmt.Println("error getting friends")
	}

	//loops through all queried records
	for rows.Next() {
		var id int
		var f_id int
		var fs_id int
		if err := rows.Scan(&id, &f_id, &fs_id); err != nil {
			panic(err)
		}
		//create Friend struct with queried data
		friend := Friends{
			Id:    id,
			F_id:  f_id,
			Fs_id: fs_id,
		}
		//add Friend to allFriends
		allFriends = append(allFriends, friend)
	}

	return allFriends
}

//adds request to database
func Addrequest(c_id int, req_id int) int {
	//add c_id and req_id to requests table
	if Checkrequest(c_id, req_id) {
		sqlInsert := "INSERT INTO requests(id,r_id) VALUES ($1,$2)"
		if _, err := Db.Exec(sqlInsert, c_id, req_id); err != nil {
			return 1 //could not send request
		} else {
			return 0 //request sent
		}
	} else {
		return 2 //request already exists
	}

}

//check if request already exists
func Checkrequest(c_id int, req_id int) bool {
	var reqCount int
	sqlSelect := "SELECT count(*) FROM requests WHERE $1 IN (id,r_id) AND $2 IN (id,r_id)"
	row := Db.QueryRow(sqlSelect, c_id, req_id)
	if err := row.Scan(&reqCount); err != nil {
		panic(err)
	}

	if reqCount == 0 {
		return true //no duplicate requests
	} else {
		return false //duplicate requests exist
	}

}

//get requests made to user
func Getincomingrequests(r_id int) []Request {
	var allRequests []Request
	sqlGet := "SELECT id FROM requests WHERE r_id = $1"
	rows, err := Db.Query(sqlGet, r_id)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var c_id int
		if err := rows.Scan(&c_id); err != nil {
			panic(err)
		}
		c_user, _ := GetuserbyID(c_id)
		c_username := c_user.Username
		req := Request{
			C_username: c_username,
		}
		allRequests = append(allRequests, req)
	}

	return allRequests
}
