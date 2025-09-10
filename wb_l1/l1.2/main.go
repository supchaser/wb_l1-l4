package main

import (
	"fmt"
	"sync"
)

func doubler(nums []int, wg *sync.WaitGroup) {
	for _, val := range nums {
		wg.Go(func() {
			fmt.Println(val * val)
		})
	}

}

func main() {
	nums := []int{2, 4, 6, 8, 10}
	wg := &sync.WaitGroup{}
	doubler(nums, wg)
	wg.Wait()
}
