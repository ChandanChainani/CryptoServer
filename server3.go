package main

import (
	"log"
	"context"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func main() {
	conn, _, _, err := ws.DefaultDialer.Dial(context.TODO(), "wss://api.hitbtc.com/api/3/public/ticker")
	if err != nil {
		log.Fatal("Connection Failed")
	}


	var message = []byte("")
	err = wsutil.WriteClientMessage(conn, ws.OpText, message)
	if err != nil {
		log.Fatal("Connection Failed")
	}

	msg, _, err := wsutil.ReadServerData(conn)
	if err != nil {
		log.Fatal("Connection Failed")
	}
}
