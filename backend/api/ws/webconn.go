package ws

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, //allow any connection connect to websocket enpoint
}

func Reader(conn *websocket.Conn) {
	for {
		// get message sent from client
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			panic(err)
		}

		fmt.Println("Message sent was: ", string(msg))
		//send it back to the client
		if err := conn.WriteMessage(msgType, msg); err != nil {
			panic(err)
		}
	}
}

func Startws(c *gin.Context) {
	//upgrade connection to websocket
	ws, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"msg": "Could not establish connection",
		})
	}
	Reader(ws)

}
