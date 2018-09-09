package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/lib/pq"
	"github.com/tmcnicol/inframe/wss"
)

func main() {
	createNotificationFunction()
	subscribeToUserNotifications()

	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	connStr := "user=thomas dbname=test sslmode=disable"
	listener := pq.NewListener(connStr, 10*time.Second, time.Minute, reportProblem)
	err := listener.Listen("events")
	if err != nil {
		panic(err)
	}

	fmt.Println("Start monitoring PostgreSQL...")

	s := wss.NewWSServer()
	go s.StartWSServer()

	for {
		notificationListener(listener, s)
	}
}

func createNotificationFunction() {
	filepath := "./queries/notification_function.sql"
	file, _ := ioutil.ReadFile(filepath) // TODO Handle error

	query := string(file)

	exectureQuery(query)

}

func subscribeToUserNotifications() {
	filepath := "./queries/subscribe.sql"
	file, _ := ioutil.ReadFile(filepath) // TODO Handle error
	query := string(file)

	exectureQuery(query)

}

func exectureQuery(query string) {
	connStr := "user=thomas dbname=test sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error during query:", err)
		return
	}

	fmt.Println(rows.Columns())
}

func notificationListener(l *pq.Listener, s *wss.WSServer) {
	for {
		select {
		case n := <-l.Notify:
			fmt.Println("Received update", n.Extra)
			s.Broadcast <- []byte(n.Extra)
		}
	}
}
