package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/google/uuid"
	"errors"
)

type BookRepo interface {
	GetById(ctx context.Context, id uuid.UUID) (*Book, error)
	Create(ctx context.Context, b *Book) error
	GetAllBooks(ctx context.Context) ([]*Book, error)
	Update(ctx context.Context, book *Book) error
	Delete(ctx context.Context, id uuid.UUID) error
	Patch(ctx context.Context, id uuid.UUID, title *string, author *string, year *int, price *float64) error
}

type bookRepo struct {
	db *pgxpool.Pool
}

func NewBookRepo(db *pgxpool.Pool) BookRepo {
	return &bookRepo{db: db}
}


func (r *bookRepo) GetById(ctx context.Context, id uuid.UUID) (*Book, error) {
	query := "SELECT id, title, author, year, price FROM books WHERE id = $1"
	b := &Book{} 
	
	err := r.db.QueryRow(ctx, query, id).Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.Price)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r *bookRepo) Create(ctx context.Context, b *Book) error {
	query := "INSERT INTO books (id, title, author, year, price) VALUES ($1, $2, $3, $4, $5)"
	
	_, err := r.db.Exec(ctx, query, b.ID, b.Title, b.Author, b.Year, b.Price)
	return err
}

func (r *bookRepo) Patch(ctx context.Context, id uuid.UUID, title *string, author *string, year *int, price *float64) error {
	query := "UPDATE books SET "
	args := []interface{}{}
	argIdx := 1

	if title != nil {
		query += fmt.Sprintf("title = $%d, ", argIdx)
		args = append(args, *title)
		argIdx++
	}
	if author != nil {
		query += fmt.Sprintf("author = $%d, ", argIdx)
		args = append(args, *author)
		argIdx++
	}
	if year != nil {
		query += fmt.Sprintf("year = $%d, ", argIdx)
		args = append(args, *year)
		argIdx++
	}
	if price != nil {
		query += fmt.Sprintf("price = $%d, ", argIdx)
		args = append(args, *price)
		argIdx++
	}

	if len(args) == 0 {
		return nil
	}

	query = query[:len(query)-2] + fmt.Sprintf(" WHERE id = $%d", argIdx)
	args = append(args, id)

	cmdTag, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("book not found")
	}

	return nil
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

