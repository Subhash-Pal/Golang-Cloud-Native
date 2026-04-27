// The Decorator Design Pattern is a structural pattern because it helps flexibly structure objects and classes.
// It allows behaviour to be added to individual objects without modifying the original object or other instances of the same class.

// Example: Coffee Pricing System

// Step 1: Define the Component Interface
package main

import "fmt"

// Component Interface
type Coffee interface {
	Cost() float64
	Description() string
}

// Step 2: Create a Concrete Component (Base Coffee)
// Concrete Component
type SimpleCoffee struct{}

func (s *SimpleCoffee) Cost() float64 {
	return 5.0
}

func (s *SimpleCoffee) Description() string {
	return "Simple Coffee"
}

// Step 3: Create Decorators

// Decorator Base
type CoffeeDecorator struct {
	Coffee Coffee
}

func (c *CoffeeDecorator) Cost() float64 {
	return c.Coffee.Cost()
}

func (c *CoffeeDecorator) Description() string {
	return c.Coffee.Description()
}

// Concrete Decorator: Milk
type MilkDecorator struct {
	Coffee Coffee
}

func (m *MilkDecorator) Cost() float64 {
	return m.Coffee.Cost() + 1.5 // Milk adds $1.5
}

func (m *MilkDecorator) Description() string {
	return m.Coffee.Description() + ", Milk"
}

// Concrete Decorator: Sugar
type SugarDecorator struct {
	Coffee Coffee
}

func (s *SugarDecorator) Cost() float64 {
	return s.Coffee.Cost() + 0.5 // Sugar adds $0.5
}

func (s *SugarDecorator) Description() string {
	return s.Coffee.Description() + ", Sugar"
}

// Step 4 : Using decorators
func main() {
	// Base Coffee
	coffee := &SimpleCoffee{}
	fmt.Println(coffee.Description(), "Cost:", coffee.Cost())

	// Coffee with Milk
	coffeeWithMilk := &MilkDecorator{Coffee: coffee}
	fmt.Println(coffeeWithMilk.Description(), "Cost:", coffeeWithMilk.Cost())

	// Coffee with Milk and Sugar
	coffeeWithMilkAndSugar := &SugarDecorator{Coffee: coffeeWithMilk}
	fmt.Println(coffeeWithMilkAndSugar.Description(), "Cost:", coffeeWithMilkAndSugar.Cost())
}
