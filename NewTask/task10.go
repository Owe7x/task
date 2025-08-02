package main

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

}
