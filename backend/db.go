package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", "data.db")

	if err != nil {
		log.Fatal(err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS guests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		room_number TEXT NOT NULL,
		daily_rate INTEGER NOT NULL,
		check_in_date DATE NOT NULL,
		contact TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS notifications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		guest_id INTEGER NOT NULL,
		
		period_number INTEGER NOT NULL,
		sent_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (guest_id, period_number)
	);`

	_, err = db.Exec(schema)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
