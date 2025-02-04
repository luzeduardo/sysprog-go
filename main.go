package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	mainC1P35SignalingChannel()
	// mainC1P30BufferedChannels()
	// mainC1ChannelWaitGroup()
	// mainC1Channel()
	// mainC1P22Mutexes()
}

func mainC1P35SignalingChannel() {
	var wg sync.WaitGroup
	signalChannel := make(chan bool)

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Goroutine 1: Waiting for signal")
		<-signalChannel
		fmt.Println("Goroutine 1: Received signal")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Goroutine 2 is about to send a signal")
		signalChannel <- true
		fmt.Println("Goroutine 2: Sent signal")
	}()

	wg.Wait()
	fmt.Println("All goroutines have finished their jobs")
}

func mainC1P30BufferedChannels() {
	clownChannel := make(chan int, 3)
	clowns := 5

	go func() {
		defer close(clownChannel)
		for clownID := range clownChannel {
			ballon := fmt.Sprintf("Ballon %d", clownID)
			fmt.Printf("Driver: Drove the car with %s inside \n", ballon)

			time.Sleep(time.Millisecond * 500)
			fmt.Printf("Driver: Clown finished with %s, the car is ready for more!\n", ballon)
		}
	}()

	var wg sync.WaitGroup

	for clown := 1; clown <= clowns; clown++ {
		wg.Add(1)

		go func(clownID int) {
			defer wg.Done()
			ballon := fmt.Sprintf("Ballon %d", clownID)
			fmt.Printf("Clown %d: Hopped into the car with %s\n", clownID, ballon)
			select {
			case clownChannel <- clownID:
				fmt.Printf("Clown %d: Finished with %s\n", clownID, ballon)
			default:
				fmt.Printf("Clown %d: Ops the car is full, can't fit %s!\n", clownID, ballon)
			}
		}(clown)

		fmt.Println("Circus car ride is over")
	}

	wg.Wait()
	fmt.Println("Circus car ride is over")
}

func mainC1ChannelWaitGroup() {
	balls := make(chan string)
	// create a WaitGroup to wait for the goroutines to finish
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
		// so the range loop in the main goroutine can finish
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

func mainC1P22Mutexes() {
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
