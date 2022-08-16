package database

//User datatype
type User struct {
	username string `json:"username"`
	password string `json:"password"`
}

//checks if the new username already exists
func checkDuplicates(username string) bool {

	sqlSelect := "SELECT username from users"
	rows, err := Db.Query(sqlSelect)
	if err != nil {
		panic(err)
	}

	//iterates through every record and checks if username is the same as new username
	for rows.Next() {
		var currentusername string
		if err := rows.Scan(&currentusername); err != nil {
			panic(err)
		}
		if currentusername == username {
			return true
		}
	}

	return false
}

//Adds new user to database, and returns a response code, which will produce a different response message depending on the code
func Adduser(newuser User) int {
	username := newuser.username
	password := newuser.password

	//if username does not exist, add new user to database
	if !checkDuplicates(username) {
		sqlInsert := "INSERT INTO users (username,password) VALUES ($1,$2)"

		if _, err := Db.Exec(sqlInsert, username, password); err != nil {
			panic(err)
		}

		return 0 //user successfully created
	} else if checkDuplicates(username) {
		return 1 //username already exists
	}

	return 2 //could not create user
}
