package db

import (
	"fmt"
	"time"

	"github.com/lib/pq"
)

type Subscription struct {
	notification chan []byte
	listener     *pq.Listener
	db           *DB
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
	for {
		select {
		case n := <-s.listener.Notify:
			fmt.Println(n)
		}
	}
}
