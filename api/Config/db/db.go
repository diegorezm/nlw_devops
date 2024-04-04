package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Databse struct {
	Connection *sql.DB
}


func NewDatabse(dbUrl string) *Databse {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Not able to open the databse!")
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute * 3)
	return &Databse{Connection: db}
}
