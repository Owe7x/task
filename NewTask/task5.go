package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {

	wg := sync.WaitGroup{}

	var counter int64

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}

	wg.Wait()

	fmt.Println("Counter:", counter)

}
