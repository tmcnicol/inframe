package main

import (
	"fmt"

	"github.com/tmcnicol/inframe/db"
	"github.com/tmcnicol/inframe/server"
)

func main() {
	db := db.NewDB("test", "thomas")
	sub, err := db.NewSubscription()

	if err != nil {
		fmt.Println("error creating subsctiption", err)
	}

	fmt.Println("Subscribing to notifications")
	go sub.Run()

	fmt.Println("Starting webserver")
	server := server.NewServer()
	go server.StartServer()

	for n := range sub.Notification {
		server.Broadcast <- n
	}
}
