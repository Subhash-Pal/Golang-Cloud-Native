package main

/*
#include <stdlib.h>

int addNumbers(int a, int b) {
	return a + b;
}

int multiplyNumbers(int a, int b) {
	return a * b;
}
*/
import "C"

import "fmt"

func main() {
	a := 10
	b := 5

	sum := C.addNumbers(C.int(a), C.int(b))
	product := C.multiplyNumbers(C.int(a), C.int(b))

	fmt.Println("CGO Integration Example")
	fmt.Printf("Addition using C function: %d + %d = %d\n", a, b, int(sum))
	fmt.Printf("Multiplication using C function: %d * %d = %d\n", a, b, int(product))
}
