package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/ncruces/go-sqlite3/driver"
)

var DB *sql.DB

func Connect_to_db() {
	var err error
	DB, err = sql.Open("sqlite3", "app.DB")
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
		enc_key TEXT NOT NULL
	);
	`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Tables initialised successfully!")
}
