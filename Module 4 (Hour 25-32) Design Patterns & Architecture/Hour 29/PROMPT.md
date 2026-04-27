# Prompt - Hour 29 Factory and Strategy Pattern

Write Go code for an `order pricing` system that uses both the Factory pattern and the Strategy pattern.

## Requirements

- Define a pricing strategy interface with a method like `Calculate(total float64) float64`
- Implement multiple strategies such as `RegularPricing`, `FestivalDiscount`, and `PremiumMemberDiscount`
- Use a factory function to return the correct strategy based on an input type
- Add an `OrderService` that applies the chosen strategy to an order amount
- Handle invalid strategy names gracefully
- Use only the Go standard library
- Include a runnable `main.go` with multiple examples

## Output expectations

- The code should clearly separate object creation from behavior selection
- Strategies should be interchangeable without changing the service logic
- Print final prices for different order types

## Goal

The final code should show how Factory and Strategy work together to choose and execute runtime behavior cleanly.
