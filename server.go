package main

import (
	"log"

	"time"
	"os"
	"os/signal"
	"syscall"

	"regexp"
	"net/http"
	"encoding/json"

	"golang.org/x/net/websocket"
)

const (
	ORIGIN = "https://api.hitbtc.com"
	URL    = "wss://api.hitbtc.com/api/3/ws/public"
	SYNC_TIME = 1
)

var SYMBOLS = map[string]map[string]interface{}{
	"BTCUSDT": map[string]interface{}{
		"id": "BTC",
		"fullName": "BITCOIN",
		"feeCurrency": "USD",
	},
	"ETHBTC": map[string]interface{}{
		"id": "ETH",
		"fullName": "Ethereum",
		"feeCurrency": "BTC",
	},
}
var SYMBOLS_KEY_VALUE_MAPPING = map[string]string{
	"ask": "a",
	"bid": "b",
	"last": "c",
	"open": "o",
	"low": "l",
	"high": "h",
}

var pattern = regexp.MustCompile("^/currency/(all|[^/]+)$")

func handler(w http.ResponseWriter, r *http.Request) {
	matches := pattern.FindStringSubmatch(r.URL.Path)
	if len(matches) == 0 {
		http.NotFound(w, r)
		return
	}

	var body interface{}
	if data, ok := SYMBOLS[matches[1]]; ok {
		body = data
	} else if (matches[1] == "all") {
		v := make([]interface{}, 0, len(SYMBOLS))
		for  _, value := range SYMBOLS {
			v = append(v, value)
		}
		body = v
	} else {
		http.NotFound(w, r)
		return
	}

	respData, err := json.Marshal(body)
	if err != nil {
		// handle error
		log.Warn(err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(respData)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func GetCryptoDataThroughSocket() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(SYNC_TIME * time.Second)

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

	for {
		select {
			case <-ticker.C:
				var message map[string]interface{}
				if err = websocket.JSON.Receive(c, &message); err != nil {
					log.Println("read:", err)
					return
				}
				// log.Printf("recv: %s", message)
				
				if data, ok := message["data"]; ok {
					if message, ok = data.(map[string]interface{}); ok {
						for k := range message {
							if v, ok := message[k].(map[string]interface{}); ok {
								for m, n := range SYMBOLS_KEY_VALUE_MAPPING {
									SYMBOLS[k][m] = v[n]
								}
							}
						}
					}
				}
			case <-interrupt:
				log.Println("interrupt")

				// Cleaning
				// close the connection by sending a close message to server if required
				// waiting (with timeout) for the server to close the connection.
				ticker.Stop()

				<-time.After(time.Second)
				return
		}
	}
}

func main() {
	http.HandleFunc("/currency", NotFoundHandler)
	http.HandleFunc("/currency/", handler)

	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	GetCryptoDataThroughSocket()
}
