package main

type BookHandler struct {
	repo BookRepo
}

func NewBookHandler(repo BookRepo) *BookHandler {
	return &BookHandler{repo: repo}
}
