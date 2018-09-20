package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Struct to hold the web socket conenctions
// for convienience the server is attached to this
// as well but should probably spilt these out and compose them.
type Server struct {
	clients   []*client
	Broadcast chan []byte
	register  chan *client
}

func NewServer() *Server {
	return &Server{
		clients:   []*client{},
		register:  make(chan *client),
		Broadcast: make(chan []byte),
	}
}

func (w *Server) startWebSocketServer() {
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
	http.ServeFile(w, r, "server/home.html")
}

var upgrader = websocket.Upgrader{
	WriteBufferSize: 1024,
}

func serveWebsocket(s *Server, w http.ResponseWriter, r *http.Request) {
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

func (s *Server) StartServer() {
	go s.startWebSocketServer()

	fmt.Println("Serving at localhost:8080")
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
