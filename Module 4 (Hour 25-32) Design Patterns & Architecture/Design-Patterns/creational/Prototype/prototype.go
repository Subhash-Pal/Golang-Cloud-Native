// The Prototype Design Pattern is a creational pattern that allows objects to be cloned instead of creating new instances from scratch. This is useful when object creation is costly or complex. In Golang, we achieve this by implementing a Clone method for the prototype object.

// Implementation of Prototype Pattern in Golang
// Let's consider a scenario where we have a Shape interface with different types of shapes (e.g., Circle and Rectangle). Each shape should implement the Clone method.

// Step 1: Define the Prototype Interface
package main

import "fmt"

// Step 2: Implement Concrete Prototypes
// Prototype interface
type Shape interface {
	Clone() Shape
	GetInfo() string
}

// Concrete type: Circle
type Circle struct {
	Radius int
	Color  string
}

// Clone method for Circle
func (c *Circle) Clone() Shape {
	return &Circle{Radius: c.Radius, Color: c.Color}
}

// GetInfo method for Circle
func (c *Circle) GetInfo() string {
	return fmt.Sprintf("Circle: Radius = %d, Color = %s", c.Radius, c.Color)
}

// Concrete type: Rectangle
type Rectangle struct {
	Width  int
	Height int
	Color  string
}

// Clone method for Rectangle
func (r *Rectangle) Clone() Shape {
	return &Rectangle{Width: r.Width, Height: r.Height, Color: r.Color}
}

// GetInfo method for Rectangle
func (r *Rectangle) GetInfo() string {
	return fmt.Sprintf("Rectangle: Width = %d, Height = %d, Color = %s", r.Width, r.Height, r.Color)
}

// Step 3: Using the Prototype Pattern
func main() {
	// Create an original circle
	originalCircle := &Circle{Radius: 10, Color: "Red"}
	fmt.Println("Original:", originalCircle.GetInfo())

	// Clone the circle
	clonedCircle := originalCircle.Clone()
	fmt.Println("Cloned:", clonedCircle.GetInfo())

	// Modify the cloned circle
	clonedCircle.(*Circle).Color = "Blue"
	fmt.Println("Modified Clone:", clonedCircle.GetInfo())
	fmt.Println("Original Unchanged:", originalCircle.GetInfo())

	// Create and clone a rectangle
	originalRectangle := &Rectangle{Width: 20, Height: 10, Color: "Green"}
	clonedRectangle := originalRectangle.Clone()
	fmt.Println("Original Rectangle:", originalRectangle.GetInfo())
	fmt.Println("Cloned Rectangle:", clonedRectangle.GetInfo())
}
