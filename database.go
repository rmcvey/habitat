package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func initDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "trust_data.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS habitats (
		id INTEGER PRIMARY KEY,
		name TEXT UNIQUE,
		ecosystem_m INTEGER,
		ecosystem_k INTEGER,
		ecosystem_b TEXT,
		trust_score REAL
	)`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS visit_times (
		id INTEGER PRIMARY KEY,
		habitat_id INTEGER,
		visit_time TIMESTAMP,
		FOREIGN KEY (habitat_id) REFERENCES habitats (id)
	)`)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
