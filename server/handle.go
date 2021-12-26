package server

import (
	"net/http"

	"github.com/go-chatroom/logic"
)

var rootDir string

func RegisterHandle() {
	inferRootDir()

	// Processing broadcast message
	go logic.Broadcaster.Start()

	http.HandleFunc("/", homeHandleFunc)
	http.HandleFunc("/ws", WebSocketHandleFunc)
}
