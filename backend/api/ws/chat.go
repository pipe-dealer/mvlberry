package ws

import (
	"fmt"
	"time"

	"github.com/mvlberry/backend/api/database"
)

// ChatChannel acts as a hub where clients in the hub can only interact with other clients in the hub
// chan are channels where values are sent to and retrieved from
// by using this, concurrent goroutines can be used to handle multiple ChatChannels simultaneously with greater efficiency
type Chat struct {
	Id      int              //Chat id, this will the friendship id
	Joined  chan *Client     //channel can only store clients
	Left    chan *Client     //channel can only store clients
	Clients map[*Client]bool //channel map of client to a boolean
	Send    chan Message     //channel stores message received
}

var Allchats []*Chat

// creates a new empty chat and adds to existing chats
func Newchannel(f_id int) *Chat {
	newChat := &Chat{
		Id:      f_id,
		Joined:  make(chan *Client),
		Left:    make(chan *Client),
		Clients: make(map[*Client]bool),
		Send:    make(chan Message),
	}

	Allchats = append(Allchats, newChat)
	return newChat

}

// deletes chat
func removeChat(i int) {
	//replace ith element with last element
	Allchats[i] = Allchats[len(Allchats)-1]
	//new truncated slice
	Allchats = Allchats[:len(Allchats)-1]
}

// check if date of message sent was today or yesterday
func checkDate(msgDate time.Time) int {
	currentTime := time.Now()          //get current time
	cy, cm, cd := currentTime.Date()   //get year, month and day of current date
	msgY, msgM, msgD := msgDate.Date() //get year, month and day of message date

	yesterdayY, yesterdayM, yesterdayD := currentTime.AddDate(0, 0, -1).Date() //get year, month and day of yesterday

	if cy == msgY && cm == msgM && cd == msgD { //if current date matches message date
		return 0 //today
	} else if yesterdayY == msgY && yesterdayM == msgM && yesterdayD == msgD { //if yesterday's date matches message date
		return 1 //yesterday
	} else {
		return 2
	}
}

// this will be a gorountine which listens for values received from channels
func (chatchannel *Chat) Start() {
	for {

		select {
		//if a new user joins the ChatChannel, add it to list of clients and informs all other clients that a new user has joined
		//if there is a value in chatchannel.Joined, retrieve it and store in variable client

		case client := <-chatchannel.Joined:
			currentTime := time.Now()                                          //get time user joined chat
			friend_id := database.GetfriendshipbyID(chatchannel.Id, client.Id) //get id of friend
			allmessages := database.Getmessages(client.Id, friend_id.F_id)     //get chat history

			chatchannel.Clients[client] = true
			//create system message
			sysMsg := Message{
				MessageType: 1,
				Msg:         fmt.Sprintf("SYS: %v has joined", client.Username),
				DateSent:    currentTime,
			}
			fmt.Printf("[%v]  %v has joined channel %v \n", sysMsg.DateSent.Format("2006-01-02 15:04"), client.Id, chatchannel.Id)

			//get chat history and sends them to client
			for _, v := range allmessages {
				sender, _ := database.GetuserbyID(v.Sender_id) //get sender of message
				date := v.Date_sent.Format("2006-01-02 15:04") //get date message was sent

				if res := checkDate(v.Date_sent); res == 0 { //if message was sent today, format date accordingly
					date = fmt.Sprintf("Today at %v", v.Date_sent.Format("15:04"))
				} else if res == 1 { //if message was sent yesterday, format date accordingly
					date = fmt.Sprintf("Yesterday at %v", v.Date_sent.Format("15:04"))
				}
				//send message to client
				client.Conn.WriteMessage(1, []byte(fmt.Sprintf("[%v] %v: %v", date, sender.Username, v.Msg_text)))
			}

			//informs connected users that a new user has joined chat
			for client := range chatchannel.Clients {
				client.Conn.WriteMessage(sysMsg.MessageType, []byte(fmt.Sprintf("[%v] %v", fmt.Sprintf("Today at %v", currentTime.Format("15:04")), sysMsg.Msg)))
			}

		//if a client leaves the Chat
		case client := <-chatchannel.Left:
			currentTime := time.Now() //get time user left chat

			delete(chatchannel.Clients, client) //remove client from list of connected users
			//create system message
			sysMsg := Message{
				MessageType: 1,
				Msg:         fmt.Sprintf("%v has left", client.Username),
				DateSent:    currentTime,
			}
			fmt.Printf("[%v]  %v has left channel %v \n", sysMsg.DateSent.Format("2006-01-02 15:04"), client.Id, chatchannel.Id)
			//if there are no clients in chat, delete chat channel
			if len(chatchannel.Clients) == 0 {
				for i, v := range Allchats {
					if v == chatchannel {
						removeChat(i)
						break
					}
				}
				//otherwise inform all other clients that user has left
			} else {
				for client := range chatchannel.Clients {
					client.Conn.WriteMessage(sysMsg.MessageType, []byte(fmt.Sprintf("[%v] %v", fmt.Sprintf("Today at %v", time.Now().Format("15:04")), sysMsg.Msg)))
				}
			}
			//prints number of channels
			fmt.Println("Number of chats:", len(Allchats))

		//get message sent from a client and sends it to all clients in Chat
		case msg := <-chatchannel.Send:
			//sends message from client to all connected users, including sender of message
			for client := range chatchannel.Clients {
				client.Conn.WriteMessage(msg.MessageType, []byte(fmt.Sprintf("[%v] %v: %v", fmt.Sprintf("Today at %v", time.Now().Format("15:04")), msg.Sender_username, msg.Msg)))
			}
			friend_id := database.GetfriendshipbyID(chatchannel.Id, msg.Sender_id) //get id of user message was sent to
			//create Message struct with details of message sent
			newMsg := database.Message{
				Sender_id:   msg.Sender_id,
				Receiver_id: friend_id.F_id,
				Msg_text:    msg.Msg,
				Date_sent:   msg.DateSent,
			}
			//save message to chat history
			database.Addmessage(newMsg)

		}
	}
}
