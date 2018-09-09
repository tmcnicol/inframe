package wss

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WSServer struct {
	clients   []*client
	Broadcast chan []byte
	register  chan *client
}

func NewWSServer() *WSServer {
	return &WSServer{
		clients:   []*client{},
		register:  make(chan *client),
		Broadcast: make(chan []byte),
	}
}

func (w *WSServer) run() {
	for {
		select {
		case client := <-w.register:
			fmt.Println("New client registered")
			w.clients = append(w.clients, client)

		case message := <-w.Broadcast:
			for _, client := range w.clients {
				client.send <- message
			}
		}
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "wss/home.html")
}

var upgrader = websocket.Upgrader{
	WriteBufferSize: 1024,
}

func serveWebsocket(s *WSServer, w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(conn)
	s.register <- client

	go client.writePump()
}

func (s *WSServer) StartWSServer() {
	go s.run()

	// Route requests
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWebsocket(s, w, r)
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
