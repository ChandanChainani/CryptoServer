package main

import (
	"fmt"
	"os"
	"log"
	// "io/ioutil"
	"net/http"
	"encoding/json"
	"path/filepath"
	"regexp"
	"reflect"
)

func GetConfig(name string) interface{} {
	cwd, err := os.Getwd()

	file, err := os.Open(filepath.Join(cwd, "config", name))
	if err != nil {
		log.Fatal("configuration file not found")
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var config interface{}

	err = decoder.Decode(&config)
	if err != nil {
		log.Println("error:", err)
	}

	return config
}

var SYMBOLS = map[string]interface{}{
	"BTCUSD": map[string]string{
		"fullName": "BITCOIN",
		"feeCurrency": "USD",
	},
	"ETHBTC": map[string]string{
		"fullName": "Ethereum",
		"feeCurrency": "BTC",
	},
}

var pattern = regexp.MustCompile("^/currency/(all|[^/]+)$")
func handler(w http.ResponseWriter, r *http.Request) {
	matches := pattern.FindStringSubmatch(r.URL.Path)
	// matchCount := len(matches)
	// fmt.Println(matchCount)
	fmt.Println(matches)

	if len(matches) == 0 {
		http.NotFound(w, r)
		return
	}

	url := "https://api.hitbtc.com/api/3/public/ticker"

	if _, ok := SYMBOLS[matches[1]]; ok {
		url += "/" + matches[1]
	}

	fmt.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Body)
	fmt.Println(reflect.TypeOf(resp.Body))

	defer resp.Body.Close()

	// var body map[string]interface{}
	var body interface{}
	json.NewDecoder(resp.Body).Decode(&body)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// w.WriteHeader(http.StatusCreated)
	// fmt.Println(body)
	jData, err := json.Marshal(body)
	if err != nil {
		// handle error
	}
	w.Write(jData)
	// json.NewEncoder(w).Encode(body)

  // fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
// https://api.hitbtc.com/api/3/public/currency/BTC
// https://api.hitbtc.com/api/3/public/ticker/ETHBTC

func main() {
	// resp, err := http.Get("https://api.hitbtc.com/api/3/public/ticker/ETHBTC")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer resp.Body.Close()

	http.HandleFunc("/currency", NotFoundHandler)
	http.HandleFunc("/currency/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
	fmt.Println("completed")
}
