package server

import (
	"log"
	"net/http"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	"github.com/go-chatroom/logic"
)

func WebSocketHandleFunc(w http.ResponseWriter, req *http.Request) {
	// Accept receives the WebSocket handshake from the client and upgrades the connection to WebSocket.
	// If the Origin domain is different from the host, Accept will refuse the handshake unless the InsecureSkipVerify option is set (set by the third parameter AcceptOptions).
	// In other words, by default, it does not allow cross-origin requests. If an error occurs, Accept will always write an appropriate response
	conn, err := websocket.Accept(w, req, nil)
	if err != nil {
		log.Println("websocket accept error:", err)
		return
	}

	// 1. A new user comes in and builds an instance of that user
	nickname := req.FormValue("nickname")
	logic.Broadcaster.CheckUserChannel() <- nickname
	if l := len(nickname); l < 2 || l > 20 {
		log.Println("nickname illegal: ", nickname)
		wsjson.Write(req.Context(), conn, logic.NewErrorMessage("Illegal nickname, nickname length: 4-20"))
		conn.Close(websocket.StatusUnsupportedData, "nickname illegal!")
		return
	}
	if !<-logic.Broadcaster.CheckUserCanInChannel() {
		log.Println("Nickname already exists:", nickname)
		wsjson.Write(req.Context(), conn, logic.NewErrorMessage("The nickname already exists!"))
		conn.Close(websocket.StatusUnsupportedData, "nickname exists!")
		return
	}

	user := logic.NewUser(conn, nickname, req.RemoteAddr)

	// 2. Turn on the goroutine that sends messages to users
	go user.SendMessage(req.Context())

	// 3. Send a welcome message to the current user
	user.MessageChannel <- logic.NewWelcomeMessage(nickname)
	log.Println("New User:", nickname)

	// Notify all users of the arrival of new users
	msg := logic.NewNoticeMessage(nickname + "Joined the chat room")
	logic.Broadcaster.MessageChannel() <- msg

	// 4. Add this user to the user list of the broadcaster
	logic.Broadcaster.EnteringChannel() <- user
	log.Println("user:", nickname, "joins chat")

	// 5. Receive user messages
	err = user.ReceiveMessage(req.Context())

	// 6. User left
	logic.Broadcaster.LeavingChannel() <- user
	msg = logic.NewNoticeMessage(user.NickName + " Left the chat room")
	logic.Broadcaster.MessageChannel() <- msg
	log.Println("user:", nickname, "leaves chat")

	// Execute different Close according to the error during reading
	if err == nil {
		conn.Close(websocket.StatusNormalClosure, "")
	} else {
		log.Println("read from client error:", err)
		conn.Close(websocket.StatusInternalError, "Read from client error")
	}
}
