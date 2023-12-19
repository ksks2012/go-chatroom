package main

import (
	"fmt"
	"net/http"

	_ "net/http/pprof"

	"github.com/go-chatroom/global"
	log "github.com/go-chatroom/pkg/logger"
	"github.com/go-chatroom/server"
)

var (
	addr   = ":2022"
	banner = `
    ____              _____
   |    |    |   /\     |
   |    |____|  /  \    |
   |    |    | /----\   |
   |____|    |/      \  |

GChat Room：%s
`
)

func init() {
	global.Init()
}

func main() {
	fmt.Printf(banner+"\n", addr)

	server.RegisterHandle()

	log.Logger.Err(http.ListenAndServe(addr, nil))
}
