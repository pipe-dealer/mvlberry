package database

//User datatype
type User struct {
	Id       int    //user id
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
func GetuserbyID(c_id int) (User, int) {
	sqlGet := "SELECT * FROM users WHERE id = $1"
	row := Db.QueryRow(sqlGet, c_id)

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
func Getuserbyusername(c_username string) (User, int) {
	sqlGet := "SELECT * FROM users WHERE username = $1"
	row := Db.QueryRow(sqlGet, c_username)

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
