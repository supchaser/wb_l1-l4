package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const T = 3 * time.Second

func goroutineFunc(wg *sync.WaitGroup) {
	defer wg.Done()
	defer fmt.Println("runtime.Goexit начал завершение горутины, но перед этим он вызовет все defer'ы")

	// имитация работы горутиныы
	time.Sleep(T)

	// runtime.Goexit - завершает только ту горутину, которая его вызывает
	runtime.Goexit()
}

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go goroutineFunc(wg)

	wg.Wait()
}
