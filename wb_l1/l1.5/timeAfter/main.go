package main

import (
	"fmt"
	"time"
)

const N = 5   // время, по истечению которого программа завершится
const M = 100 // кол-во значений, которое будет сгенерировано
const T = 100 // константа, на которую домнажается time.Millisecond для имитации работы горутины

// Решение аналогично решению с контекстом, но в качестве таймаута используетя time.After

func generator() <-chan int {
	out := make(chan int)

	go func() {
		for i := range M {
			out <- i * i
			time.Sleep(T * time.Millisecond)
		}

		close(out)
	}()

	return out
}

func main() {
	// создаем timer
	// time.After возвращает инициализированный каналs
	timer := time.After(N * time.Second)
	out := generator()

	for {
		select {
		// здесь мы ждём N секунд по прихода писателя
		// через N секунд канал отправит текущее время и канал закроется, поэтому мы сможем провалиться в этот кейс
		case <-timer:
			fmt.Println("canceled by timer")
			return
		// читаем из канала
		case v, ok := <-out:
			// если канал закрыт, ывыходим
			if !ok {
				fmt.Println("canceled by closed channel")
				return
			}

			// пишем прочитанное из канала значение в stdout
			fmt.Println(v)
		}
	}
}
