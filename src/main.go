package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{}

// Message object
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	fs := http.FileServer(http.Dir(filepath.Join("..", "public")))
	http.Handle("/", fs)
	http.HandleFunc("/ws", handleConnections)
	// Start listening for messages
	go handleMessages()

	log.Println("starting to listen on port :8000")
	if err := http.ListenAndServe("0.0.0.0:8000", nil); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Error upgrading connection: ", err)
	}
	defer ws.Close()
	clients[ws] = true

	for {
		var msg Message
		if err := ws.ReadJSON(&msg); err != nil {
			log.Println("error: ", err)
			delete(clients, ws)
			break
		}
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			if err := client.WriteJSON(msg); err != nil {
				log.Println("error writing message: ", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
