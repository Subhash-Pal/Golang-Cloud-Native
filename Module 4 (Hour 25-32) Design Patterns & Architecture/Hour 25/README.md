# Hour 25 - Clean Architecture Theory

This hour introduces the theory behind Clean Architecture in Go. The main idea is to keep core business rules independent from frameworks, databases, and transport details.

## Key concepts

- entities contain core business data and rules
- use cases coordinate application behavior
- interfaces define contracts for external dependencies
- infrastructure provides implementations such as in-memory storage, databases, or HTTP handlers
- dependencies should point inward toward stable business abstractions

## Learning goals

- understand why architecture boundaries reduce coupling
- identify the role of entities, use cases, and adapters
- separate business logic from framework-specific code
- prepare for folder-based implementation in the next hour

## Suggested deliverable

Create a small task management example showing entities, repository interfaces, and use cases with dependency inversion.

## Files

- `main.go`: runnable task management demo with entity, repository abstraction, use case, and in-memory implementation
- `PROMPT.md`: coding prompt for practice
- `go.mod`: module definition for the hour

## Run

```powershell
cd "D:\training_golang\Module 4 (Hour 25-32) Design Patterns & Architecture\Hour 25"
go run .
```
