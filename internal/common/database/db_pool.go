package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

var db *sql.DB

func createDatabase(db *sql.DB) {
	query := `CREATE DATABASE IF NOT EXISTS medialibrary
					ENCODING    'utf8'
					LC_COLLATE  'en_US.utf8'
					LC_CTYPE    'en_US.utf8';`
	_, err := db.Query(query)
	if err != nil {
		log.Fatalf("Could not create database: %v", err)
	}
}

func createTables(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS bookstore (
					id SERIAL PRIMARY KEY,
					title TEXT NOT NULL,
					author TEXT NOT NULL,
					year INT NOT NULL,
					description TEXT
					);`

	_, err := db.Query(query)
	if err != nil {
		log.Fatalf("Could not create tables: %v", err)
	}
}

func ConfigDBPool(host string, port int, user string, pass string, dbname string, dbMaxIdleConns int, dbMaxOpenConns int, dbMaxIdleTime time.Duration) {
	log.Printf("DB configs:\nHost: %v\nPort: %v\nUser: %v\ndbName: %v\nMax idle conns: %v\nMax open conns: %v\nMax idle time: %v",
		host, port, user, dbname, dbMaxIdleConns, dbMaxOpenConns, dbMaxIdleTime)
	initTemplate := fmt.Sprintf("host=%s post=%d user=%s password=%s ", host, port, user, pass)

	initTable := fmt.Sprintf("%s dbname=%s sslmode=disable", initTemplate, dbname)
	initDB := fmt.Sprintf("%s sslmode=disable", initTemplate)

	var err error

	db, err = sql.Open("postgres", initDB)
	if err != nil {
		log.Fatalf("%v", err)
	}

	createDatabase(db)

	db, err = sql.Open("postgres", initTable)
	if err != nil {
		panic(err)
	}

	createTables(db)

	db.SetMaxIdleConns(dbMaxIdleConns)
	db.SetMaxOpenConns(dbMaxOpenConns)
	db.SetConnMaxIdleTime(dbMaxIdleTime)

	log.Println("Connected to database successfully")
}

func GetDB() *sql.DB {
	if err := db.Ping(); err != nil {
		log.Fatalf("Cannot ping DB: %v", err)
	}
	return db
}
