package database

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitPostgres(connStr string) {
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("PostgreSQL connection error: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("PostgreSQL ping error: %v", err)
	}

	createTables()
	seedData()
	log.Println("PostgreSQL connected successfully!")
}

func createTables() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT,
		price NUMERIC(10, 2) NOT NULL,
		image_url TEXT,
		stock INT DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS cart_items (
		id SERIAL PRIMARY KEY,
		user_id INT REFERENCES users(id) ON DELETE CASCADE,
		product_id INT REFERENCES products(id) ON DELETE CASCADE,
		quantity INT DEFAULT 1
	);

	CREATE TABLE IF NOT EXISTS videos (
		id SERIAL PRIMARY KEY,
		title VARCHAR(150) NOT NULL,
		youtube_id VARCHAR(50) NOT NULL,
		description TEXT
	);
	`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
}

func seedData() {
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if count == 0 {
		DB.Exec(`INSERT INTO products (name, description, price, image_url, stock) VALUES 
			('Arduino Uno R3', 'Популярная плата arduino uno r3', 25.50, 'https://images.unsplash.com/photo-1553406830-ef2513450d76?q=80&w=1631&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=NHwxMjA3fDA%3D', 15),
			('ESP32 Wi-Fi + Bluetooth', 'Мощный модуль для IoT проектов и умного дома', 12.30, 'https://images.unsplash.com/photo-1789036070970-62285231046a?q=80&w=1170&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=NHwxMjA3fDA%3D', 30),
			('Датчик расстояния HC-SR04', 'Ультразвуковой дальномер для роботов', 3.20, 'https://cdn.insales-shop.ru/images/products/1/4578/441561570/rangefinder-ultrasonic-hc-sr04.2.jpg', 50);`)
	}

	DB.QueryRow("SELECT COUNT(*) FROM videos").Scan(&count)
	if count == 0 {
		DB.Exec(`INSERT INTO videos (title, youtube_id, description) VALUES 
			('1 Урок по Python. Print(). Переменные. Типы данных.', '3zLD1APZUn4', 'Базовое введение.'),
			('2 Урок по Python.  Input(). Условные операторы if, else, elif. Внутренние операторы and, or.', 'c4zItcGifE8', 'Создаем первые условия.');`),
      ('3 Урок по Python. Циклы for, while.', '8p0KC4akylQ', 'Задачи с циклами.');`)
	}
}
