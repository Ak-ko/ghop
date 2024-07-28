package db

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

func InitDB (config mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", config.FormatDSN()); if err != nil {
		log.Fatal(err)
	}

	checkConnection(db)

	return db, nil
}

func checkConnection(db *sql.DB) {
	dbErr := db.Ping(); if dbErr != nil {
		log.Fatal(dbErr)
	}

	log.Println("Database Connected Successfully")
}