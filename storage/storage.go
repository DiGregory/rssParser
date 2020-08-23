package storage

import (
	"database/sql"
	"os"
	"io/ioutil"
	"fmt"
	_ "github.com/lib/pq"
)

type NewsStorage struct {
	DB *sql.DB
}

type NewsStorager interface {
	CreateNews([]News)
}
type News struct {
	ID          int
	Title       string
	Description string
	Link        string
}

func (s *NewsStorage) CreateNews(news []*News) (error) {
	tx, err := s.DB.Begin()
	if err != nil {
		tx.Rollback()
		return   err
	}

	for _, n := range news {
		rows, err := tx.Query("INSERT INTO news(title,description,link) VALUES ( $1,$2,$3) ON CONFLICT DO NOTHING;",
			n.Title,n.Description,n.Link)
		rows.Close()
		if err != nil {
			tx.Rollback()
			return  err
		}
	}
	tx.Commit()

	return nil
}
func NewConn(driver, host, port, user, password, name string) (*NewsStorage, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, name)
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	initSQLFile, err := os.Open("./storage/init.sql")
	if err != nil {
		return nil, err
	}
	defer initSQLFile.Close()

	initQuery, err := ioutil.ReadAll(initSQLFile)
	if err != nil {
		return nil, err
	}

	_, err = db.Query(string(initQuery))
	if err != nil {
		return nil, err
	}

	fmt.Println("Connection with db was set up")
	return &NewsStorage{DB: db}, nil
}
