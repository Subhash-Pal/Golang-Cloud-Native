package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// Book represents the domain model.
type Book struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
}

// BookRepository acts like the repository contract from the domain/use case layer.
type BookRepository interface {
	Create(book Book) error
	FindAll() []Book
}

type InMemoryBookRepository struct {
	mu    sync.Mutex
	books []Book
}

func (r *InMemoryBookRepository) Create(book Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.books = append(r.books, book)
	return nil
}

func (r *InMemoryBookRepository) FindAll() []Book {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]Book, len(r.books))
	copy(out, r.books)
	return out
}

type BookUseCase struct {
	repo BookRepository
}

func NewBookUseCase(repo BookRepository) *BookUseCase {
	return &BookUseCase{repo: repo}
}

func (uc *BookUseCase) CreateBook(book Book) error {
	if book.Title == "" || book.Author == "" {
		return fmt.Errorf("title and author are required")
	}
	return uc.repo.Create(book)
}

func (uc *BookUseCase) ListBooks() []Book {
	return uc.repo.FindAll()
}

type HTTPHandler struct {
	uc *BookUseCase
}

func NewHTTPHandler(uc *BookUseCase) *HTTPHandler {
	return &HTTPHandler{uc: uc}
}

func (h *HTTPHandler) Books(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listBooks(w)
	case http.MethodPost:
		h.createBook(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *HTTPHandler) createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	if err := h.uc.CreateBook(book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(book)
}

func (h *HTTPHandler) listBooks(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(h.uc.ListBooks())
}

func main() {
	repo := &InMemoryBookRepository{}
	useCase := NewBookUseCase(repo)
	handler := NewHTTPHandler(useCase)

	http.HandleFunc("/books", handler.Books)

	log.Println("server running on http://127.0.0.1:8080")
	log.Println("try POST /books then GET /books")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
