package main

import (
	"fmt"
	"sync"
	"time"
)

const N = 5                      // кол-во итераций в цикле for
const T = 500 * time.Millisecond // время имитации работы горутины

func goroutineFunc(wg *sync.WaitGroup) {
	defer wg.Done()
	counter := 0
	// пока counter < N горутина будет работать
	for counter < N {
		counter++
		time.Sleep(T)
	}

	fmt.Println("completed by condition")
}

func main() {
	// wait group'а здесь необходима для того, чтобы поджоинить основную горутину с goroutineFunc
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go goroutineFunc(wg)

	wg.Wait()
}
