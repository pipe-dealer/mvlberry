package ws

import (
	"fmt"
	"time"
)

// ChatChannel acts as a hub where clients in the hub can only interact with other clients in the hub
// chan are channels where values are sent to and retrieved from
// by using this, concurrent goroutines can be used to handle multiple ChatChannels simultaneously with greater efficiency
type Chat struct {
	Id      int
	Joined  chan *Client     //channel can only store clients
	Left    chan *Client     //channel can only store clients
	Clients map[*Client]bool //channel map of client to a boolean
	Send    chan Message     //channel stores message received
}

var Allchannels []*Chat

// creates a new empty ChatChannel and adds to existsing channels
func Newchannel(f_id int) *Chat {
	newChannel := &Chat{
		Id:      f_id,
		Joined:  make(chan *Client),
		Left:    make(chan *Client),
		Clients: make(map[*Client]bool),
		Send:    make(chan Message),
	}

	Allchannels = append(Allchannels, newChannel)
	return newChannel

}

func removeChat(i int) {
	Allchannels[i] = Allchannels[len(Allchannels)-1]
	Allchannels = Allchannels[:len(Allchannels)-1]
}

// this will be a gorountine that runs indefinately, listening for values received from channels
func (chatchannel *Chat) Start() {
	for {

		select {
		//if a new user joins the ChatChannel, add it to list of clients and informs all other clients that a new user has joined
		//if there is a value in chatchannel.Joined, retrieve it and store in variable client

		case client := <-chatchannel.Joined:
			currentTime := time.Now()

			chatchannel.Clients[client] = true
			sysMsg := Message{
				MessageType: 1,
				Msg:         []byte(fmt.Sprintf("[%v]	%v has joined", currentTime.Format(time.ANSIC), client.Username)),
			}
			fmt.Printf("%v has joined channel %v \n", client.Id, chatchannel.Id)
			//sends sysMsg to all clients in Chat
			for client := range chatchannel.Clients {
				client.Conn.WriteMessage(sysMsg.MessageType, sysMsg.Msg)
			}

		//if a client leaves the Chat
		case client := <-chatchannel.Left:
			currentTime := time.Now()

			delete(chatchannel.Clients, client)
			sysMsg := Message{
				MessageType: 1,
				Msg:         []byte(fmt.Sprintf("[%v]	%v has left", currentTime.Format(time.ANSIC), client.Username)),
			}
			fmt.Printf("%v has left channel %v \n", client.Id, chatchannel.Id)
			//if there are no clients in chat, delete chat channel
			if len(chatchannel.Clients) == 0 {
				for i, v := range Allchannels {
					if v == chatchannel {
						removeChat(i)
						break
					}
				}
				//otherwise inform all other clients that user has left
			} else {
				for client := range chatchannel.Clients {
					client.Conn.WriteMessage(sysMsg.MessageType, sysMsg.Msg)
				}
			}
			//prints number of channels
			fmt.Println("Number of chats:", len(Allchannels))

		//get message sent from a client and sends it to all clients in Chat
		case msg := <-chatchannel.Send:
			for client := range chatchannel.Clients {
				if err := client.Conn.WriteMessage(msg.MessageType, msg.Msg); err != nil {
					fmt.Println("Could not send message to client", client.Username)
				}
			}

		}
	}
}
