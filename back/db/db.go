package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewDb() {
	dbStr := "postgres://postgres:postgres@127.0.0.1:5433/postgres?sslmode=disable"
	c, err := sql.Open("postgres", dbStr)
	if err != nil {
		panic(err)
	}

	err = c.Ping()
	if err != nil {
		panic(err)
	}

	r, err := c.Query("select * from hui")
	if err != nil {
		panic(err)
	}

	var a any

	for r.Next() {
		err := r.Scan(&a)
		if err != nil {
			panic(err)
		}
		fmt.Println(555, a)
	}

}
