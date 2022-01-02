package logic

import (
	"expvar"
	"fmt"
	"log"

	"github.com/go-chatroom/global"
)

func init() {
	expvar.Publish("message_queue", expvar.Func(calcMessageQueueLen))
}

func calcMessageQueueLen() interface{} {
	fmt.Println("===len=:", len(Broadcaster.messageChannel))
	return len(Broadcaster.messageChannel)
}

// logic/broadcast.go
// broadcaster
type broadcaster struct {
	// all users
	users map[string]*User

	// All channels are managed in a unified manner, which can avoid external misuse

	enteringChannel chan *User
	leavingChannel  chan *User
	messageChannel  chan *Message

	//Determine whether the user with the nickname can enter the chat room (duplicate or not): true => yes, false => no
	checkUserChannel      chan string
	checkUserCanInChannel chan bool
}

var Broadcaster = &broadcaster{
	users: make(map[string]*User),

	enteringChannel: make(chan *User),
	leavingChannel:  make(chan *User),
	messageChannel:  make(chan *Message, global.MessageQueueLen),

	checkUserChannel:      make(chan string),
	checkUserCanInChannel: make(chan bool),
}

// logic/broadcast.go

// Start starts the broadcaster
// needs to be run in a new goroutine because it will not return
func (b *broadcaster) Start() {
	for {
		select {
		case user := <-b.enteringChannel:
			// new user enters
			b.users[user.NickName] = user

			b.sendUserList()

			OfflineProcessor.Send(user)
		case user := <-b.leavingChannel:
			// user leaves
			delete(b.users, user.NickName)
			// Avoid goroutine leaks
			user.CloseMessageChannel()

			b.sendUserList()
		case msg := <-b.messageChannel:
			// send message to all users
			if msg.ToUser == "" {
				// Send a message to all online users
				for _, user := range b.users {
					if user.UID == msg.User.UID {
						continue
					}
					user.MessageChannel <- msg
				}
			} else {
				if user, ok := b.users[msg.ToUser]; ok {
					user.MessageChannel <- msg
				} else {
					// The other party is not online or the user does not exist, just ignore the message
					log.Println("user:", msg.ToUser, "not exists!")
				}
			}

			OfflineProcessor.Save(msg)
		case nickname := <-b.checkUserChannel:
			if _, ok := b.users[nickname]; ok {
				b.checkUserCanInChannel <- false
			} else {
				b.checkUserCanInChannel <- true
			}
		}
	}
}

func (b *broadcaster) UserEntering(u *User) {
	b.enteringChannel <- u
}

func (b *broadcaster) UserLeaving(u *User) {
	b.leavingChannel <- u
}

func (b *broadcaster) Broadcast(msg *Message) {
	b.messageChannel <- msg
}

func (b *broadcaster) CanEnterRoom(nickname string) bool {
	b.checkUserChannel <- nickname

	return <-b.checkUserCanInChannel
}

func (b *broadcaster) sendUserList() {
	// To avoid deadlock, there is the possibility that the list that the user sees is not updated in time
	if len(b.messageChannel) < global.MessageQueueLen {
		b.messageChannel <- NewUserListMessage(b.users)
	} else {
		log.Println("The concurrency of messages is too large, causing MessageChannel congestion. . .")
	}
}
