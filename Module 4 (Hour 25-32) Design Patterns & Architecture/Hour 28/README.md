# Hour 28 - Dependency Injection

Dependency Injection means providing a component with the dependencies it needs instead of letting it create them internally.

## Key concepts

- constructor injection is the most common and clear form in Go
- services should depend on interfaces when behavior may vary
- explicit dependencies make code easier to test and reason about
- dependency injection is often used together with repositories, services, and transport handlers

## Learning goals

- identify tightly coupled code
- refactor service creation to use injected dependencies
- swap implementations without editing business logic
- understand how dependency injection supports clean architecture

## Suggested deliverable

Create a notification service that works with different senders and loggers by injecting abstractions through constructors.

## Files

- `main.go`: dependency injection demo with email and SMS senders
- `PROMPT.md`: coding prompt for practice
- `go.mod`: module definition for the hour

## Run

```powershell
cd "D:\training_golang\Module 4 (Hour 25-32) Design Patterns & Architecture\Hour 28"
go run .
```
