package main

import (
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

func main() {
	origin := "https://api.hitbtc.com"
	url := "wss://api.hitbtc.com/api/3/ws/public"
	conn, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	var msg = map[string]interface{}{
		"method": "subscribe",
		"ch": "ticker/1s",
		"params": map[string]interface{}{
			"symbols": []string{"ETHBTC","BTCUSDT"},
		},
		"id": 123,
	}
	if err := websocket.JSON.Send(conn, msg); err != nil {
		log.Fatal(err)
	}

	var data interface{}
	if err = websocket.JSON.Receive(conn, &data); err != nil {
		log.Fatal(err)
	}

	fmt.Println(data)
}
