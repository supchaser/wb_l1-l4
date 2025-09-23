package main

import (
	"context"
	"fmt"
	"time"
)

const N = 5   // время, по истечению которого программа завершится
const M = 100 // кол-во значений, которое будет сгенерировано
const T = 100 // константа, на которую домнажается time.Millisecond для имитации работы горутины

// функция генератор
// принимает контекст, возвращает канал на чтение
func generator(ctx context.Context) <-chan int {
	// создаём канал
	out := make(chan int)

	// все блокирующие операции должны выполняться в горутине
	go func() {
		// закрываем канал out отложенным вызовом
		defer close(out)

		for i := range M {
			select {
			// по истечению N секунд контекст отменится
			// ctx.Done() вернёт закрытый канал
			// чтение из закрытого канал - nil value
			// следовательно мы провалимся в этот кейс и завершим нашу горутину
			case <-ctx.Done():
				fmt.Println("canceled by context")
				return
			// пишем в канал out
			case out <- i:
				// имитация работы
				time.Sleep(100 * time.Millisecond)
			}

		}
	}()

	return out
}

func main() {
	// создаём контекст с таймаутом в N секунд
	// это означает, что через N секунд контекст отменится и программа начнёт своё завершение
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*N)
	defer cancel()

	// получаем канал out из функции generator
	out := generator(ctx)

	// запускаем бесконечный цикл, но можно было сделать цикл на M итераций
	for {
		select {
		// если контекст отменился, то завершаем нашу программу
		case <-ctx.Done():
			fmt.Println("canceled by context")
			return
		// пытаемся читать из канала
		case v, ok := <-out:
			// если канал закрыт, то завершаем нашу программу
			if !ok {
				fmt.Println("canceled by closed channel")
				return
			}
			select {
			// если контекст отменился, то завершаем нашу программу
			case <-ctx.Done():
				fmt.Println("canceled by context")
				return
			default:
				// по дефолту пишем наше значение, прочитанное из канала в stdout
				fmt.Println(v)
			}
		}
	}
}
