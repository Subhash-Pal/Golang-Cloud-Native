// The Strategy Pattern is a behavioral design pattern that allows you to define a family of algorithms, encapsulate them in separate structs, and make them interchangeable at runtime. It promotes composition over inheritance and is useful when you need to select an algorithm dynamically.

// Step 1: Define the Strategy Interface
package main

import "fmt"

// PaymentStrategy defines a common interface for different payment methods
type PaymentStrategy interface {
	Pay(amount float64)
}

// Step 2: Implement Different Payment Strategies
// CreditCardPayment strategy
type CreditCardPayment struct{}

func (c *CreditCardPayment) Pay(amount float64) {
	fmt.Printf("Paid $%.2f using Credit Card\n", amount)
}

// PayPalPayment strategy
type PayPalPayment struct{}

func (p *PayPalPayment) Pay(amount float64) {
	fmt.Printf("Paid $%.2f using PayPal\n", amount)
}

// BitcoinPayment strategy
type BitcoinPayment struct{}

func (b *BitcoinPayment) Pay(amount float64) {
	fmt.Printf("Paid $%.2f using Bitcoin\n", amount)
}

// Step 3: Context to Use Strategies
// PaymentContext holds a reference to a payment strategy
type PaymentContext struct {
	strategy PaymentStrategy
}

// SetStrategy allows changing the payment strategy at runtime
func (p *PaymentContext) SetStrategy(strategy PaymentStrategy) {
	p.strategy = strategy
}

// ExecutePayment processes the payment using the selected strategy
func (p *PaymentContext) ExecutePayment(amount float64) {
	if p.strategy == nil {
		fmt.Println("No payment strategy selected!")
		return
	}
	p.strategy.Pay(amount)
}

// Step 4: Use the Strategy Pattern
func main() {
	context := &PaymentContext{}

	// Use Credit Card Payment
	context.SetStrategy(&CreditCardPayment{})
	context.ExecutePayment(100.00)

	// Switch to PayPal Payment
	context.SetStrategy(&PayPalPayment{})
	context.ExecutePayment(50.00)

	// Switch to Bitcoin Payment
	context.SetStrategy(&BitcoinPayment{})
	context.ExecutePayment(200.00)
}
