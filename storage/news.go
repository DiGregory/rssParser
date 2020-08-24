package storage

import (
	"database/sql"
)

type NewsStorager interface {
	CreateNews(news []*News) error
	GetNews(limit, offset *int32) ([]*News, error)
}

type NewsStorage struct {
	DB *sql.DB
}

func NewNewsStorage(conn *sql.DB) *NewsStorage {
	return &NewsStorage{DB: conn}
}

type News struct {
	ID          int32
	Title       string
	Description string
	Link        string
}

func (s *NewsStorage) CreateNews(news []*News) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	for _, n := range news {
		_, err := tx.Exec(
			`INSERT INTO news(title, description, link) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING;`,
			n.Title, n.Description, n.Link,
		)
		if err != nil {
			return tx.Rollback()
		}
	}
	return tx.Commit()
}

func (s *NewsStorage) GetNews(limit, offset *int32) ([]*News, error) {
	rows, err := s.DB.Query(
		`SELECT * FROM news LIMIT coalesce($1, 10) OFFSET coalesce($2, 0);`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	news := make([]*News, 0)
	for rows.Next() {
		var n News
		if err = rows.Scan(&n.ID, &n.Title, &n.Description, &n.Link); err != nil {
			return nil, err
		}
		news = append(news, &n)
	}
	return news, nil
}
