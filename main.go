package main

import (
	"fmt"
	"sync"
)

func runGoroutine() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		fmt.Println("Hello from goroutine!")
	}()

	wg.Wait()
}

func main() {
	runGoroutine()
}
