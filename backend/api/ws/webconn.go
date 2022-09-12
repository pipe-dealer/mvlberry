package ws

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// upgrader variable, which will be used to upgrade connection to websocket
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, //allow any connection connect to websocket enpoint
}

// adds client to channel if it exists, otherwise create a new channel
func serve(clientUsername string, clientId int, chatId int, ws *websocket.Conn) {
	var chat *Chat
	f := false
	//checks if channel ID exists, add client to that channel
	for _, v := range Allchats {
		if chatId == v.Id {
			chat = v
			f = true
			break
		}
	}
	//otherwise create a new channel with chatId
	if !f {
		chat = Newchannel(chatId)
	}

	//start up chat
	go chat.Start()

	//create client
	client := Client{
		Id:       clientId,
		Username: clientUsername,
		Conn:     ws,
		Chats:    chat,
	}

	fmt.Println("Number of chats:", len(Allchats))

	//add client to chat
	chat.Joined <- &client
	//check for actions from client
	client.Read()

}

func Startws(clientUsername string, clientId int, chatId int, c *gin.Context) {
	//upgrade connection to websocket
	ws, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"msg": "Could not establish connection",
		})
	}
	//add client to channel
	serve(clientUsername, clientId, chatId, ws)

}
