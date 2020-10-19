package sum

import (
	"sync"
)

//Regular r
func Regular() int64 {
	sum := int64(0)
	for i := 0; i < 20000000; i++ {
		sum++
	}
	return sum
}

//Concurrently c
func Concurrently() int64 {
	wg := sync.WaitGroup{}
	wg.Add(2)
	mu := sync.Mutex{}

	sum := int64(0)

	go func() {
		defer wg.Done()
		val := int64(0)
		for i := 0; i < 10000000; i++ {
			val++
		}
		mu.Lock()
		defer mu.Unlock()
		sum += val
	}()
	go func() {
		defer wg.Done()
		val := int64(0)
		for i := 0; i < 10000000; i++ {
			val++
		}
		mu.Lock()
		defer mu.Unlock()
		sum += val
	}()
	wg.Wait()
	return sum
}
