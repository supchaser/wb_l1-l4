package main

import (
	"fmt"
	"sync"
)

// Функция doubler принимает массив чисел и wait group'у,
// на каждой итерации цикла создается горутина,
// которая конкурирует с другими созданными за запись в stdout.
// Wait Group'а здесь необходима для того, чтобы синхронизировать выполнение всех горутин
// и наша программа не завершится, пока не дождется выполнения всех горутин.
func doubler(nums []int, wg *sync.WaitGroup) {
	for _, num := range nums {
		wg.Go(func() {
			fmt.Println(num * num)
		})
	}

}

func main() {
	nums := []int{2, 4, 6, 8, 10}
	wg := &sync.WaitGroup{}
	doubler(nums, wg)

	// Дожидаемся выполнения всех горутин
	wg.Wait()
}
