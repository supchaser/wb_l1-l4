package main

import (
	"fmt"
	"sync"
	"time"
)

const T1 = 500 * time.Millisecond
const T2 = 2 * time.Second

func goroutineFunc(wg *sync.WaitGroup, stopChan chan struct{}) {
	defer wg.Done()
	for {
		select {
		// если в этот канал никто не пишет, то мы не проваливаемся в этот кейс
		// как только в main мы передадим в канал значение то отработает этот случай и наша горутина завершится
		case <-stopChan:
			fmt.Println("canceled by stop channel")
			return
		// имитируем работу по дефолту
		default:
			time.Sleep(T1)
		}
	}
}

func main() {
	stopChan := make(chan struct{})
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go goroutineFunc(wg, stopChan)

	time.Sleep(T2) // вынужденная мера
	stopChan <- struct{}{}
	wg.Wait()
}
