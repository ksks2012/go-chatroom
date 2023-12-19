package server

import (
	"fmt"
	"net/http"

	"github.com/go-chatroom/logic"
	log "github.com/go-chatroom/pkg/logger"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func WebSocketHandleFunc(w http.ResponseWriter, req *http.Request) {
	// Accept receives the WebSocket handshake from the client and upgrades the connection to WebSocket.
	// If the Origin domain is different from the host, Accept will refuse the handshake unless the InsecureSkipVerify option is set (set by the third parameter AcceptOptions).
	// In other words, by default, it does not allow cross-origin requests. If an error occurs, Accept will always write an appropriate response
	conn, err := websocket.Accept(w, req, nil)
	if err != nil {
		log.Logger.Error().Msg(fmt.Sprintf("Websocket accept error: %v", err))
		return
	}

	// 1. A new user comes in and builds an instance of that user
	token := req.FormValue("token")
	nickname := req.FormValue("nickname")
	// logic.Broadcaster.CheckUserChannel(nickname) <- nickname
	// if l := len(nickname); l < 2 || l > 20 {
	// 	log.Logger.Print("nickname illegal: ", nickname)
	// 	wsjson.Write(req.Context(), conn, logic.NewErrorMessage("Illegal nickname, nickname length: 4-20"))
	// 	conn.Close(websocket.StatusUnsupportedData, "nickname illegal!")
	// 	return
	// }
	if !logic.Broadcaster.CanEnterRoom(nickname) {
		log.Logger.Info().Msg(fmt.Sprintf("Nickname already exists: %v", nickname))
		wsjson.Write(req.Context(), conn, logic.NewErrorMessage("The nickname already exists!"))
		conn.Close(websocket.StatusUnsupportedData, "nickname exists!")
		return
	}

	user := logic.NewUser(conn, token, nickname, req.RemoteAddr)

	// 2. Turn on the goroutine that sends messages to users
	go user.SendMessage(req.Context())

	// 3. Send a welcome message to the current user
	user.MessageChannel <- logic.NewWelcomeMessage(user)
	log.Logger.Info().Msg(fmt.Sprintf("New User: %v", nickname))

	// Notify all users of the arrival of new users
	msg := logic.NewUserEnterMessage(user)
	logic.Broadcaster.Broadcast(msg)

	// 4. Add this user to the user list of the broadcaster
	logic.Broadcaster.UserEntering(user)
	log.Logger.Info().Msg(fmt.Sprintf("user: %v joins chat", nickname))

	// 5. Receive user messages
	err = user.ReceiveMessage(req.Context())

	// 6. User left
	logic.Broadcaster.UserLeaving(user)
	msg = logic.NewUserLeaveMessage(user)
	logic.Broadcaster.Broadcast(msg)
	log.Logger.Info().Msg(fmt.Sprintf("user: %v leaves chat", nickname))

	// Execute different Close according to the error during reading
	if err == nil {
		conn.Close(websocket.StatusNormalClosure, "")
	} else {
		log.Logger.Error().Msg(fmt.Sprintf("Read from client error: %v", err))
		conn.Close(websocket.StatusInternalError, "Read from client error")
	}
}
