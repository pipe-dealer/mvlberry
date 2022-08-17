package database

import "fmt"

//User datatype
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
		fmt.Println(username)
		return 0 //user successfully created
	} else if checkDuplicates(username) == 1 {
		return 1 //username already exists
	} else {
		return 2 //could not create user
	}

}
