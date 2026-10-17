package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"robotics-project/database"
	"robotics-project/middleware"
	"robotics-project/models"
	"github.com/gorilla/mux"
)

// crud
func GetProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name, description, price, image_url, stock FROM products")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.Stock)
		products = append(products, p)
	}
	json.NewEncoder(w).Encode(products)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := database.DB.QueryRow(
		"INSERT INTO products (name, description, price, image_url, stock) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		p.Name, p.Description, p.Price, p.ImageURL, p.Stock,
	).Scan(&p.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec("UPDATE products SET name=$1, description=$2, price=$3, image_url=$4, stock=$5 WHERE id=$6",
		p.Name, p.Description, p.Price, p.ImageURL, p.Stock, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Product updated successfully"})
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	_, err := database.DB.Exec("DELETE FROM products WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted successfully"})
}

// cart
func GetCart(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middleware.UserIDKey).(int)
	rows, err := database.DB.Query(`
		SELECT c.id, c.quantity, p.id, p.name, p.price, p.image_url 
		FROM cart_items c JOIN products p ON c.product_id = p.id WHERE c.user_id = $1`, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []models.CartItem
	for rows.Next() {
		var item models.CartItem
		item.Product = &models.Product{}
		rows.Scan(&item.ID, &item.Quantity, &item.Product.ID, &item.Product.Name, &item.Product.Price, &item.Product.ImageURL)
		items = append(items, item)
	}
	json.NewEncoder(w).Encode(items)
}

func AddToCart(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middleware.UserIDKey).(int)
	var body struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec("INSERT INTO cart_items (user_id, product_id, quantity) VALUES ($1, $2, $3)",
		userId, body.ProductID, body.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Added to cart"})
}
