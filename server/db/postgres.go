package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type database struct {
	db *sql.DB
}

func newDatabase() (*database, error) {

	db, err := sql.Open("postgres", "postgres://root:root@localhost:5432/chat?sslmode=disable")
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &database{db: db}, nil
}

func (DB *database) Close() {
	DB.db.Close()
}

func (DB *database) GetDB() *sql.DB {
	return DB.db
}

var DB, _ = newDatabase()
