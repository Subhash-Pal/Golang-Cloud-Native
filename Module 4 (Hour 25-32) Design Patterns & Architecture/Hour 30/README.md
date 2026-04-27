# Hour 30 - Adapter and Observer Pattern

This lesson combines integration and event notification patterns. Adapter helps incompatible interfaces work together, while Observer helps the system react to events without tight coupling.

## Key concepts

- Adapter wraps an existing component and exposes a new interface
- Observer supports one-to-many notifications when state changes or events occur
- these patterns are useful when integrating legacy systems and side effects such as logging or notifications

## Learning goals

- adapt a legacy dependency into a cleaner interface
- notify multiple subscribers after an action completes
- keep side effects separate from the core payment flow
- understand how these patterns support extensible system design

## Suggested deliverable

Create a payment processor that adapts a legacy gateway and notifies observers when a payment is completed.

## Files

- `main.go`: Adapter plus Observer pattern demo for payments
- `PROMPT.md`: coding prompt for practice
- `go.mod`: module definition for the hour

## Run

```powershell
cd "D:\training_golang\Module 4 (Hour 25-32) Design Patterns & Architecture\Hour 30"
go run .
```
