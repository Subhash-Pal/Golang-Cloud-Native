// Implementation of the Builder Pattern in Golang
// The Builder design pattern is a creational pattern used to construct complex objects step by step. It provides a clear and readable way to create objects with many optional parameters without using a long constructor with numerous arguments.

// When to Use the Builder Pattern
// When an object has multiple optional parameters.
// When you need to create different variations of an object.
// When you want to make object creation readable and maintainable.

package main

import "fmt"

// Car is the complex object we want to build.
type Car struct {
	Brand  string
	Model  string
	Color  string
	Engine string
	Wheels int
}

// CarBuilder is the builder struct.
type CarBuilder struct {
	car Car
}

// NewCarBuilder creates a new instance of CarBuilder.
func NewCarBuilder() *CarBuilder {
	return &CarBuilder{}
}

// SetBrand sets the brand of the car.
func (b *CarBuilder) SetBrand(brand string) *CarBuilder {
	b.car.Brand = brand
	return b
}

// SetModel sets the model of the car.
func (b *CarBuilder) SetModel(model string) *CarBuilder {
	b.car.Model = model
	return b
}

// SetColor sets the color of the car.
func (b *CarBuilder) SetColor(color string) *CarBuilder {
	b.car.Color = color
	return b
}

// SetEngine sets the engine type of the car.
func (b *CarBuilder) SetEngine(engine string) *CarBuilder {
	b.car.Engine = engine
	return b
}

// SetWheels sets the number of wheels.
func (b *CarBuilder) SetWheels(wheels int) *CarBuilder {
	b.car.Wheels = wheels
	return b
}

// Build finalizes and returns the constructed Car.
func (b *CarBuilder) Build() Car {
	return b.car
}

func main() {

	// Using the builder to create a car object.
	carBuilder := NewCarBuilder()

	carBuilder = carBuilder.SetBrand("Tesla")
	carBuilder = carBuilder.SetModel("Model S")
	carBuilder = carBuilder.SetColor("Red")
	carBuilder = carBuilder.SetEngine("Electric")
	carBuilder = carBuilder.SetWheels(4)

	car := carBuilder.Build()

	fmt.Printf("Car Details: %+v\n", car)
}
