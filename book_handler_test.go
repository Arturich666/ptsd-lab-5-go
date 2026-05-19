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

// Tests
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