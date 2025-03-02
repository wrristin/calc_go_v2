package main

import (
	"calc_service/internal/agent"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < agent.ComputingPower; i++ {
		wg.Add(1)
		go agent.Worker(&wg)
	}
	wg.Wait()
}
