package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
	"robotics-project/database"
	"robotics-project/middleware"
	"robotics-project/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	err = database.DB.QueryRow(
		"INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id, created_at",
		u.Username, u.Email, string(hashed),
	).Scan(&u.ID, &u.CreatedAt)

	if err != nil {
		http.Error(w, "User already exists or database error", http.StatusBadRequest)
		return
	}

	u.PasswordHash = ""
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds models.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var u models.User
	err := database.DB.QueryRow("SELECT id, username, email, password_hash FROM users WHERE username = $1", creds.Username).
		Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash)

	if err != nil || bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(creds.PasswordHash)) != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"user_id":  u.ID,
		"username": u.Username,
		"exp":      expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middleware.JwtKey)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString, "username": u.Username})
}

func Profile(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middleware.UserIDKey).(int)

	var u models.User
	err := database.DB.QueryRow("SELECT id, username, email, created_at FROM users WHERE id = $1", userId).
		Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(u)
}
