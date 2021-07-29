package database

import (
	"database/sql"
	"log"
	"time"
)

func CreateDBpool(dbURL string, dbMaxIdleConns int, dbMaxOpenConns int, dbMaxIdleTime time.Duration) *sql.DB{
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err!=nil {
		panic(err)
	}

	db.SetMaxIdleConns(dbMaxIdleConns)
	db.SetMaxOpenConns(dbMaxOpenConns)
	db.SetConnMaxIdleTime(dbMaxIdleTime)

	log.Println("Connected to database succesfully")

	return db
}
