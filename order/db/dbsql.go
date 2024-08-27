package db

// import (
// 	"context"
// 	"database/sql"

// 	"github.com/hsxflowers/restaurante-digital/order/domain"
// 	"github.com/labstack/gommon/log"
// )

// type SQLStore struct {
// 	db *sql.DB
// }

// func NewSQLStore(db *sql.DB) *SQLStore {
// 	return &SQLStore{
// 		db: db,
// 	}
// }

// func (s *SQLStore) GetItem(ctx context.Context, tag string) (*domain.Item, error) {
// 	var cat domain.Item
// 	var url, tagResult string

// 	query := "SELECT url, tag FROM cats WHERE tag = $1 ORDER BY RANDOM() LIMIT 1"
// 	row := s.db.QueryRowContext(ctx, query, tag)

// 	err := row.Scan(&url, &tagResult)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, exceptions.New(exceptions.ErrCatNotFound, err)
// 		}
// 		log.Error("Error fetching cat from database: ", err)
// 		return nil, exceptions.New(exceptions.ErrInternalServer, err)
// 	}

// 	cat.Url = url
// 	cat.Tag = tagResult

// 	return &cat, nil
// }

// func (s *SQLStore) Create(ctx context.Context, cat *domain.Cat) error {
// 	_, err := s.db.ExecContext(ctx, "INSERT INTO cats (url, tag) VALUES ($1, $2)", cat.Url, cat.Tag)
// 	if err != nil {
// 		log.Error("Error creating cat in database: ", err)
// 		return exceptions.New(exceptions.ErrInternalServer, err)
// 	}
// 	return nil
// }
