package server

import (
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	send chan []byte
	conn *websocket.Conn
}

func newClient(conn *websocket.Conn) *client {
	return &client{
		send: make(chan []byte),
		conn: conn,
	}
}

const writeWait = 10 * time.Second

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *client) writePump() {
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}
