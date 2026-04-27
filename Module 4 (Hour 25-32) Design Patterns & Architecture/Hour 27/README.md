# Hour 27 - Repository Pattern

The Repository pattern abstracts data access behind an interface so the application can work with domain objects without depending on storage details.

## Key concepts

- repositories expose collection-like operations for domain objects
- business logic should depend on interfaces, not concrete storage
- implementations can be swapped without changing use cases
- repositories often improve testing because mocks or in-memory versions are easy to provide

## Learning goals

- design repository interfaces around business needs
- separate storage logic from service logic
- understand how the pattern supports clean architecture
- prepare for dependency injection in the next lesson

## Suggested deliverable

Create a user management example with a repository interface, an in-memory implementation, and a service layer that performs validation.

## Files

- `main.go`: repository pattern demo with duplicate email validation
- `PROMPT.md`: coding prompt for practice
- `go.mod`: module definition for the hour

## Run

```powershell
cd "D:\training_golang\Module 4 (Hour 25-32) Design Patterns & Architecture\Hour 27"
go run .
```
