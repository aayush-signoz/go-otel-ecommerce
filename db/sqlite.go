package db

import (
	"database/sql"
	"log"
)

var DB *sql.DB

func SetupDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./ecommerce.db")
	if err != nil {
		log.Fatal(err)
	}

	DB.Exec(`CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	)`)

	DB.Exec(`CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		product_id INTEGER NOT NULL,
		quantity INTEGER NOT NULL,
		user_id TEXT
	)`)

	DB.Exec(`INSERT OR IGNORE INTO products (name) VALUES 
		('Book'), ('Laptop'), ('Phone')
	`)
}
