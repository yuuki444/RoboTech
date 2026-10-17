package handlers

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
	"robotics-project/database"
	"robotics-project/models"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan models.ChatMessage)
var mutex = &sync.Mutex{}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer ws.Close()

	mutex.Lock()
	clients[ws] = true
	mutex.Unlock()

	collection := database.MongoDatabase.Collection("chat_messages")
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}}).SetLimit(50)
	cursor, err := collection.Find(context.TODO(), bson.M{}, opts)
	if err == nil {
		var history []models.ChatMessage
		if cursor.All(context.TODO(), &history) == nil {
			for _, msg := range history {
				ws.WriteJSON(msg)
			}
		}
	}

	for {
		var msg models.ChatMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			mutex.Lock()
			delete(clients, ws)
			mutex.Unlock()
			break
		}

		msg.Timestamp = time.Now()

		collection.InsertOne(context.TODO(), msg)

		broadcast <- msg
	}
}

func HandleMessages() {
	for {
		msg := <-broadcast
		mutex.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}
