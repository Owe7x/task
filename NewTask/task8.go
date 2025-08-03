package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func main() {
	sliceUrl := []string{
		"https://google.com",
		"https://facebook.com",
		"https://golang.org",
		"https://yandex.ru",
	}

	mapUrl := FetchURLs(sliceUrl)

	fmt.Println(mapUrl)
}

func FetchURLs(slice []string) map[string]string {
	const workerPool = 10

	var (
		mu          sync.Mutex
		wg          sync.WaitGroup
		ctx, cancel = context.WithCancel(context.Background())
	)

	defer cancel()

	results := make(map[string]string)
	jobs := make(chan string)

	// воркеры
	for i := 0; i < workerPool; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := http.Client{Timeout: 5 * time.Second}

			for url := range jobs {
				req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
				if err != nil {
					fmt.Println("Ошибка создания запроса:", err)
					cancel()
					return
				}

				resp, err := client.Do(req)
				if err != nil {
					fmt.Println("Ошибка запроса или контекст отменён:", err)
					cancel()
					return
				}

				// закрываем с обработкой ошибки
				func() {
					defer func() {
						if cerr := resp.Body.Close(); cerr != nil {
							fmt.Println("Ошибка закрытия тела ответа:", cerr)
						}
					}()

					buf := make([]byte, 100)
					n, readErr := resp.Body.Read(buf)
					if readErr != nil && readErr != io.EOF {
						fmt.Println("Ошибка чтения тела:", readErr)
						cancel()
						return
					}

					mu.Lock()
					results[url] = fmt.Sprintf("Status: %d, Body: %s", resp.StatusCode, string(buf[:n]))
					mu.Unlock()
				}()
			}
		}()
	}

	// подача задач
	go func() {
		for _, url := range slice {
			select {
			case <-ctx.Done():
				return // сразу выйти, а не break
			case jobs <- url:
			}
		}
		close(jobs)
	}()

	wg.Wait()
	return results
}
