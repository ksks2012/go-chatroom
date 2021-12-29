package logic

import "time"

// Message sent to the user
type Message struct {
	// Which user sent the message
	User    *User     `json:"user"`
	Type    int       `json:"type"`
	Content string    `json:"content"`
	MsgTime time.Time `json:"msg_time"`

	ToUser  string   `json:"to_user"`
	AtsUser []string `json:"ats_user"`

	Users map[string]*User `json:"users"`
}

const (
	MsgTypeNormal   = iota // normal user message
	MsgTypeSystem          // System message
	MsgTypeError           // error message
	MsgTypeUserList        // Send the current user list
)

func NewMessage(user *User, content string) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeNormal,
		Content: content,
		MsgTime: time.Now(),
	}
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

func NewUserListMessage(users map[string]*User) *Message {
	return &Message{
		User:    SystemUser,
		Type:    MsgTypeUserList,
		MsgTime: time.Now(),
		Users:   users,
	}
}
