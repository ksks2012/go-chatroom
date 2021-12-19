package main

import (
	"fmt"
	"log"
	"net/http"

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

func main() {
	fmt.Printf(banner+"\n", addr)

	server.RegisterHandle()

	log.Fatal(http.ListenAndServe(addr, nil))
}
