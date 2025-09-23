package main

import (
	"fmt"
	"sync"
	"time"
)

const T = 500 * time.Millisecond
const N = 5

func goroutineFunc(wg *sync.WaitGroup, ch <-chan int) {
	defer wg.Done()
	for {
		_, ok := <-ch
		// если канал закрыт, то ok == false => выполнится условие ниже
		if !ok {
			fmt.Println("canceled by closed channel")
			return
		}

		time.Sleep(T)
	}
}

func main() {
	ch := make(chan int)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go goroutineFunc(wg, ch)

	// в цикле пишем несколько значений в канал
	for i := range N {
		ch <- i
	}

	// закрываем канал
	close(ch)
	wg.Wait()
}
