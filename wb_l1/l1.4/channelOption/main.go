package main

import (
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// функция воркер, нужна для имитации работы
func worker(stopCh <-chan struct{}, wg *sync.WaitGroup) {
	// уменьшаем счётчик вэйт группы (-1)
	defer wg.Done()

	for {
		select {
		// если в main канал stopCh закроется, то мы провалимся в этот case и завершим наш worker
		case <-stopCh:
			return
		// по дефолту спим рандомное время (макс 5 секунд)
		default:
			time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
		}
	}
}

func main() {
	// канал, из которого будут читать горутины для корректного заврешения
	stopCh := make(chan struct{})
	// вейт группа, необходима для того, чтобы дождаться выполнения всех горутин
	wg := &sync.WaitGroup{}
	// создаем 5 воркеров
	for range 5 {
		// увеличиваем счетчик воркеров на 1
		wg.Add(1)
		go worker(stopCh, wg)
	}

	// создаем буф канал для получения системных сигналов
	sigCh := make(chan os.Signal, 1)
	// подписываемся каналом sigCh на сигнала SIGINT (CTRL + C) и SIGTERM
	// эти сигналы будут отправлены в канал sigCh вместо стандартных обработчиков гошного рантайма
	// но может возникнуть вопрос: подписываемся на сигнал SIGINT, а в функцию мы передали os.Interrupt, почему?
	// это сделано для обеспечения кроссплотформенности, чтобы и на Windows и на Unix-подобных системах os.Interrupt преобразовывался в конкретный сигнал
	// под копотом в Go os Interrupt это: var Interrupt = sysyscall.SIGINT (пакет os)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// чтение из канала - блокирующая операция, поэтому на этом этапе блокируем наш main до момента
	// пока не получим сигнал завершения
	<-sigCh
	// после разблокировки, закрываем канал stopCh и тем самым инициируем завершение всех воркеров
	close(stopCh)

	// дожидаемся выполнения всех горутин
	wg.Wait()

}
