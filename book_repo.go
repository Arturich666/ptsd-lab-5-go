package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookRepo interface {
	GetById(ctx context.Context, id uuid.UUID) (*Book, error)
	Create(ctx context.Context, b *Book) error
}

type bookRepo struct {
	db *pgxpool.Pool
}

func NewBookRepo(db *pgxpool.Pool) BookRepo {
	return &bookRepo{db: db}
}

func (r *bookRepo) GetById(ctx context.Context, id uuid.UUID) (*Book, error) {
query := "SELECT id, title, author, year, price FROM books where id=$1"
			b := &Book{} 
			err := r.db.QueryRow(ctx, query, id).Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.Price)
			if err != nil {
				fmt.Errorf("get book by id", err)
				http.Error(w, "book not found", http.StatusNotFound)
				return nil, err
			}
			return b, nill
}


func (r *bookRepo) Create(ctx context.Context, id uuid.UUID) error {}
