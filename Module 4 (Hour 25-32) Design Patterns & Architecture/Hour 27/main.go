package main

import (
	"fmt"
	"strings"
	"sync"
)

type User struct {
	ID    int
	Name  string
	Email string
}

type UserRepository interface {
	Create(user User) error
	FindByID(id int) (User, bool)
	FindAll() []User
	Delete(id int) error
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

func (r *InMemoryUserRepository) FindByID(id int) (User, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	user, ok := r.users[id]
	return user, ok
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

func (r *InMemoryUserRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.users, id)
	return nil
}

type UserService struct {
	repo UserRepository
}
////////////////////////////////
func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(user User) error {
	for _, existing := range s.repo.FindAll() {
		if strings.EqualFold(existing.Email, user.Email) {
			return fmt.Errorf("email already exists: %s", user.Email)
		}
	}
	return s.repo.Create(user)
}

func main() {
	repo := NewInMemoryUserRepository()
	service := NewUserService(repo)

	users := []User{
		{ID: 1, Name: "Asha", Email: "asha@example.com"},
		{ID: 2, Name: "Rohan", Email: "rohan@example.com"},
		{ID: 3, Name: "Duplicate", Email: "asha@example.com"},
	}

	for _, user := range users {
		if err := service.Register(user); err != nil {
			fmt.Println("register error:", err)
			continue
		}
		fmt.Println("registered:", user.Name)
	}

	fmt.Println("\nall users:")
	for _, user := range repo.FindAll() {
		fmt.Printf("- #%d %s <%s>\n", user.ID, user.Name, user.Email)
	}
}
