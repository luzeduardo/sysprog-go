package main

import (
	"fmt"
	"sync"
)

func main() {
	mainC1ChannelWaitGroup()
	// mainC1Channel()
	// mainC1P43Mutexes()
}

func mainC1ChannelWaitGroup() {
	balls := make(chan string)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		throwBalls("red", balls)
	}()

	go func() {
		defer wg.Done()
		throwBalls("blue", balls)
	}()

	go func() {
		wg.Wait()
		// close the channel after the goroutines are done
		close(balls)
	}()

	for color := range balls {
		fmt.Printf("%s ball received\n", color)
	}
}

func mainC1Channel() {
	balls := make(chan string)
	go throwBalls("red", balls)
	fmt.Println(<-balls, "received")
}

func throwBalls(color string, balls chan string) {
	fmt.Printf("throwing the %s ball\n", color)
	balls <- color
}

func mainC1P43Mutexes() {
	m := sync.Mutex{}
	fmt.Println("Total Items Packed: ", PackItems(&m, 0))
}

func PackItems(m *sync.Mutex, totalItems int) int {
	const workers = 2
	const itemsPerWorker = 1000

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < itemsPerWorker; j++ {
				m.Lock()
				itemsPacked := totalItems
				itemsPacked++
				totalItems = itemsPacked
				m.Unlock()
			}
		}(i)
	}

	wg.Wait()
	return totalItems
}
