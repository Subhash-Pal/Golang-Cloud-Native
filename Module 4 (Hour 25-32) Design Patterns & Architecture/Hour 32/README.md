# Hour 32 - Mock Test: Refactor Existing API to Clean Architecture

This mock test combines the main ideas from the full module. The learner starts with a tightly coupled API and refactors it into a cleaner, layered design.

## What this test measures

- understanding of clean architecture boundaries
- ability to separate handlers, use cases, and repositories
- correct use of interfaces and dependency injection
- confidence refactoring existing code without changing behavior

## Suggested challenge

Take a simple user API that stores data in memory and move it into:

- domain entities
- use cases
- repository interfaces
- repository implementations
- HTTP delivery layer 
- application wiring in `main.go`

## Success criteria

- routes still behave the same after refactoring
- the project is easier to test and extend
- storage and transport details are isolated from business logic

## Files

- `main.go`: refactored user API with repository, use case, and HTTP handler layers
- `PROMPT.md`: mock test prompt
- `go.mod`: module definition for the hour

## Run

```powershell
cd "D:\training_golang\Module 4 (Hour 25-32) Design Patterns & Architecture\Hour 32"
go run .
```

Test the API:

```powershell
Invoke-RestMethod -Method Post -Uri "http://127.0.0.1:8081/users" -ContentType "application/json" -Body '{"id":1,"name":"Asha"}'
Invoke-RestMethod -Method Get -Uri "http://127.0.0.1:8081/users"
Invoke-RestMethod -Method Get -Uri "http://127.0.0.1:8081/users/1"
```
