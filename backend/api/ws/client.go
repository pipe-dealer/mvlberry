package ws

import (
	"github.com/gorilla/websocket"
)

// Client datatype, contains information about user, websocket connection and channel it is connected to
type Client struct {
	// Id      int
	Conn        *websocket.Conn
	ChatChannel *ChatChannel
}

// message datatype
type Message struct {
	MessageType int
	Msg         []byte
}

// client function that runs until client closes the websocket connection
func (c *Client) Read() {
	defer func() {
		c.ChatChannel.Left <- c
		c.Conn.Close()
	}()

	for {
		//get message sent from client
		msgType, msg, err := c.Conn.ReadMessage()
		if err != nil {
			panic(err)
		}
		message := Message{
			MessageType: msgType,
			Msg:         msg,
		}
		//passes it to send channel, where message will be processed and sent to all other clients
		c.ChatChannel.Send <- message
	}
}
