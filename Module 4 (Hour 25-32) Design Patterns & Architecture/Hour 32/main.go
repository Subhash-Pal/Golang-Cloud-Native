package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserRepository interface {
	Create(user User) error
	FindAll() []User
	FindByID(id int) (User, bool)
}

type InMemoryUserRepository struct {
	mu    sync.Mutex
	users map[int]User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[int]User),
	}
}

func (r *InMemoryUserRepository) Create(user User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) FindAll() []User {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]User, 0, len(r.users))
	for _, user := range r.users {
		out = append(out, user)
	}
	return out
}

func (r *InMemoryUserRepository) FindByID(id int) (User, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	user, ok := r.users[id]
	return user, ok
}
/////
type UserUseCase struct {
	repo UserRepository
}

func NewUserUseCase(repo UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (uc *UserUseCase) CreateUser(user User) error {
	if user.ID == 0 || strings.TrimSpace(user.Name) == "" {
		return fmt.Errorf("id and name are required")
	}
	return uc.repo.Create(user)
}

func (uc *UserUseCase) ListUsers() []User {
	return uc.repo.FindAll()
}

func (uc *UserUseCase) GetUser(id int) (User, bool) {
	return uc.repo.FindByID(id)
}
////////////////////////////////////////////////////////
type UserHandler struct {
	uc *UserUseCase
}

func NewUserHandler(uc *UserUseCase) *UserHandler {
	return &UserHandler{uc: uc}
}

func (h *UserHandler) Users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createUser(w, r)
	case http.MethodGet:
		h.listUsers(w)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) UserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idText := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idText)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	user, ok := h.uc.GetUser(id)
	if !ok {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	if err := h.uc.CreateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) listUsers(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(h.uc.ListUsers())
}

func main() {
	repo := NewInMemoryUserRepository()
	uc := NewUserUseCase(repo)
	handler := NewUserHandler(uc)

	http.HandleFunc("/users", handler.Users)
	http.HandleFunc("/users/", handler.UserByID)

	log.Println("clean architecture demo running on http://127.0.0.1:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
