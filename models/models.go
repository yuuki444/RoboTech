package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type Product struct {
	ID          int     `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Description string  `json:"description" db:"description"`
	Price       float64 `json:"price" db:"price"`
	ImageURL    string  `json:"image_url" db:"image_url"`
	Stock       int     `json:"stock" db:"stock"`
}

type CartItem struct {
	ID        int      `json:"id" db:"id"`
	UserID    int      `json:"user_id" db:"user_id"`
	ProductID int      `json:"product_id" db:"product_id"`
	Quantity  int      `json:"quantity" db:"quantity"`
	Product   *Product `json:"product,omitempty"`
}

type Video struct {
	ID          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	YouTubeID   string `json:"youtube_id" db:"youtube_id"`
	Description string `json:"description" db:"description"`
}

type Comment struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TargetType string             `json:"target_type" bson:"target_type"`
	TargetID   int                `json:"target_id" bson:"target_id"`
	UserID     int                `json:"user_id" db:"user_id" bson:"user_id"`
	Username   string             `json:"username" bson:"username"`
	Text       string             `json:"text" bson:"text"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
}

type ChatMessage struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username"`
	Message   string             `json:"message" bson:"message"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
}
