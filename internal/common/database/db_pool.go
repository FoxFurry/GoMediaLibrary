package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"time"
)


var db *sql.DB

func ConfigDBPool(dbURL string, database string,initTable string,dbMaxIdleConns int, dbMaxOpenConns int, dbMaxIdleTime time.Duration){
	log.Printf("DB configs:\nURL: %v\nMax idle conns: %v\nMax open conns: %v\nMax idle time: %v", dbURL, dbMaxIdleConns, dbMaxOpenConns, dbMaxIdleTime)


	var err error
	db, err = sql.Open("postgres", dbURL)

	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err!=nil {
		panic(err)
	}

	dropExisting := `DROP DATABASE IF EXISTS medialibrary;`
	db.Query(dropExisting)

	createNewDB := `CREATE DATABASE medialibrary
					ENCODING    'utf8'
					LC_COLLATE  'en_US.utf8'
					LC_CTYPE    'en_US.utf8';`

	db.Query(createNewDB)

	db, _ = sql.Open("postgres", database)

	if err != nil {
		panic(err)
	}

	c, ioErr := ioutil.ReadFile(initTable)
	if ioErr != nil {
		log.Fatalf("%v", err)
	}
	db.Query(string(c))



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
