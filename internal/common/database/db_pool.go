package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"time"
)


var db *sql.DB

func ConfigDBPool(dbURL string, dbMaxIdleConns int, dbMaxOpenConns int, dbMaxIdleTime time.Duration){
	log.Printf("DB configs:\nURL: %v\nMax idle conns: %v\nMax open conns: %v\nMax idle time: %v", dbURL, dbMaxIdleConns, dbMaxOpenConns, dbMaxIdleTime)


	var err error
	db, err = sql.Open("postgres", dbURL)

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
}

func GetDB() *sql.DB {
	if err:=db.Ping(); err != nil {
		log.Fatalf("Cannot ping DB: %v", err)
	}
	return db
}
