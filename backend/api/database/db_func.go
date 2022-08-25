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
func Getuser(u_id int) User {
	sqlGet := "SELECT * FROM users WHERE id = $1"
	row := Db.QueryRow(sqlGet, u_id)

	var id int
	var username string
	var password string

	//assigns query data to each variable initialised above
	if err := row.Scan(&id, &username, &password); err != nil {
		panic(err)
	}
	user := User{
		Id:       id,
		Username: username,
		Password: password,
	}

	return user

}

//gets friends of user
func Getfriends(u_id string) []Friends {
	var allFriends []Friends

	sqlGet := "SELECT * FROM friendships WHERE id in (SELECT users.id FROM users WHERE users.username = $1)"

	rows, err := Db.Query(sqlGet, u_id)
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
