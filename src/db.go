package main

import (
	"database/sql"
	"log"
	"os"
	"time"
)

func initDB() *sql.DB {
	conn := connectToDB()
	if conn == nil {
		log.Panic("cannot connect to database")
	}
	return conn
}

func connectToDB() *sql.DB {
	counts := 0
	dsn := os.Getenv("DSN")
	// DSN="host=localhost port=5432 user=postgres password=password dbname=concurrency sslmode=disable timezone=UTC connect_timeout=5"
	for {
		conn, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres is not yet ready. err: ", err)
		} else {
			log.Printf("Connected successful")
			return conn
		}

		if counts > 10 {
			return nil
		}
		log.Printf("sleep 1 sec")
		time.Sleep(1 * time.Second)
		counts++
	}
}

func openDB(dsn string) (*sql.DB, error) {
	// DSN="host=localhost port=5432 user=postgres password=password dbname=concurrency sslmode=disable timezone=UTC connect_timeout=5"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
