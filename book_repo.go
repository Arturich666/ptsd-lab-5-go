package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookRepo interface {
	GetById(ctx context.Context, id uuid.UUID) (*Book, error)
}

type bookRepo struct {
	db *pgxpool.Pool
}
