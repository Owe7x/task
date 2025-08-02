package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	const n = 5

	for i := 1; i <= n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(i)
		}()
	}

	wg.Wait()
}
