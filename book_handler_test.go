package main

import (
	"github.com/google/uuid"
)

type mockBookRepo struct {
	books map[uuid.UUID]*Book
}

func NewMockBookRepo() *mockBookRepo {
	return &mockBookRepo{books: make(map[uuid.UUID]*Book)}
}
