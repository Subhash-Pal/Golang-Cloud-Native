# Prompt - Hour 25 Clean Architecture Theory

Write Go code that demonstrates the core ideas of Clean Architecture for a simple `task management` API.

## Requirements

- Create clear layers for `domain`, `use case`, `interface`, and `infrastructure`
- Define an entity such as `Task` with fields like `ID`, `Title`, `Completed`, and `CreatedAt`
- Add one use case such as `CreateTask` and one use case such as `ListTasks`
- Define repository interfaces in the business layer, not in the infrastructure layer
- Show dependency direction from outer layers toward abstractions in inner layers
- Use only the Go standard library
- Include a small `main.go` that wires dependencies and runs a basic example

## Output expectations

- Keep the code beginner friendly and well organized
- Add short comments only where the architecture boundary needs explanation
- Print a sample result from the use case execution

## Goal

The final code should help a learner understand how Clean Architecture separates business rules from delivery and storage details.
