/*
2. Factory Pattern
The Factory pattern provides an interface for creating objects in a superclass but allows subclasses to alter the type of objects that will be created. It’s useful when the creation logic is complex or varies based on conditions.
*/
package main

import "fmt"

// Product interface
type Product interface {
	Use() string
}

// Concrete products
type ProductA struct{}
type ProductB struct{}

func (p *ProductA) Use() string {
	return "Using Product A"
}

func (p *ProductB) Use() string {
	return "Using Product B"
}

// Factory function
func CreateProduct(productType string) Product {
	switch productType {
	case "A":
		return &ProductA{}
	case "B":
		return &ProductB{}
	default:
		return nil
	}
}

func main() {
	productA := CreateProduct("A")
	productB := CreateProduct("B")

	fmt.Println(productA.Use()) // Output: Using Product A
	fmt.Println(productB.Use()) // Output: Using Product B
}

/*
Explanation: The CreateProduct function acts as a factory that creates instances of different products (ProductA or ProductB) based on the input type.
*/
