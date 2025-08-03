package main

import (
	"fmt"
	"sync"
)

func main() {
	ch1 := make(chan int)
	const workerPool = 10
	var wg sync.WaitGroup

	go func() {
		for i := 0; i < 50; i++ {
			ch1 <- i
		}
		close(ch1)
	}()

	outs := fanOut(ch1, workerPool)

	wg.Add(len(outs))

	for i, out := range outs {
		go func(id int, c <-chan int) {
			defer wg.Done()
			for v := range c {
				fmt.Printf("worker %d output: %d\n", id+1, v)
			}
		}(i, out)
	}

	wg.Wait()
}

func fanOut(ch1 <-chan int, workerPool int) []<-chan int {
	outs := make([]<-chan int, workerPool)
	for i := 0; i < workerPool; i++ {
		outs[i] = pipeline(ch1)
	}
	return outs
}

func pipeline(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range in {
			out <- v * 2
		}
	}()
	return out
}
