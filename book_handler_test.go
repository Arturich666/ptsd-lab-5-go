package main

import (
	"github.com/google/uuid"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

type mockBookRepo struct {
	books map[uuid.UUID]*Book
}

func NewMockBookRepo() *mockBookRepo {
	return &mockBookRepo{books: make(map[uuid.UUID]*Book)}
}

func (m *mockBookRepo) GetById(ctx context.Context, id uuid.UUID) (*Book, error) {
	book, ok := m.books[id]
	if !ok {
		return nil, errors.New("book not found")
	}
	return book, nil
}

func (m *mockBookRepo) Create(ctx context.Context, book *Book) error {
	m.books[book.ID] = book
	return nil
}

func (m *mockBookRepo) Patch(ctx context.Context, id uuid.UUID, title, author *string, year *int, price *float64) error {
	book, ok := m.books[id]
	if !ok {
		return errors.New("book not found")
	}
	if title != nil {
		book.Title = *title
	}
	if author != nil {
		book.Author = *author
	}
	if year != nil {
		book.Year = *year
	}
	if price != nil {
		book.Price = *price
	}
	return nil
}

func (m *mockBookRepo) GetAllBooks(ctx context.Context) ([]*Book, error) {
	var list []*Book
	for _, b := range m.books {
		list = append(list, b)
	}
	return list, nil
}

func (m *mockBookRepo) Update(ctx context.Context, book *Book) error {
	if _, ok := m.books[book.ID]; !ok {
		return errors.New("book not found")
	}
	m.books[book.ID] = book
	return nil
}

func (m *mockBookRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if _, ok := m.books[id]; !ok {
		return errors.New("book not found")
	}
	delete(m.books, id)
	return nil
}
// Tests
func TestCreateBook_Success(t *testing.T) {
	mockRepo := NewMockBookRepo()
	handler := NewBookHandler(mockRepo)


	newBook := map[string]interface{}{
		"title":  "The Hobbit",
		"author": "J.R.R. Tolkien",
		"year":   1937,
		"price":  45.0,
	}
	body, _ := json.Marshal(newBook)

	req := httptest.NewRequest("POST", "/books", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.Create(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	if len(mockRepo.books) != 1 {
		t.Errorf("expected 1 book in mock db, got %d", len(mockRepo.books))
	}
}

func TestGetBookById_Success(t *testing.T) {
	mockRepo := NewMockBookRepo()
	id := uuid.New()
	mockRepo.books[id] = &Book{ID: id, Title: "Animal Farm", Author: "George Orwell", Year: 1945, Price: 12.0}

	handler := NewBookHandler(mockRepo)
	req := httptest.NewRequest("GET", "/books/"+id.String(), nil)

	r := chi.NewRouter()
	r.Get("/books/{id}", handler.GetById)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestPatchBook_Success(t *testing.T) {
	mockRepo := NewMockBookRepo()
	id := uuid.New()
	mockRepo.books[id] = &Book{ID: id, Title: "Kobzar", Author: "Taras Shevchenko", Year: 1840, Price: 100.0}

	handler := NewBookHandler(mockRepo)

	patchData := map[string]interface{}{
		"price": 250.0,
	}
	body, _ := json.Marshal(patchData)

	req := httptest.NewRequest("PATCH", "/books/"+id.String(), bytes.NewBuffer(body))
	r := chi.NewRouter()
	r.Patch("/books/{id}", handler.Patch)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if mockRepo.books[id].Price != 250.0 {
		t.Errorf("expected price to be 250.0, got %v", mockRepo.books[id].Price)
	}
	if mockRepo.books[id].Title != "Kobzar" {
		t.Errorf("expected title to remain 'Kobzar', got %v", mockRepo.books[id].Title)
	}
}

func TestPatchBook_ValidationError(t *testing.T) {
	mockRepo := NewMockBookRepo()
	id := uuid.New()
	mockRepo.books[id] = &Book{ID: id, Title: "Kobzar", Author: "Taras Shevchenko", Year: 1840, Price: 100.0}

	handler := NewBookHandler(mockRepo)

	patchData := map[string]interface{}{
		"title": "",
	}
	body, _ := json.Marshal(patchData)

	req := httptest.NewRequest("PATCH", "/books/"+id.String(), bytes.NewBuffer(body))
	r := chi.NewRouter()
	r.Patch("/books/{id}", handler.Patch)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request due to empty title validation, got %v", status)
	}

	if !strings.Contains(rr.Body.String(), "Title cannot be empty") {
		t.Errorf("expected error message 'Title cannot be empty', got: %s", rr.Body.String())
	}
}

func TestGetAllBooks(t *testing.T) {
	mockRepo := NewMockBookRepo()
	id := uuid.New()
	mockRepo.books[id] = &Book{ID: id, Title: "1984", Author: "George Orwell", Year: 1949, Price: 15.0}

	handler := NewBookHandler(mockRepo)
	req := httptest.NewRequest("GET", "/books", nil)
	rr := httptest.NewRecorder()

	handler.GetAllBooks(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestUpdateBook_Validation(t *testing.T) {
	mockRepo := NewMockBookRepo()
	id := uuid.New()
	mockRepo.books[id] = &Book{ID: id, Title: "1984", Author: "George Orwell", Year: 1949, Price: 15.0}

	handler := NewBookHandler(mockRepo)

	invalidBook := Book{Title: "1984", Author: "George Orwell", Year: 1949, Price: -5.0}
	body, _ := json.Marshal(invalidBook)

	req := httptest.NewRequest("PUT", "/books/"+id.String(), bytes.NewBuffer(body))
	r := chi.NewRouter()
	r.Put("/books/{id}", handler.Update)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("expected validation to fail with 400 Bad Request, got %v", status)
	}
}

func TestDeleteBook(t *testing.T) {
	mockRepo := NewMockBookRepo()
	id := uuid.New()
	mockRepo.books[id] = &Book{ID: id, Title: "To Kill a Mockingbird", Author: "Harper Lee", Year: 1960, Price: 10.0}

	handler := NewBookHandler(mockRepo)
	req := httptest.NewRequest("DELETE", "/books/"+id.String(), nil)

	r := chi.NewRouter()
	r.Delete("/books/{id}", handler.Delete)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	if len(mockRepo.books) != 0 {
		t.Errorf("expected book to be deleted from mock db")
	}
}