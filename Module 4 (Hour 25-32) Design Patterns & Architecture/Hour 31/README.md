# Hour 31 - Event-Driven Internal Architecture

Event-driven internal architecture allows components within the same application to communicate through events instead of direct calls.

## Key concepts

- events represent important things that happened in the system
- publishers emit events without knowing who handles them
- subscribers react to events independently
- internal event buses improve modularity without requiring external messaging systems

## Learning goals

- model domain events in a Go application
- decouple workflows using publish-subscribe style communication
- understand when an internal event bus is enough
- prepare for larger distributed event-driven systems later in the course

## Suggested deliverable

Build a small e-commerce flow where order creation publishes an internal event that triggers inventory and notification handlers.

## Files

- `main.go`: internal event bus demo with order creation handlers
- `PROMPT.md`: coding prompt for practice
- `go.mod`: module definition for the hour

## Run

```powershell
cd "D:\training_golang\Module 4 (Hour 25-32) Design Patterns & Architecture\Hour 31"
go run .
```
