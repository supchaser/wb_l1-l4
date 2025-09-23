package main

import (
	"fmt"
	"time"
)

const N = 5   // время, по истечению которого программа завершится
const M = 100 // кол-во значений, которое будет сгенерировано
const T = 100 // константа, на которую домнажается time.Millisecond для имитации работы горутины

func generator() <-chan int {
	out := make(chan int)

	go func() {
		for i := range M {
			time.Sleep(T * time.Millisecond)
			out <- i + i
		}

		close(out)
	}()

	return out
}

func main() {
	// вариант с созданием структуры timer
	// time.After - это обёртка над time.NewTimer, поэтому можно сделать еще и такой вариант
	timer := time.NewTimer(N * time.Second)

	// также можно самому создать таймер с помощью канала
	ch := make(chan int)
	go func() {
		time.Sleep(N * 2 * time.Second)
		close(ch)
	}()

	out := generator()

	for {
		select {
		case <-timer.C:
			fmt.Println("canceled by timer")
			return
		case <-ch:
			fmt.Println("canceled by channel")
			return
		case v, ok := <-out:
			if !ok {
				fmt.Println("canceled by closed channel")
				return
			}
			fmt.Println(v)
		}
	}
}
