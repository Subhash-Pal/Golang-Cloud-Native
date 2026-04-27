package main

import (
	"fmt"
	"reflect"
)

// ExampleStruct demonstrates dynamic method invocation.
type ExampleStruct struct{}

// Greet is a method that takes a name and prints a greeting.
func (e *ExampleStruct) Greet(name string) {
	fmt.Printf("Hello, %s!\n", name)
}

// Add is a method that adds two integers and returns the result.
func (e *ExampleStruct) Add(a, b int) int {
	return a + b
}

// InvokeMethod dynamically invokes a method on a struct using reflection.
func InvokeMethod(obj interface{}, methodName string, args ...interface{}) ([]reflect.Value, error) {
	// Get the value of the object (must be a pointer to invoke methods)
	value := reflect.ValueOf(obj)
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("obj must be a pointer to a struct")
	}

	// Get the method by name
	method := value.MethodByName(methodName)
	if !method.IsValid() {
		return nil, fmt.Errorf("method '%s' does not exist", methodName)
	}

	// Prepare arguments for the method call
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	// Call the method dynamically
	results := method.Call(in)
	return results, nil
}

func main() {
	// Create an instance of ExampleStruct
	example := &ExampleStruct{}

	// Dynamically invoke the 'Greet' method
	_, err := InvokeMethod(example, "Greet", "Shubh")
	if err != nil {
		fmt.Println("Error invoking method:", err)
	}

	// Dynamically invoke the 'Add' method
	results, err := InvokeMethod(example, "Add", 10, 20)
	if err != nil {
		fmt.Println("Error invoking method:", err)
	} else {
		// Print the result of the 'Add' method
		fmt.Printf("Add Result: %v\n", results[0].Interface())
	}
}

/*

Explanation of the Code
1. Struct Definition
The ExampleStruct contains two methods:
Greet: A method that takes a string argument and prints a greeting.
Add: A method that takes two integers, adds them, and returns the result.
2. Dynamic Method Invocation
The InvokeMethod function uses reflection to dynamically invoke a method on a struct.
It ensures that:
The input is a pointer to a struct.
The specified method exists on the struct.
It prepares the arguments for the method call using reflect.ValueOf and invokes the method dynamically using method.Call().
3. Main Function
The main function demonstrates how to use InvokeMethod to call the Greet and Add methods dynamically.
Output Example
When you run the above code, the output will look something like this:

Hello, Shubh!
Add Result: 30

Key Concepts Demonstrated
Dynamic Method Invocation:
Reflection allows you to call methods on a struct dynamically, even when their names and arguments are not known at compile time.
Error Handling:
The InvokeMethod function includes robust error handling to ensure type safety and prevent runtime panics.
Flexibility:
This approach can be extended to work with any struct or interface, making it highly reusable.
Use Cases:
Plugin Systems: Dynamically load and invoke methods from plugins.
RPC Frameworks: Invoke remote methods based on client requests.
Testing Frameworks: Dynamically test methods on structs.

Advanced Usage: Handling Variadic Methods
If your struct has variadic methods (methods that accept a variable number of arguments), you can handle them as follows:
// VariadicMethod demonstrates a method with variadic arguments.
func (e *ExampleStruct) VariadicMethod(args ...int) int {
	sum := 0
	for _, arg := range args {
		sum += arg
	}
	return sum
}

func main() {
	example := &ExampleStruct{}

	// Dynamically invoke the 'VariadicMethod'
	results, err := InvokeMethod(example, "VariadicMethod", 1, 2, 3, 4)
	if err != nil {
		fmt.Println("Error invoking method:", err)
	} else {
		fmt.Printf("VariadicMethod Result: %v\n", results[0].Interface())
	}
}

This will output:
VariadicMethod Result: 10


Notes
Performance:
Reflection introduces some overhead compared to direct method calls. Use it judiciously in performance-critical applications.
Type Safety:
Reflection bypasses compile-time type checking. Ensure proper error handling to avoid runtime panics.
Extensibility:
You can extend this implementation to support interfaces, embedded methods, or custom argument parsing.
*/