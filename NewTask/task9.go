package main

import (
	"fmt"
	"sync"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		for i := 0; i < 5; i++ {
			ch1 <- i
		}
		for i := 5; i < 10; i++ {
			ch2 <- i
		}
		close(ch1)
		close(ch2)
	}()

	v := Merge(ch1, ch2)

	for val := range v {
		fmt.Println(val)
	}
}

func Merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	send := func(ch <-chan int) {
		for n := range ch {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, ch := range cs {
		go send(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
