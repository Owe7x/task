package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	users := []User{{"aaa"}, {"bbb"}, {"ccc"}, {"ddd"}, {"eee"}}
	res, err := Do(context.Background(), users)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(res)
}

type User struct {
	Name string
}

func fetchByName(_ctx context.Context, _userName string) (int, error) {
	// Тут происходит сетевой поход, который по userName возвращает userID
	time.Sleep(10 * time.Millisecond) // Имитация сетевого похода
	return rand.Int() % 100000, nil
}

func Do(ctx context.Context, users []User) (map[string]int, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	collected := make(map[string]int)

	for _, u := range users {
		wg.Add(1)
		go func(u User) {
			defer wg.Done()

			userID, err := fetchByName(ctx, u.Name)
			if err != nil {
				return
			}
			mu.Lock()
			collected[u.Name] = userID
			mu.Unlock()

		}(u)

	}

	wg.Wait()

	return collected, nil
}
