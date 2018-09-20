package main

import (
	"fmt"

	"github.com/tmcnicol/inframe/db"
)

func main() {
	db := db.NewDB("test", "thomas")
	sub, err := db.NewSubscription()

	if err != nil {
		fmt.Println("error creating subsctiption", err)
	}

	sub.Run()

	// listener := pq.NewListener(connStr, 10*time.Second, time.Minute, reportProblem)
	// err := listener.Listen("events")
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Start monitoring PostgreSQL...")

	// s := wss.NewWSServer()
	// go s.StartWSServer()
}
