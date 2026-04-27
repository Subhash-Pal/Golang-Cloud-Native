// The Observer Pattern is a behavioural design pattern that defines a one-to-many dependency between objects so that when one object (the subject) changes state, all its dependents (observers) are notified and updated automatically.
// The Observer Pattern is like a subscription system. One main object (the subject) keeps a list of subscribers (observers). Whenever the main object updates, it automatically notifies all subscribers about the change.

// Implementing the Observer Pattern in Golang
// 1. Define the Observer Interface

package main

import "fmt"

// Observer interface
type Observer interface {
	Update(price float64)
}

// 2. Define the Subject Interface(Observable)
// Subject interface
type Subject interface {
	Register(observer Observer)
	Deregister(observer Observer)
	Notify()
}

// Step 3 : Create the Concreate Subject (Stock)
// Stock represents the concrete subject
type Stock struct {
	observers []Observer
	price     float64
}

// Register adds an observer
func (s *Stock) Register(observer Observer) {
	s.observers = append(s.observers, observer)
}

// Deregister removes an observer
func (s *Stock) Deregister(observer Observer) {
	for i, obs := range s.observers {
		if obs == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

// Notify updates all observers
func (s *Stock) Notify() {
	for _, observer := range s.observers {
		observer.Update(s.price)
	}
}

// SetPrice updates the stock price and notifies observers
func (s *Stock) SetPrice(price float64) {
	fmt.Printf("Stock price updated to: %.2f\n", price)
	s.price = price
	s.Notify()
}

// Step 4 : Create the Concreate Observers (Investor/Traders)
// Trader is a concrete observer
type Trader struct {
	name string
}

// Update method for Trader
func (t *Trader) Update(price float64) {
	fmt.Printf("Trader %s notified: New stock price is %.2f\n", t.name, price)
}

// Step 5 : Testing the Observer Pattern

func main() {
	// Create a stock
	stock := &Stock{}

	// Create observers
	trader1 := &Trader{name: "Alice"}
	trader2 := &Trader{name: "Bob"}

	// Register observers
	stock.Register(trader1)
	stock.Register(trader2)

	// Change stock price
	stock.SetPrice(100.50)
	stock.SetPrice(102.75)

	// Deregister an observer
	stock.Deregister(trader1)

	// Change stock price again
	stock.SetPrice(105.00)
}
