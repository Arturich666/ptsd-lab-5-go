package main

import (
	"context"
	"fmt"

	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbpool, err := pgxpool.New(context.Background(), "postgres://postgres:password@localhost:5432/postgres")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to db %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	bookRepo := NewBookRepo(dbpool)
	bookHandler := NewBookHandler(bookRepo)

	r := chi.NewRouter()
	
	r.Post("/books", bookHandler.Create)
	r.Get("/books", bookHandler.GetAllBooks)

	r.Get("/books/{id}", bookHandler.GetById)
	r.Put("/books/{id}", bookHandler.Update)
	r.Patch("/books/{id}", bookHandler.Patch)
	r.Delete("/books/{id}", bookHandler.Delete)

	

	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", r)
}
