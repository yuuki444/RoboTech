package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"robotics-project/database"
	"robotics-project/middleware"
	"robotics-project/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func GetComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetType := vars["type"] // "product" or "video"
	targetID, _ := strconv.Atoi(vars["id"])

	collection := database.MongoDatabase.Collection("comments")
	cursor, err := collection.Find(context.TODO(), bson.M{"target_type": targetType, "target_id": targetID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var comments []models.Comment
	if err = cursor.All(context.TODO(), &comments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if comments == nil {
		comments = []models.Comment{}
	}
	json.NewEncoder(w).Encode(comments)
}

func AddComment(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middleware.UserIDKey).(int)
	vars := mux.Vars(r)
	targetType := vars["type"]
	targetID, _ := strconv.Atoi(vars["id"])

	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// getting username from Postgres
	var username string
	database.DB.QueryRow("SELECT username FROM users WHERE id = $1", userId).Scan(&username)

	comment.TargetType = targetType
	comment.TargetID = targetID
	comment.UserID = userId
	comment.Username = username
	comment.CreatedAt = time.Now()

	collection := database.MongoDatabase.Collection("comments")
	result, err := collection.InsertOne(context.TODO(), comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	comment.ID = result.InsertedID.(interface{}).(primitive.ObjectID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}
