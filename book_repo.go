package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BookRepo interface {
	GetAllBooks(ctx context.Context) ([]*Book, error)
}

type bookRepo struct {
	db *pgxpool.Pool
}

func NewBookRepo(db *pgxpool.Pool) BookRepo {
	return &bookRepo{db: db}
}
func (r *bookRepo) GetAllBooks(ctx context.Context) ([]*Book, error) {
	query := "SELECT id, title, author, year, price FROM books"

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*Book

	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.Price); err != nil {
			return nil, err
		}
		books = append(books, &b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if books == nil {
		books = make([]*Book, 0)
	}

	return books, nil
}
