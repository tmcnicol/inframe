package db

import (
	"fmt"
	"time"

	"github.com/lib/pq"
)

type Subscription struct {
	notification chan []byte
	quitc        chan struct{}

	// TODO: check how these listeners work there is two ways
	// 1. Listener only listens to topic
	// 2. Listens to multiple topics.
	//
	// I assume it is 2 which means that this implementation
	// it gives us a notification for each registered listener.
	// In that case I need to determine wether I want fan out
	// to the clients to be delegated to the Db or handled internally
	// here.
	listener *pq.Listener
	db       *DB
}

func (db *DB) NewSubscription() (*Subscription, error) {
	connStr := db.getConnectionString()
	listener := pq.NewListener(connStr, 10*time.Second, time.Minute, reportError)

	err := listener.Listen("events") // TODO make this work with other events?
	if err != nil {
		panic(err)
	}

	sub := &Subscription{
		notification: make(chan []byte),
		listener:     listener,
		db:           db,
	}

	err = sub.initialiseSubscriptions()
	if err != nil {
		return nil, err
	}

	sub.subscribeToChanges()

	return sub, nil
}

func (s *Subscription) initialiseSubscriptions() error {
	err := s.createNotificationFunction()
	return err
}

func reportError(ev pq.ListenerEventType, err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (s *Subscription) createNotificationFunction() error {
	filepath := "db/queries/notification_function.sql"
	_, err := s.db.ExecuteFromFile(filepath)

	return err
}

func (s *Subscription) subscribeToChanges() error {
	filepath := "db/queries/subscribe.sql"
	_, err := s.db.ExecuteFromFile(filepath)
	return err
}

func (s *Subscription) Run() {
	// go func() {
	for {
		select {
		case n := <-s.listener.Notify:
			fmt.Println(n)
		case <-s.quitc:
			return
		}
	}
	// }()
}
