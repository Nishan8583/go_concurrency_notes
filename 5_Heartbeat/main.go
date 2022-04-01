package main

import (
	"fmt"
	"time"
)

// doWork is the concurrent process that will be giving heartbeat signal
func doWork(
	done <-chan interface{},
	pulseInterval time.Duration,
) (<-chan interface{}, <-chan interface{}) {

	result := make(chan interface{}) // holds the return value
	pulse := make(chan interface{})  // used for giving pulse signal

	go func() {
		defer close(result)
		defer close(pulse)

		pulseInterval := time.Tick(pulseInterval)
		workDone := time.After(10 * time.Second)

		sendPulse := func() {
			select {
			case <-done:
				fmt.Println("done called inside pulse sending data")
			case pulse <- 1:
			default: // default case is needed because noone might be listening to pulse
			}
		}

		// need this infinite loop, because once the pulse is sent, select statement will exit, but we want
		// the work to still be going on
		for {
			select {
			case <-done:
				fmt.Println("done was closed")
			case <-pulseInterval:
				sendPulse()
			case <-workDone:
				result <- 1
				fmt.Println("work has been completed")
			}
		}
	}()

	return pulse, result

}

func main() {
	done := make(chan interface{})
	defer close(done)

	heartbeet, result := doWork(done, 1*time.Second)
	timer := time.After(30 * time.Second)
	for {
		select {
		// need to check if either has been closed, it indicates that the goroutine has returned, with or without error
		// else condition will be handled by timer chan
		case _, ok := <-heartbeet:
			if !ok {
				fmt.Println("heartbeat closed")
			}
			fmt.Println("pulse from goroutine")
		case _, ok := <-result:
			if !ok {
				fmt.Println("result cancelleed")
			}
			fmt.Println("work completed")
			return
		case <-timer:
			// since we handeled both channel closed condition in above cases
			// we reach here when channel has not been closed in a given time
			// the gorutine is not working properly or timeout exdceeded or panicked?
			fmt.Println("time exceeded, goroutine unhealthy")
			return
		}
	}
}
