package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/ncruces/go-sqlite3/driver"
)

var DB *sql.DB

func Connect_to_db() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = ":memory:"
	}

	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Couldnt connect to DB")
	}
	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to SQLite!")
}

func Init() {
	fmt.Println("Creating table")
	query := `
	CREATE TABLE IF NOT EXISTS users(
		id TEXT PRIMARY KEY NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		mpin INT NOT NULL,
		salt TEXT NOT NULL DEFAULT ''
	);
	`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	// Migrate existing database to add salt column if missing
	_, _ = DB.Exec("ALTER TABLE users ADD COLUMN salt TEXT NOT NULL DEFAULT ''")
	query = `
	CREATE TABLE IF NOT EXISTS passwords(
		id TEXT PRIMARY KEY NOT NULL,
		user_id TEXT NOT NULL,
		title TEXT NOT NULL,
		email TEXT NOT NULL,
		password TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`

	_, err = DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Tables initialised successfully!")
}
