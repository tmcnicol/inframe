package main

import (
	"fmt"
	"time"

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

	for {
		// Game loop for want of a ctrl-c handler
		time.Sleep(time.Second)
	}
}
