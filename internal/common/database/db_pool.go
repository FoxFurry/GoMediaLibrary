package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

// ifDBExists returns true if dbname exists in database instance.
// This function works specifically and only with postgres.
func ifDBExists(db *sql.DB, dbname string) bool {
	queryCheckExists := `SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname=$1);`

	row := db.QueryRow(queryCheckExists, dbname)

	var status bool
	err := row.Scan(&status)

	if err != nil {
		log.Panicf("Could not check if db:%v exists: %v", dbname, err)
	}

	return status
}

// createDatabase drops dbname from database instance and recreate it with template1 parameters (en_US.utf-8)
func createDatabase(db *sql.DB, dbname string) {
	queryCreate := `CREATE DATABASE ` + dbname + // Since postgres does not support parameters for
		` ENCODING    'utf8'` + // schema modifiers we have no choice but to concatenate
		` LC_COLLATE  'en_US.utf8'` + // db name into the query
		` LC_CTYPE    'en_US.utf8';`

	_, err := db.Query(queryCreate)
	if err != nil {
		log.Panicf("Could not create database: %v", err)
	}
}

// Creates a bookstore table(id, title, author, year, description) for a given database instance
func createBookstore(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS bookstore (		
					id SERIAL PRIMARY KEY,
					title TEXT NOT NULL,
					author TEXT NOT NULL,
					year INT NOT NULL,
					description TEXT
					);` // Should move it to separate sql file?

	_, err := db.Query(query)
	if err != nil {
		log.Panicf("Could not create tables: %v", err)
	}
}

// Create returns database connection pool with specified parameters. Function will validate database and tables
// before returning the instance
func Create(host string, port int, user string, pass string, dbname string, dbMaxIdleConns int, dbMaxOpenConns int, dbMaxIdleTime time.Duration) *sql.DB {
	log.Printf("DB configs:\nHost: %v\nPort: %v\nUser: %v\ndbName: %v\nMax idle conns: %v\nMax open conns: %v\nMax idle time: %v",
		host, port, user, dbname, dbMaxIdleConns, dbMaxOpenConns, dbMaxIdleTime)

	initTemplate := fmt.Sprintf("host=%s port=%d user=%s password=%s ", host, port, user, pass) // Connection template

	initTable := fmt.Sprintf("%s dbname=%s sslmode=disable", initTemplate, dbname) // Connection to specific db
	initDB := fmt.Sprintf("%s sslmode=disable", initTemplate)                      // Generic connection for db create

	db, err := sql.Open("postgres", initDB) // Connection for db creation
	if err != nil {
		log.Panicf("Could not connect to DB: %v", err)
	}

	if !ifDBExists(db, dbname) { // Since postgres does not allow IF NOT EXISTS for db alter - check and create are standalone
		log.Printf("Database %v does not exists, creating it", dbname)
		createDatabase(db, dbname)
	} else {
		log.Printf("Database %v already exists, skipping creation", dbname)
	}

	db, err = sql.Open("postgres", initTable) // Connection for table creation
	if err != nil {
		log.Panicf("Could not connect to db:%v: %v", dbname, err)
	}

	createBookstore(db)

	db.SetMaxIdleConns(dbMaxIdleConns)
	db.SetMaxOpenConns(dbMaxOpenConns)
	db.SetConnMaxIdleTime(dbMaxIdleTime)

	log.Println("Connected to database successfully")

	return db
}
