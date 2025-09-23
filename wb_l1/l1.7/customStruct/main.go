package main

import (
	"fmt"
	"sync"
	"time"
)

// стр-ра ConcurrentMap
// содержит в себе поля m - мапа, просто хэш таблица
// mu - RWMutex для того, чтобы на чтение не было блокировки и несколько горутин могли бы читать данные из мапы
// также данная реализация стр-ра сделана с помощью Дженериков, чтобы использование мапы было более гибким
type ConcurrentMap[K comparable, V any] struct {
	m  map[K]V
	mu sync.RWMutex
}

// конструктор
func CreateConcurrentMap[K comparable, V any]() *ConcurrentMap[K, V] {
	return &ConcurrentMap[K, V]{
		m: make(map[K]V),
	}
}

// метод сохранения ключа и значения в мапу
func (cm *ConcurrentMap[K, V]) Store(key K, value V) {
	// для данного метода мы не можем применить RLock и RUnlock, т.к. появляется вероятность возникновения data race
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.m[key] = value
}

// метод чтения значения по ключу из мапы
func (cm *ConcurrentMap[K, V]) Load(key K) (value V, ok bool) {
	// тут мы уже можем применить RLock и RUnlock, потому что для нас не будет проблематичным, если
	// несколько горутин будут одновременно читать данные
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	value, ok = cm.m[key]
	return value, ok
}

// метод, который позволяет сохранить данные в мапу, если такого ключа еще нет
// или, если ключ уже существует, вернуть текущее значение
func (cm *ConcurrentMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	actual, loaded = cm.m[key]
	if !loaded {
		cm.m[key] = value
		actual = value
	}

	return actual, loaded
}

// метод удаления ключа из мапы
func (cm *ConcurrentMap[K, V]) Delete(key K) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.m, key)
}

// метод получения длины мапы
func (cm *ConcurrentMap[K, V]) Len() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return len(cm.m)
}

// таким образом можно сделать вывод, что если производимая операция допустима для конкурентного выполнения (например, чтение значения по ключу или получения длины мапы), то
// мы можем использовать RLock и RUnlock, для остальных операций мы должно жестко блокироваться

// также можно было бы еще использовать sync.Map

func main() {
	// Создаем мапу
	cm := CreateConcurrentMap[string, int]()
	wg := &sync.WaitGroup{}

	// сохранили 100 значений
	for i := range 100 {
		wg.Go(func() {
			key := fmt.Sprintf("key%d", i%10)
			cm.Store(key, i)
			time.Sleep(time.Microsecond * 10)
		})
	}

	// прочитали 50 значений
	for i := range 50 {
		wg.Go(func() {
			key := fmt.Sprintf("key%d", i%10)
			value, ok := cm.Load(key)
			// необходимо использовать полученное значение, что компилятор не оптимизировал вызов
			if ok {
				_ = value
			}
		})
	}

	// конкурентный LoadOrStore
	for i := range 30 {
		wg.Go(func() {
			key := fmt.Sprintf("shared_key%d", i%5)
			actual, loaded := cm.LoadOrStore(key, i*100)
			_ = actual
			_ = loaded
		})
	}

	// удаление во время чтения записи
	for i := range 20 {
		wg.Go(func() {
			key := fmt.Sprintf("key%d", i%10)
			if i%3 == 0 {
				cm.Delete(key)
			} else {
				cm.Load(key)
			}
		})
	}

	// конкуретное получение длины
	for range 10 {
		wg.Go(func() {
			length := cm.Len()
			_ = length // используем значение
		})
	}

	wg.Wait()

	fmt.Printf("финальный размер мапы: %d\n", cm.Len())

	for i := range 5 {
		key := fmt.Sprintf("key%d", i)
		if value, ok := cm.Load(key); ok {
			fmt.Printf("%s: %d\n", key, value)
		} else {
			fmt.Printf("%s: не найден\n", key)
		}
	}
}
