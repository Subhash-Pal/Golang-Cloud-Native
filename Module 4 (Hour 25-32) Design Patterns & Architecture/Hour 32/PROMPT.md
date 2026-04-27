# Prompt - Hour 32 Mock Test Refactor Existing API to Clean Architecture

Refactor an existing Go API into Clean Architecture.

## Scenario

Assume you are given a small API where routing, validation, business logic, and storage are mixed in one file. Your task is to reorganize it into clear architecture layers.

## Requirements

- keep the external API behavior the same
- extract domain entities and use cases
- define repository interfaces for persistence
- move HTTP handler code into a delivery layer
- add an in-memory repository implementation for demonstration
- use dependency injection to wire the application
- keep the project runnable with the standard library only
- include a short explanation in comments or README of what was improved

## Expected endpoints

- `POST /users`
- `GET /users`
- `GET /users/{id}`

## Goal

The final code should demonstrate your understanding of architecture boundaries, repositories, and dependency injection by turning a tightly coupled API into a maintainable Go project.
