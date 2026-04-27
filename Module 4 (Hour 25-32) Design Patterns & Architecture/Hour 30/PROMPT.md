# Prompt - Hour 30 Adapter and Observer Pattern

Write Go code for a `payment processing` example that uses both the Adapter pattern and the Observer pattern.

## Requirements

- Define a target interface for payment processing, for example `Pay(amount float64) error`
- Assume there is a legacy payment provider with an incompatible method signature
- Build an adapter so the legacy provider can be used through the new target interface
- Add an observer system so listeners are notified when a payment succeeds
- Create at least two observers such as `EmailNotifier` and `AuditLogger`
- Use only the standard library
- Demonstrate the flow in `main.go`

## Output expectations

- The adapter should isolate legacy details from the rest of the application
- Observers should be attachable without changing payment processor logic
- Print a clear sequence of payment processing and observer notifications

## Goal

The final code should teach how structural and behavioral patterns can be combined in a practical Go workflow.
