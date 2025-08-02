package main

import (
	"fmt"
)

func main() {
	sum := 0
	ch := make(chan int)

	go func() {
		sendToCh(ch)
		close(ch)
	}()

	for v := range ch {
		sum += v
	}

	fmt.Println(sum)

}

func sendToCh(ch chan<- int) {

	for i := 0; i < 5; i++ {
		ch <- i
	}
}
