package main
import (
	"context"
	"fmt"
	"time"
)
func Worker(ctx context.Context){
	for{
		select{	case <-ctx.Done():
			fmt.Println("Worker received cancellation signal, exiting...",ctx.Err()	)
				return
		default:
			fmt.Println("Worker is doing some work...")
			time.Sleep(1*time.Second)
		}
	}
}
func main(){
	ctx,cancel:=context.WithCancel(context.Background())
	go Worker(ctx)
	time.Sleep(5 * time.Second)
	cancel()//->Done ->Stop the worker 
	time.Sleep(2 * time.Second)
	fmt.Println("Main function exiting...")
}