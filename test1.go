package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Program Started")
	done := make(chan struct{})
	go func() {
		fmt.Println("Server Running")
		log.Fatal(http.ListenAndServe(":8080", nil))
		fmt.Println("Server Stopped")
		done <- struct{}{}
	}()
	fmt.Println("Server Started")
	<- done
	fmt.Println("Server Stopped")
	fmt.Println("Program Stopped")
}
