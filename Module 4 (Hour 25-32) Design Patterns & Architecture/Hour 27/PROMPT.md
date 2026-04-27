# Prompt - Hour 27 Repository Pattern

Write Go code that demonstrates the Repository pattern for a `user management` service.

## Requirements

- Define a `User` entity with fields such as `ID`, `Name`, and `Email`
- Create a repository interface with methods like `Create`, `FindByID`, `FindAll`, and `Delete`
- Implement an in-memory repository using a map and mutex for safety
- Add a service or use case layer that depends on the repository interface
- Validate duplicate email addresses before saving
- Use only the standard library
- Show usage in `main.go`

## Output expectations

- The repository should hide storage details from the business logic
- The example should be easy to replace later with a database-backed implementation
- Print meaningful output for create and list operations

## Goal

The final code should make clear why repositories improve testability and reduce direct coupling to storage mechanisms.
