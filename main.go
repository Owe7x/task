package main

import (
	"fmt"
	"sync"
)

func runGoroutine(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Hello from goroutine!")
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go runGoroutine(&wg)
	wg.Wait()
}
