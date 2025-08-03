package main

import (
	"fmt"
	"sync"
)

// Асихнронное слияние каналов

func main() {
	ch1 := make(chan int, 10)
	ch2 := make(chan int, 20)

	ch1 <- 1
	ch2 <- 2
	ch2 <- 4

	close(ch1)
	close(ch2)

	ch3 := asyncMerge[int](ch1, ch2)

	for val := range ch3 {
		fmt.Println(val)
	}
}

// todo
func asyncMerge[T any](cs ...<-chan T) <-chan T {
	ch := make(chan T)
	var wg sync.WaitGroup

	for _, val := range cs {
		wg.Add(1)
		go func(val <-chan T) {
			defer wg.Done()
			for v := range val {
				ch <- v
			}
		}(val)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}
