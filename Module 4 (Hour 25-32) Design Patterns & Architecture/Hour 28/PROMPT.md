# Prompt - Hour 28 Dependency Injection

Write Go code that shows Dependency Injection in a `notification service`.

## Requirements

- Define an interface such as `Sender` with a method like `Send(to, message string) error`
- Create at least two implementations, for example `EmailSender` and `SMSSender`
- Build a `NotificationService` that receives its dependencies through constructor injection
- Add a logger abstraction and inject that as well
- Demonstrate swapping implementations without changing service logic
- Use only standard library packages
- Include a runnable `main.go`

## Output expectations

- Constructor functions should make dependencies explicit
- The service should not instantiate its own concrete senders
- Show at least two examples using different injected implementations

## Goal

The final code should teach how dependency injection improves flexibility, testability, and separation of concerns in Go.
