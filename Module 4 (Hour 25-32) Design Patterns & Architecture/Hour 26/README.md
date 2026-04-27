# Hour 26 - Folder Structure Implementation

This hour turns architecture theory into a practical Go project layout. A clear folder structure makes it easier to scale code, onboard new developers, and test layers independently.

## Key concepts

- `cmd` contains entry points
- `internal/domain` contains entities and interfaces central to the business
- `internal/usecase` contains application logic
- `internal/repository` contains persistence implementations
- `internal/delivery` contains transport adapters such as HTTP handlers

## Learning goals

- translate architecture ideas into directories and packages
- avoid mixing HTTP, business, and data access code
- understand where interfaces and implementations should live
- build a maintainable project structure for real applications

## Suggested deliverable

Build a small book store API with standard library HTTP handlers and an in-memory repository using a Clean Architecture style layout.

## Files

- `main.go`: book store API example with repository, use case, and HTTP handler
- `PROMPT.md`: coding prompt for practice
- `go.mod`: module definition for the hour

## Run

```powershell
cd "D:\training_golang\Module 4 (Hour 25-32) Design Patterns & Architecture\Hour 26"
go run .
```

Then test:

```powershell
Invoke-RestMethod -Method Post -Uri "http://127.0.0.1:8080/books" -ContentType "application/json" -Body '{"id":1,"title":"Go Patterns","author":"Team","price":499}'
Invoke-RestMethod -Method Get -Uri "http://127.0.0.1:8080/books"
```
