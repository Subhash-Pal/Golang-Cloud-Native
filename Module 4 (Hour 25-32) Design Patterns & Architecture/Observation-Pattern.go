/*
3. Observer Pattern
The Observer pattern defines a one-to-many dependency between objects so that when one object changes state, all its dependents are notified and updated automatically. It’s commonly used in event-driven systems.
*/
package main

import "fmt"

// Observer interface
type Observer interface {
	Update(string)
}

// Subject (Observable)
type Subject struct {
	observers []Observer
	state     string
}

func (s *Subject) Attach(observer Observer) {
	s.observers = append(s.observers, observer)
}

func (s *Subject) SetState(newState string) {
	s.state = newState
	s.NotifyAll()
}

func (s *Subject) NotifyAll() {
	for _, observer := range s.observers {
		observer.Update(s.state)
	}
}

// Concrete Observer
type ConcreteObserver struct {
	name string
}

func (o *ConcreteObserver) Update(state string) {
	fmt.Printf("%s received update: %s\n", o.name, state)
}

func main() {
	subject := &Subject{}
	observer1 := &ConcreteObserver{name: "Observer 1"}
	observer2 := &ConcreteObserver{name: "Observer 2"}

	subject.Attach(observer1)
	subject.Attach(observer2)

	subject.SetState("State 1") // Output: Observer 1/2 received update: State 1
	subject.SetState("State 2") // Output: Observer 1/2 received update: State 2
}

/*

Explanation: The Subject maintains a list of observers and notifies them whenever its state changes. Observers implement the Update method to react to state changes.
*/
