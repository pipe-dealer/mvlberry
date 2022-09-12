package ws

import (
	"time"

	"github.com/gorilla/websocket"
)

// Client datatype, contains information about user, websocket connection and channel it is connected to
type Client struct {
	Id       int             //client id
	Username string          //client username
	Conn     *websocket.Conn //pointer to websocket connection
	Chats    *Chat           //current chat it is connected to
}

// message datatype
type Message struct {
	MessageType     int //this will typically be 1
	Sender_id       int
	Sender_username string
	Msg             string //message is sent as bytes
	DateSent        time.Time
}

// client function that runs until client closes the websocket connection
func (c *Client) Read() {
	defer func() {
		c.Chats.Left <- c
		c.Conn.Close()
	}()

	for {
		//get message sent from client
		msgType, msg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		currentTime := time.Now()

		message := Message{
			MessageType:     msgType,
			Sender_id:       c.Id,
			Sender_username: c.Username,
			Msg:             string(msg),
			DateSent:        currentTime,
		}
		//passes it to send channel, where message will be processed and sent to all other clients
		c.Chats.Send <- message
	}
}
