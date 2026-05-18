package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BookRepo interface {
	GetById(ctx context.Context, id uuid.UUID) (*Book, error)
	Create(ctx context.Context, b *Book) error
	GetAllBooks(ctx context.Context) ([]*Book, error)
	Update(ctx context.Context, book *Book) error
	Delete(ctx context.Context, id uuid.UUID) error
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

func (r *bookRepo) Update(ctx context.Context, book *Book) error {
	query := `
		UPDATE books 
		SET title = $1, author = $2, year = $3, price = $4 
		WHERE id = $5
	`

	cmdTag, err := r.db.Exec(ctx, query, book.Title, book.Author, book.Year, book.Price, book.ID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("book not found")
	}

	return nil
}

func (r *bookRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM books WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("book not found")
	}

	return nil
}

