package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var sendConnCount = make(chan bool)
var upgrader = websocket.Upgrader{}

// Message object
type Message struct {
	Type string `json:"type"`
	Msg  json.RawMessage
}

// Chat object
type Chat struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

// ConnCount object
type ConnCount struct {
	Count int `json:"count"`
}

func main() {
	fs := http.FileServer(http.Dir(filepath.Join("..", "public")))
	http.Handle("/", fs)
	//TODO: Look up PingHandler and PongHandler and CloseMessage
	http.HandleFunc("/ws", handleConnections)
	// Start listening for messages
	go handleMessages()
	go sendNumberOfConnections()

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
	sendConnCount <- true

	for {
		var chat Chat
		var msg Message
		if err = ws.ReadJSON(&chat); err != nil {
			log.Println("error: ", err)
			delete(clients, ws)
			sendConnCount <- true
			break
		}
		msg.Type = "chat"
		msg.Msg, _ = json.Marshal(chat)
		if err != nil {
			log.Fatal("Error marshalling chat message: ", err)
		}
		broadcast <- msg
	}
}

func sendNumberOfConnections() {
	for {
		<-sendConnCount
		conCount := ConnCount{}
		msg := Message{}
		msg.Type = "count"
		conCount.Count = len(clients)
		msg.Msg, _ = json.Marshal(conCount)
		for client := range clients {
			if err := client.WriteJSON(msg); err != nil {
				// we won't close clients here
				log.Println("error writing message: ", err)
			}
		}
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
				sendConnCount <- true
			}
		}
	}
}
