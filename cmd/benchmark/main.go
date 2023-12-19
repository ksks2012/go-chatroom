package main

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/go-chatroom/logic"
	log "github.com/go-chatroom/pkg/logger"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var (
	userNum       int           // number of users
	loginInterval time.Duration // User login time interval
	msgInterval   time.Duration // The interval for sending messages to the same user
)

func init() {
	flag.IntVar(&userNum, "u", 500, "Number of logged-in users")
	flag.DurationVar(&loginInterval, "l", 5e9, "User login time interval")
	flag.DurationVar(&msgInterval, "m", 1*time.Minute, "User sending message interval")
}

func main() {
	flag.Parse()

	for i := 0; i < userNum; i++ {
		go UserConnect("user" + strconv.Itoa(i))
		time.Sleep(loginInterval)
	}

	select {}
}

func UserConnect(nickname string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, "ws://127.0.0.1:2022/ws?nickname="+nickname, nil)
	if err != nil {
		log.Logger.Error().Msg(fmt.Sprintf("Dial error: %v", err))
		return
	}
	defer conn.Close(websocket.StatusInternalError, "Internal error!")

	go sendMessage(conn, nickname)

	ctx = context.Background()

	for {
		var message logic.Message
		err = wsjson.Read(ctx, conn, &message)
		if err != nil {
			log.Logger.Error().Msg(fmt.Sprintf("Receive msg error: %v", err))
			continue
		}

		if message.ClientSendTime.IsZero() {
			continue
		}
		if d := time.Now().Sub(message.ClientSendTime); d > 1*time.Second {
			fmt.Printf("Received server response (%d): %#v\n", d.Milliseconds(), message)
		}
	}

	conn.Close(websocket.StatusNormalClosure, "")
}

func sendMessage(conn *websocket.Conn, nickname string) {
	ctx := context.Background()
	i := 1
	for {
		msg := map[string]string{
			"content":   "Message from " + nickname + ":" + strconv.Itoa(i),
			"send_time": strconv.FormatInt(time.Now().UnixNano(), 10),
		}
		err := wsjson.Write(ctx, conn, msg)
		if err != nil {
			log.Logger.Error().Msg(fmt.Sprintf("send msg error: %v; nickname: %v; no: %v", err, nickname, i))
		}
		i++

		time.Sleep(msgInterval)
	}
}
