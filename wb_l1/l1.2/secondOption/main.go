package main

import "fmt"

// Функция doublerWithChannels принимает массив чисел, и канал на запись.
// Конкурентно пишем в этот канал
func doublerWithChannels(nums []int, result chan<- int) {
	for _, num := range nums {
		go func() {
			result <- num * num
		}()
	}
}

func main() {
	nums := []int{2, 4, 6, 8, 10}
	results := make(chan int)

	doublerWithChannels(nums, results)

	// читаем из канала
	for range len(nums) {
		fmt.Println(<-results)
	}
}
