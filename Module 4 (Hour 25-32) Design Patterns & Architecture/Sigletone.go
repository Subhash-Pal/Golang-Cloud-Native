package main

import "fmt"

// Singleton struct
type Singleton struct{}

var instance *Singleton

// GetInstance ensures only one instance of Singleton is created
func GetInstance() *Singleton {
	if instance == nil {
		instance = &Singleton{}
	}
	return instance
}

func main() {
	s1 := GetInstance()
	s2 := GetInstance()

	// Both s1 and s2 point to the same instance
	fmt.Println(s1 == s2) // Output: true
}

/*
1. Singleton Pattern
The Singleton pattern ensures that a class has only one instance and provides a global point of access to it. This is useful when exactly one object is needed to coordinate actions across the system.
Explanation: The GetInstance function ensures that only one instance of the Singleton struct is created. Subsequent calls return the same instance.
2. Lazy Initialization
The instance of the Singleton is created only when it is needed, which can save resources if the instance is never used.
Explanation: The instance variable is initialized to nil, and the actual instance is created only when GetInstance is called for the first time.
*/
