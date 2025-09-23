package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const T1 = 500 * time.Millisecond
const T2 = 5 * time.Second

func goroutineFunc(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		// по истечению таймаута контекст вернёт закрытый канал
		// чтение из закрытого канала - nil value => мы провалимся в этот кейс и завершим нашу горутину
		case <-ctx.Done():
			fmt.Println("canceled by context")
			return
		// имитация работы
		default:
			time.Sleep(T1)
		}
	}
}

func main() {
	// можно было бы сделать обычный context.WithCancel и напрямую вызвать его отмену cancel()

	// создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), T2)
	defer cancel()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go goroutineFunc(ctx, wg)

	wg.Wait()
}
