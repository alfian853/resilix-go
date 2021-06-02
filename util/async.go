package util

import "sync"

func AsyncWgRunner(fun func(), wg *sync.WaitGroup){
	go asyncWgRunner(fun, wg)
}

func asyncWgRunner(fun func(), wg *sync.WaitGroup){
	defer wg.Done()
	fun()
}

