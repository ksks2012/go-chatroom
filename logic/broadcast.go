package logic

import "log"

const (
	MessageQueueLen = 8
)

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
	messageChannel:  make(chan *Message, MessageQueueLen),

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
		case user := <-b.leavingChannel:
			// user leaves
			delete(b.users, user.NickName)
			// Avoid goroutine leaks
			user.CloseMessageChannel()

			b.sendUserList()
		case msg := <-b.messageChannel:
			// send message to all users
			for _, user := range b.users {
				if user.UID == msg.User.UID {
					continue
				}
				user.MessageChannel <- msg
			}
		case nickname := <-b.checkUserChannel:
			if _, ok := b.users[nickname]; ok {
				b.checkUserCanInChannel <- false
			} else {
				b.checkUserCanInChannel <- true
			}
		}
	}
}

func (b *broadcaster) EnteringChannel() chan<- *User {
	return b.enteringChannel
}

func (b *broadcaster) LeavingChannel() chan<- *User {
	return b.leavingChannel
}

func (b *broadcaster) MessageChannel() chan<- *Message {
	return b.messageChannel
}

func (b *broadcaster) CheckUserChannel() chan<- string {
	return b.checkUserChannel
}

func (b *broadcaster) CheckUserCanInChannel() <-chan bool {
	return b.checkUserCanInChannel
}

func (b *broadcaster) sendUserList() {
	// To avoid deadlock, there is the possibility that the list that the user sees is not updated in time
	if len(b.messageChannel) < MessageQueueLen {
		b.messageChannel <- NewUserListMessage(b.users)
	} else {
		log.Println("The concurrency of messages is too large, causing MessageChannel congestion. . .")
	}
}
