package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// Задача на контекст отмены и тикеры(Авито)

func init() {
	rand.NewSource(time.Now().UnixNano())
}

// Есть функция, работающая неопределённо долго и возвращающая число.
// Её тело нельзя изменять (представим, что внутри сетевой запрос) .
func unpredictableFunc() int64 {
	rnd := rand.Int63n(5000)
	time.Sleep(time.Duration(rnd) * time.Millisecond)

	return rnd

}

// Нужно изменить функцию обёртку, которая будет работать с заданным таймаутом (например, 1 секунду).
// Если "длинная" функция отработала за это время - отлично, возвращаем результат.
// Если нет - возвращаем ошибку. Результат работы в этом случае нам не важен.
//
// Дополнительно нужно измерить, сколько выполнялась эта функция (просто вывести
// Сигнатуру функцию обёртки менять можно.
func predictableFunc(ctx context.Context) int64 {
	timer := time.NewTimer(1 * time.Second)
	now := time.Now()

	ch := make(chan int64)
	res := int64(0)

	go func() {
		fmt.Println("Go work 1")
		ch <- unpredictableFunc()
		close(ch)
	}()

	select {

	case <-timer.C:
		fmt.Println("timed out")
		return 0
	case <-ctx.Done():
		fmt.Println("cancelled")
		return 0
	case v, ok := <-ch:
		if !ok {
			return 0
		}
		res = v
	}

	after := time.Since(now)

	fmt.Println(after)

	return res
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	fmt.Println("started")
	fmt.Println(predictableFunc(ctx))
	fmt.Println("finished")
}
