package main

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID     uuid.UUID `json:"id"`
	Title  string    `json:"title"`
	Author string    `json:"author"`
	Year   int       `json:"year"`
	Price  float64   `json:"price"`
}

func (b *Book) Validate() error {
	if b.Title == "" {
		return errors.New("Заголовок не може бути порожнім")
	}
	if b.Author == "" {
		return errors.New("Автор не може бути порожнім")
	}
	if b.Year <= 0 || b.Year > time.Now().Year() {
		return errors.New("Некоректний рік")
	}
	if b.Price <= 0 {
		return errors.New("Ціна повинна бути додатним числом")
	}
	return nil
}