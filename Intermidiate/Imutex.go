package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup
var counter int
var mutex sync.Mutex

func main() {
	wg.Add(2)
	go incrementor("Foo:")
	go incrementor("Bar:")
	wg.Wait()
	fmt.Println("Final Counter:", counter)
}

func decrementor(s string) {	
	for i := 0; i < 20; i++ {
		time.Sleep(time.Duration(rand.Intn(20)) * time.Millisecond)
		mutex.Lock()										
		counter--
		fmt.Println(s, i, "Counter:", counter)
		mutex.Unlock()										
	}	
	wg.Done()
}// go run -race main.go
// vs
// go run main.go
/*
In this code, we have two goroutines incrementing a shared counter variable. The mutex is used to ensure that only one goroutine can access the counter at a time, preventing race conditions. The `incrementor` function simulates some work by sleeping for a random duration before locking the mutex, incrementing the counter, and unlocking it. The `decrementor` function (commented out) would do the opposite, decrementing the counter in a similar fashion.	
To see the effects of race			conditions, you can run the code with the `-race` flag, which will detect
any race			
conditions and report them. Running without the `-race` flag may not show any issues, but it can lead to unpredictable behavior due to race conditions		
*/
func incrementor(s string) {
	for i := 0; i < 20; i++ {
		time.Sleep(time.Duration(rand.Intn(20)) * time.Millisecond)
		mutex.Lock()
		counter++
		fmt.Println(s, i, "Counter:", counter)
		mutex.Unlock()
	}
	wg.Done()
}

// go run -race main.go
// vs
// go run main.go
