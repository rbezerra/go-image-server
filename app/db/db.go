package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Panic(err)
	}

	db.SetMaxOpenConns(500)
	db.SetConnMaxLifetime(time.Hour)
	db.SetMaxIdleConns(100)

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
}
