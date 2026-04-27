package main

import (
	"context"
	"fmt"
	"time"
)

// Task is the domain entity. It represents core business data.
type Task struct {
	ID        int
	Title     string
	Completed bool
	CreatedAt time.Time
}

// TaskRepository is an abstraction owned by the business layer.
type TaskRepository interface {
	Save(ctx context.Context, task Task) error
	FindAll(ctx context.Context) ([]Task, error)
}

// TaskUseCase contains application rules.
type TaskUseCase struct {
	repo TaskRepository
}

func NewTaskUseCase(repo TaskRepository) *TaskUseCase {
	return &TaskUseCase{repo: repo}
}

func (uc *TaskUseCase) CreateTask(ctx context.Context, id int, title string) error {
	if title == "" {
		return fmt.Errorf("title is required")
	}

	task := Task{
		ID:        id,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	return uc.repo.Save(ctx, task)
}

func (uc *TaskUseCase) ListTasks(ctx context.Context) ([]Task, error) {
	return uc.repo.FindAll(ctx)
}

// InMemoryTaskRepository is an infrastructure detail.
type InMemoryTaskRepository struct {
	tasks []Task
}

func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{tasks: make([]Task, 0)}
}

func (r *InMemoryTaskRepository) Save(_ context.Context, task Task) error {
	r.tasks = append(r.tasks, task)
	return nil
}

func (r *InMemoryTaskRepository) FindAll(_ context.Context) ([]Task, error) {
	out := make([]Task, len(r.tasks))
	copy(out, r.tasks)
	return out, nil
}

func main() {
	ctx := context.Background()

	repo := NewInMemoryTaskRepository()
	useCase := NewTaskUseCase(repo)

	_ = useCase.CreateTask(ctx, 1, "Learn clean architecture")
	_ = useCase.CreateTask(ctx, 2, "Build a Go example")

	tasks, err := useCase.ListTasks(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Tasks from use case layer:")
	for _, task := range tasks {
		fmt.Printf("- #%d %s (completed=%t)\n", task.ID, task.Title, task.Completed)
	}
}
