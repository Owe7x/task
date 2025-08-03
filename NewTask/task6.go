package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	input := make(chan int)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go StartBatchProcessor(ctx, input)

	go func() {
		for i := 1; i <= 20; i++ {
			input <- i
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
	}()

	<-ctx.Done()
	fmt.Println("Main: processing stopped")
}

func StartBatchProcessor(ctx context.Context, input chan int) {
	slice := make([]int, 0)
	timer := time.NewTimer(2 * time.Second)

	for {
		select {
		case i, ok := <-input:
			if !ok {
				fmt.Println("Channel closed", slice)
				return
			}

			slice = append(slice, i)

			if len(slice) == 5 {
				fmt.Println("Batch process:", slice)
				slice = []int{}
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(2 * time.Second)
			}
		case <-timer.C:
			if len(slice) > 0 {
				fmt.Println("TimeOut. Final batch:", slice)
				slice = []int{}

			}
			timer.Reset(2 * time.Second)
		case <-ctx.Done():
			return
		}
	}
}
