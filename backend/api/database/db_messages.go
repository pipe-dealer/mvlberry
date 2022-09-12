package database

import (
	"fmt"
	"time"
)

// message record in table
type Message struct {
	Msg_id      int       //message id
	Sender_id   int       //id of user who sent the message
	Receiver_id int       //id of user the message was sent to
	Msg_text    string    //message content
	Date_sent   time.Time //timestamp when message was sent
}

// adds sent message to messages table
func Addmessage(message Message) {
	sqlInsert := "INSERT INTO messages(sender_id,receiver_id,msg_text,date_sent) VALUES($1,$2,$3,$4)"
	Db.Exec(sqlInsert, message.Sender_id, message.Receiver_id, message.Msg_text, message.Date_sent)
}

// get messages between two users
func Getmessages(c_id, r_id int) []Message {
	var allmessages []Message
	//select all records in messages that contains users' id
	sqlGet := "SELECT * FROM messages WHERE $1 in (sender_id,receiver_id) AND $2 in (sender_id,receiver_id) ORDER BY date_sent ASC"

	rows, err := Db.Query(sqlGet, c_id, r_id)
	if err != nil {
		fmt.Println(err)
	}

	//gets each record's data and stores it in a Message struct then adds it to allmessages
	for rows.Next() {
		var msg_id int
		var sender_id int
		var receiver_id int
		var msg_text string
		var date_sent time.Time

		//gets values from record and stores them in each variable
		if err := rows.Scan(&msg_id, &sender_id, &receiver_id, &msg_text, &date_sent); err != nil {
			fmt.Println(err)
		}

		//create Message struct with record details and add it allmessages slice
		message := Message{
			Msg_id:      msg_id,
			Sender_id:   sender_id,
			Receiver_id: receiver_id,
			Msg_text:    msg_text,
			Date_sent:   date_sent,
		}

		allmessages = append(allmessages, message)
	}
	return allmessages
}
