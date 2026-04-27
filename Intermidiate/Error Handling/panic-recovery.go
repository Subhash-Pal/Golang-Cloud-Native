package main

import "fmt"

func main() {
    // 1. Defer a function to handle potential panics
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from panic:", r)
        }
    }()

    fmt.Println("Starting division...")
    
    // 2. Trigger a panic (division by zero)
    result := divide(10, 0)
    
    // This line will never execute because of the panic
    fmt.Println("Result:", result)
}

func divide(a, b int) int {
    if b == 0 {
        panic("cannot divide by zero") // Manual panic
    }
    return a / b
}
