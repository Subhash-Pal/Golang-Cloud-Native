# Prompt - Hour 26 Folder Structure Implementation

Write Go code that implements a Clean Architecture folder structure for a `book store` service.

## Requirements

- Organize the project into folders such as `cmd`, `internal/domain`, `internal/usecase`, `internal/repository`, and `internal/delivery/http`
- Add a `Book` entity with fields like `ID`, `Title`, `Author`, and `Price`
- Implement one use case to create a book and another to fetch all books
- Add an in-memory repository implementation
- Add simple HTTP handlers using the standard library `net/http`
- Keep each package focused on a single responsibility
- Add a `main.go` entry point inside `cmd`

## Output expectations

- The code should compile as one runnable project
- Folder names and package responsibilities should be easy for beginners to follow
- Include minimal routing for `POST /books` and `GET /books`

## Goal

The final code should teach how project structure supports architecture, readability, and long-term maintainability.
