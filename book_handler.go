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
	"encoding/json"
	"net/http"
)

type BookHandler struct {
	repo BookRepo
	func NewBookHandler() *BookHandler {
		return &BookHandler{}
	}
}

func NewBookHandler(repo BookRepo) *BookHandler {
	return &BookHandler{repo: repo}
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
func (h *BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.repo.GetAllBooks(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (h *BookHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid UUID format", http.StatusBadRequest)
		return
	}

	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	book.ID = id
	if err := book.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.repo.Update(r.Context(), &book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid UUID format", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
