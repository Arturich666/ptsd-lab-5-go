package main

import (
	"encoding/json"
	"net/http"
)

type BookHandler struct {
	repo BookRepo
}

func NewBookHandler(repo BookRepo) *BookHandler {
	return &BookHandler{repo: repo}
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
