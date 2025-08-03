package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Задача Оптимизация через горутины и worker pool

func main() {
	rand.NewSource(time.Now().UnixNano())

	const (
		workerPool = 300
		maxId      = 100000
		idByf      = 100
		TimeOut    = 10 * time.Second
	)

	var wg sync.WaitGroup

	jobs := make(chan int, idByf)

	ctx, timeOutCancel := context.WithTimeout(context.Background(), TimeOut)
	defer timeOutCancel()

	url := "https://www.wildberries.ru/getPurchasesByID?id=%d"

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(jobs)
		for i := 0; i < maxId; i++ {
			select {
			case jobs <- i:
			case <-ctx.Done():
				return
			}
		}

	}()

	for i := 0; i < workerPool; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range jobs {
				resp := getPurchasesByID(ctx, url, id)
				if resp == "" {
					continue
				}
				fmt.Println(resp)
			}
		}()
	}
	wg.Wait()
}

func getPurchasesByID(ctx context.Context, urlTpl string, id int) string {
	delay := time.Duration(rand.Intn(3000)) * time.Millisecond

	select {
	case <-time.After(delay):
		fullUrl := fmt.Sprintf(urlTpl, id)
		return fmt.Sprintf("id: %d url: %s", id, fullUrl)
	case <-ctx.Done():
		return ""

	}
}
