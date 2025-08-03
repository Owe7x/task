package main

import (
	"fmt"
	"sync"
)

func main() {
	var mu sync.Mutex
	wg := sync.WaitGroup{}

	counter := 0

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()

	fmt.Println("Counter:", counter)

}
