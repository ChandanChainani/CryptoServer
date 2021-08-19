package main

import (
	"fmt"

	"log"
	"os"
	"os/signal"
	"time"

	"golang.org/x/net/websocket"
)

const (
	ORIGIN = "https://api.hitbtc.com"
	URL    = "wss://api.hitbtc.com/api/3/ws/public"
)

func main() {
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	log.Printf("connecting to %s", URL)

	// c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	c, err := websocket.Dial(URL, "", ORIGIN)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	var msg = map[string]interface{}{
		"method": "subscribe",
		"ch": "ticker/1s",
		"params": map[string]interface{}{
			"symbols": []string{"ETHBTC","BTCUSDT"},
		},
		"id": 123,
	}
	if err := websocket.JSON.Send(c, msg); err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			var message interface{}
			if err = websocket.JSON.Receive(c, &message); err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	for {
		select {
			case <-done:
				return
			case <-interrupt:
				log.Println("interrupt")

				// Cleaning
				// close the connection by sending a close message to server if required
				// waiting (with timeout) for the server to close the connection.

				select {
					case <-done:
					case <-time.After(time.Second):
				}
				return
		}
	}

	fmt.Println("completed")
}
