package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func createDatabase(db *sql.DB) {
	queryDelete := `DROP DATABASE medialibrary;`
	queryCreate := `CREATE DATABASE medialibrary
					ENCODING    'utf8'
					LC_COLLATE  'en_US.utf8'
					LC_CTYPE    'en_US.utf8';`

	_, err := db.Query(queryDelete)
	if err != nil {
		log.Panicf("Could not delete database: %v", err)
	}

	_, err = db.Query(queryCreate)
	if err != nil {
		log.Panicf("Could not create database: %v", err)
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
		log.Panicf("Could not create tables: %v", err)
	}
}

func CreateDBPool(host string, port int, user string, pass string, dbname string, dbMaxIdleConns int, dbMaxOpenConns int, dbMaxIdleTime time.Duration) *sql.DB {
	log.Printf("DB configs:\nHost: %v\nPort: %v\nUser: %v\ndbName: %v\nMax idle conns: %v\nMax open conns: %v\nMax idle time: %v",
		host, port, user, dbname, dbMaxIdleConns, dbMaxOpenConns, dbMaxIdleTime)
	initTemplate := fmt.Sprintf("host=%s port=%d user=%s password=%s ", host, port, user, pass)

	initTable := fmt.Sprintf("%s dbname=%s sslmode=disable", initTemplate, dbname)
	initDB := fmt.Sprintf("%s sslmode=disable", initTemplate)

	var err error

	db, err := sql.Open("postgres", initDB)
	if err != nil {
		log.Panicf("%v", err)
	}

	createDatabase(db)

	db, err = sql.Open("postgres", initTable)
	if err != nil {
		log.Panicf("%v", err)
	}

	createTables(db)

	db.SetMaxIdleConns(dbMaxIdleConns)
	db.SetMaxOpenConns(dbMaxOpenConns)
	db.SetConnMaxIdleTime(dbMaxIdleTime)

	log.Println("Connected to database successfully")

	return db
}
