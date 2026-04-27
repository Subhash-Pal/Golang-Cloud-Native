# Prompt - Hour 31 Event-Driven Internal Architecture

Write Go code that demonstrates an internal event-driven architecture for an `e-commerce` application.

## Requirements

- Define domain events such as `OrderCreated` and `InventoryReserved`
- Build a simple in-memory event bus using channels or subscriber lists
- Publish an event after an order is created
- Add at least two handlers, for example inventory update and email notification
- Keep the event bus internal to the application, not an external broker
- Show how components remain loosely coupled through events
- Use only the standard library

## Output expectations

- The event flow should be easy to follow for learners
- Each handler should react independently to the published event
- Print logs showing publication and handling order

## Goal

The final code should help learners understand how event-driven design can improve modularity even inside a single Go application.
