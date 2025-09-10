package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

// Функция workerPoll принимает кол-во воркеров, и канал на чтение.
// В цикле создаем необходимое кол-во горутин
func workerPool(numWorkers int, ch <-chan string) {
	// Создаём Wait Group'у для того чтобы синхронизировать работу воркеров
	wg := &sync.WaitGroup{}
	// Добавляем в счетчик кол-во горутин, которых необходимо будет дождаться
	wg.Add(numWorkers)
	for i := range numWorkers {
		go worker(i, ch, wg)
	}

	wg.Wait()
}

func worker(id int, ch <-chan string, wg *sync.WaitGroup) {
	// После завершения работы воркера уменьшаем счетчик
	defer wg.Done()
	for str := range ch {
		fmt.Printf("Worker №%d, data: %s \n", id, str)
	}
}

func main() {
	// Проверка аргументов командной строки
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <number_of_workers>")
		return
	}

	// получаем кол-во воркеров из cmd
	numWorkers, err := strconv.Atoi(os.Args[1])
	if err != nil || numWorkers <= 0 {
		fmt.Println("Please provide a valid positive number of workers")
		return
	}

	ch := make(chan string)

	// в отдельной горутине запускаем воркер пул
	go workerPool(numWorkers, ch)

	// создаем канал пустых стр-ур для Graceful Shutdown
	stop := make(chan struct{})
	// создаем канал, который будет сигнализировать о прерывании работы программы
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// в отдельной горутине читаем из канал sigCh,
	// т.к. чтение из канал блокирующая операция, то эта горутина будет висеть,
	// пока нам не удастся прочитать из канала
	go func() {
		<-sigCh
		// после того как мы прочитаем из канала закрываем канал stop
		close(stop)
	}()

	i := 0
	for {
		// блокирующая конструкция select
		// когда мы закроем канал stop, то мы сможем прочитать из него nil-value
		// после этого мы закромем наш канал ch и завершим программу
		select {
		case <-stop:
			close(ch)
			return
		// по дефолту пишем в канал ch рандомные строки
		default:
			ch <- fmt.Sprintf("%s", strconv.Itoa(i)) + rand.Text()
			i++
		}
	}
}
