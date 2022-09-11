//sql commands for friendship table

package database

import "fmt"

type Friends struct {
	Id    int //user id
	F_id  int //friend id
	Fs_id int // friendship id
}

// gets friends of user
func Getfriends(c_username string) []Friends {
	var allFriends []Friends

	sqlGet := "SELECT * FROM friendships WHERE id in (SELECT users.id FROM users WHERE users.username = $1)"

	rows, err := Db.Query(sqlGet, c_username)
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

// create friendship
func Addfriend(req_id int, c_id int, r_id int) int {
	//create friendship id by getting number of rows in table and dividing by 2
	var fs_id int
	sqlGet := "SELECT COUNT(*) FROM friendships"
	row := Db.QueryRow(sqlGet)
	row.Scan(&fs_id)
	fs_id = (fs_id / 2) + 1

	//create two friendship rows
	sqlInsert := "INSERT INTO friendships(id,f_id,fs_id) VALUES ($1,$2,$3),($2,$1,$3)"
	if _, err := Db.Exec(sqlInsert, c_id, r_id, fs_id); err != nil {
		fmt.Println(c_id, r_id, fs_id)
		return 1 //could not create friendship
	} else {
		Db.Exec("DELETE FROM requests WHERE req_id = $1", req_id)
		return 0 //friendship created successfully
	}
}

func GetfriendshipbyID(friendship_id, c_id int) Friends {
	sqlGet := "SELECT * FROM friendships WHERE fs_id = $1 AND id = $2"
	row := Db.QueryRow(sqlGet, friendship_id, c_id)

	var id int
	var f_id int
	var fs_id int
	row.Scan(&id, &f_id, &fs_id)
	return Friends{
		Id:    id,
		F_id:  f_id,
		Fs_id: fs_id,
	}
}
