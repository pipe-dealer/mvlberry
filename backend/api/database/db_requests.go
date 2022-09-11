package database

type Request struct {
	Req_id int `json:"req_id"`
	Id     int `json:"id"`
	R_id   int `json:"r_id"`
}

type ReceivedRequest struct {
	C_username   string `json:"c_username"`
	Req_username string `json:"req_username"`
}

type RequestForm struct {
	Accepted bool `json:"accepted"`
	Req_id   int  `json:"req_id"`
}

// get request details by request ID
func GetrequestbyID(req_id int) Request {
	var id int
	var r_id int
	sqlGet := "SELECT id,r_id FROM requests WHERE req_id = $1"
	row := Db.QueryRow(sqlGet, req_id)
	row.Scan(&id, &r_id)

	return Request{
		Req_id: req_id,
		Id:     id,
		R_id:   r_id,
	}

}

// adds request to database
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

// check if request already exists
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

// get requests made to user
func Getincomingrequests(r_id int) []Request {
	var allRequests []Request
	sqlGet := "SELECT req_id,id FROM requests WHERE r_id = $1"
	rows, err := Db.Query(sqlGet, r_id)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var req_id int
		var c_id int
		if err := rows.Scan(&req_id, &c_id); err != nil {
			panic(err)
		}
		req := Request{
			Req_id: req_id,
			Id:     c_id,
		}
		allRequests = append(allRequests, req)
	}

	return allRequests
}

func Deleterequest(req_id int) int {
	_, err := Db.Exec("DELETE FROM requests WHERE req_id = $1", req_id)
	if err != nil {
		return 1
	} else {
		return 0
	}
}
