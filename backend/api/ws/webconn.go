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

var count = 0

// adds client to channel if it exists, otherwise create a new channel
func serve(count int, ws *websocket.Conn) {
	var chatchannel *ChatChannel
	f := false
	//using a simple counter as channel ID, if channel ID exists, add client to that channel
	for _, v := range Allchannels {
		if count == v.Id {
			chatchannel = v
			f = true
			break
		}
	}
	//otherwise create a new channel with counter as channel ID
	if !f {
		chatchannel = Newchannel(count)
	}
	go chatchannel.Start()

	client := Client{
		Conn:        ws,
		ChatChannel: chatchannel,
	}

	fmt.Println("Number of chats:", len(Allchannels))

	chatchannel.Joined <- &client
	client.Read()

}

func Startws(c *gin.Context) {
	ws, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"msg": "Could not establish connection",
		})
	}
	//counter will alternate between 1 and 2, which determines which channel new clients will go into
	count += 1
	if count > 2 {
		count = 1
	}

	serve(count, ws)

}
