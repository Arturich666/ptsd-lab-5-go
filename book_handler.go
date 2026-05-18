package main

import (
	"context"
	"net/http"
	"os"
    "fmt"
	"github.com/go-chi/chi/v5"
	
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/google/uuid"

	"encoding/json"

	"strings"
)


type BookHandler struct {
	repo BookRepo
	func NewBookHandler() *BookHandler {
		return &BookHandler{}
	}
}

	func GetById(w http.ResponseWriter, r *http.Request) {
			idStr := chi.URLParam(r, "id")
			id, err := uuid.Parse(idStr)
			if err != nil {
				http.Error(w, "invalid id", http.StatusBadRequest)
			}
			
			query := "SELECT id, title, author, year, price FROM books where id=$1"
			b := &Book{} 
			err = dbpool.QueryRow(r.Context(), query, id).Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.Price)
			if err != nil {
				fmt.Errorf("get book by id", err)
				http.Error(w, "book not found", http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(b)

	}
