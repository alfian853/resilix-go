package util

import "sync"

func AsyncWgRunner(fun func(), wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		fun()
	}()
}
