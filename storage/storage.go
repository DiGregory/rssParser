package storage

import (
	"database/sql"
	"io/ioutil"
	"fmt"
	_ "github.com/lib/pq"
)

func NewConn(driver, host, port, user, password, name string) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`,
		host, port, user, password, name,
	)
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	initQuery, err := ioutil.ReadFile("./storage/init.sql")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(string(initQuery))
	if err != nil {
		return nil, err
	}

	return db, nil
}
