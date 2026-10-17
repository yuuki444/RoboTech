package main

import (
	"log"
	"net/http"
	"os"
	"robotics-project/database"
	"robotics-project/handlers"
	"robotics-project/middleware"
	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	pgUrl := os.Getenv("POSTGRES_URL")
	if pgUrl == "" {
		pgUrl = "postgres://robotics_user:secure_password@localhost:5432/robotics_db?sslmode=disable"
	}

	mongoUri := os.Getenv("MONGO_URI")
	if mongoUri == "" {
		mongoUri = "mongodb://localhost:27017"
	}

	database.InitPostgres(pgUrl)
	database.InitMongo(mongoUri)

	r := mux.NewRouter()

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/"))).Methods("GET")

	r.HandleFunc("/api/register", handlers.Register).Methods("POST")
	r.HandleFunc("/api/login", handlers.Login).Methods("POST")
	

	s := r.PathPrefix("/api").Subrouter()
	s.Use(middleware.AuthMiddleware)
	s.HandleFunc("/profile", handlers.Profile).Methods("GET")

	r.HandleFunc("/api/products", handlers.GetProducts).Methods("GET")
	r.HandleFunc("/api/products", handlers.CreateProduct).Methods("POST")
	r.HandleFunc("/api/products/{id}", handlers.UpdateProduct).Methods("PUT")
	r.HandleFunc("/api/products/{id}", handlers.DeleteProduct).Methods("DELETE")

	s.HandleFunc("/cart", handlers.GetCart).Methods("GET")
	s.HandleFunc("/cart", handlers.AddToCart).Methods("POST")

	r.HandleFunc("/api/videos", handlers.GetVideos).Methods("GET")

	r.HandleFunc("/api/comments/{type}/{id}", handlers.GetComments).Methods("GET")
	s.HandleFunc("/api/comments/{type}/{id}", handlers.AddComment).Methods("POST")

	go handlers.HandleMessages()
	r.HandleFunc("/ws", handlers.HandleConnections)

	log.Printf("Server started on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
