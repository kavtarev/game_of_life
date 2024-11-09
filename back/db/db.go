package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

type Storage struct {
	Con *sql.DB
}

func NewDb() (Storage, error) {
	useFakeStorage := os.Getenv("USE_FAKE_STORAGE")
	if useFakeStorage != "" {
		return Storage{}, nil
	}
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dbStr := fmt.Sprintf("postgres://%v:%v@127.0.0.1:%v/%v?sslmode=disable", user, pass, port, dbName)
	c, err := sql.Open("postgres", dbStr)
	if err != nil {
		return Storage{}, err
	}

	err = c.Ping()
	if err != nil {
		return Storage{}, err
	}

	return Storage{Con: c}, err
}
