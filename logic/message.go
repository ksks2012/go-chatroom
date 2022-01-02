package logic

import (
	"time"

	"github.com/spf13/cast"
)

// Message sent to the user
type Message struct {
	// Which user sent the message
	User    *User     `json:"user"`
	Type    int       `json:"type"`
	Content string    `json:"content"`
	MsgTime time.Time `json:"msg_time"`

	ClientSendTime time.Time `json:"client_send_time"`

	ToUser  string   `json:"to_user"`
	AtsUser []string `json:"ats_user"`

	Users []*User `json:"users"`
}

const (
	MsgTypeNormal   = iota // normal user message
	MsgTypeSystem          // System message
	MsgTypeError           // error message
	MsgTypeUserList        // Send the current user list
)

func NewMessage(user *User, content string, clientTime string) *Message {
	message := &Message{
		User:    user,
		Type:    MsgTypeNormal,
		Content: content,
		MsgTime: time.Now(),
	}
	if clientTime != "" {
		message.ClientSendTime = time.Unix(0, cast.ToInt64(clientTime))
	}
	return message
}

func NewWelcomeMessage(user *User) *Message {
	return &Message{
		User:    SystemUser,
		Type:    MsgTypeSystem,
		Content: user.NickName + " Hello, welcome to the chat room!",
		MsgTime: time.Now(),
	}
}

func NewNoticeMessage(content string) *Message {
	return &Message{
		User:    SystemUser,
		Type:    MsgTypeSystem,
		Content: content,
		MsgTime: time.Now(),
	}
}

func NewErrorMessage(content string) *Message {
	return &Message{
		User:    SystemUser,
		Type:    MsgTypeSystem,
		Content: content,
		MsgTime: time.Now(),
	}
}

func NewUserListMessage(users []*User) *Message {
	return &Message{
		User:    SystemUser,
		Type:    MsgTypeUserList,
		MsgTime: time.Now(),
		Users:   users,
	}
}
