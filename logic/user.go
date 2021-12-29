package logic

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var globalUID uint32 = 0

// logic/user.go
type User struct {
	UID            int           `json:"uid"`
	NickName       string        `json:"nickname"`
	EnterAt        time.Time     `json:"enter_at"`
	Addr           string        `json:"addr"`
	MessageChannel chan *Message `json:"-"`

	conn *websocket.Conn
}

var SystemUser = &User{}

func NewUser(conn *websocket.Conn, nickname string, addr string) *User {
	user := &User{
		NickName:       nickname,
		Addr:           addr,
		EnterAt:        time.Now(),
		MessageChannel: make(chan *Message, 32),

		conn: conn,
	}

	if user.UID == 0 {
		user.UID = int(atomic.AddUint32(&globalUID, 1))
	}

	return user
}

func (u *User) String() string {
	return "UID:" + strconv.Itoa(u.UID) + ";nickname:" + u.NickName + ";" +
		u.EnterAt.Format("2006-01-02 15:04:05 +8000") + " Enter chat room"
}

func (u *User) SendMessage(ctx context.Context) {
	for msg := range u.MessageChannel {
		wsjson.Write(ctx, u.conn, msg)
	}
}

func (u *User) CloseMessageChannel() {
	close(u.MessageChannel)
}

// logic/user.go
func (u *User) ReceiveMessage(ctx context.Context) error {
	var (
		receiveMsg map[string]string
		err        error
	)
	for {
		err = wsjson.Read(ctx, u.conn, &receiveMsg)
		if err != nil {
			// Determine whether the connection is closed, normally closed, not considered as an error
			var closeErr websocket.CloseError
			if errors.As(err, &closeErr) {
				return nil
			}

			return err
		}

		// Send content to chat room
		sendMsg := NewMessage(u, receiveMsg["content"])
		if strings.HasPrefix(sendMsg.Content, "@") {
			sendMsg.ToUser = strings.SplitN(sendMsg.Content, " ", 2)[0][1:]
		}
		reg := regexp.MustCompile(`@[^\s@]{2,20}`)
		sendMsg.AtsUser = reg.FindAllString(sendMsg.Content, -1)

		Broadcaster.Broadcast(sendMsg)
	}
}
