# Hour 29 - Factory and Strategy Pattern

This hour combines two useful patterns. Factory helps create the right object, while Strategy helps choose the right behavior at runtime.

## Key concepts

- Factory centralizes object creation logic
- Strategy encapsulates interchangeable algorithms
- the two patterns often work well together when behavior depends on configuration or user choice
- the service code stays simpler because it depends on abstractions

## Learning goals

- know when to move branching logic into a factory
- design interchangeable algorithms behind a shared interface
- combine creation and behavior patterns in one example
- improve maintainability by reducing `switch` logic in core services

## Suggested deliverable

Build an order pricing example where a factory selects a pricing strategy and a service uses it to calculate the final amount.

## Files

- `main.go`: Factory plus Strategy pattern demo for order pricing
- `PROMPT.md`: coding prompt for practice
- `go.mod`: module definition for the hour

## Run

```powershell
cd "D:\training_golang\Module 4 (Hour 25-32) Design Patterns & Architecture\Hour 29"
go run .
```
