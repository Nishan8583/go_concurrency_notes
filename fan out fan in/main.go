package main

import (
	"fmt"
	"math/big"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// fanIn receives the usual done channel, and the slice of channels that have been fanned out
// it then loops through each of those fanned out channels and puts that channel in a multiplex stream
// NOTE that there will be seperate goroutine for each channel responsible for putting that channel in
// that particular stream, it waits for each channel to be put in that multiplex stream on a seperate goroutine
// mean while it returns the multiplex stream
func fanIn(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} {

	var wg sync.WaitGroup // make sure all channels have been drained

	multiplexedStream := make(chan interface{})

	// puts the channels in the multiplexed stream
	multiplex := func(c <-chan interface{}) {
		defer wg.Done()

		for i := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- i:
			}
		}
	}

	// wait till the number of channels
	wg.Add(len(channels))
	for _, c := range channels {
		// for each channel, a seperate goroutine to put in multiplex channel
		go multiplex(c)
	}

	// a groutine to wait for all channels to be put in the multiplex channel,
	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

func reapeatFn(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	valueStream := make(chan interface{})

	go func() {
		defer close(valueStream)
		for {
			select {
			case <-done:
				return
			case valueStream <- fn():
			}
		}
	}()
	return valueStream
}

func toInt(done <-chan interface{}, valueStream <-chan interface{}) <-chan int {
	intStream := make(chan int)

	go func() {
		defer close(intStream)
		for v := range valueStream {
			select {
			case <-done:
				return
			case intStream <- v.(int):
			}
		}
	}()

	return intStream
}

// take takes num number of values off the valueStream channel
func take(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{})

	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()

	return takeStream
}

func primerFinder(done <-chan interface{}, inStream <-chan int) <-chan interface{} {
	primerNumbers := make(chan interface{})

	go func() {
		defer close(primerNumbers)
		for i := range inStream {
			select {
			case <-done:
				return
			default:
				v := int64(i)
				if bv := big.NewInt(v).ProbablyPrime(0); bv {
					primerNumbers <- v
				}
			}
		}
	}()
	return primerNumbers
}

func main() {

	// done channel
	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	// function that generates random number
	rand := func() interface{} { return rand.Intn(500000) }

	// generate random number and put it in the strea,
	randIntStream := toInt(done, reapeatFn(done, rand))

	// getting the number of goroutines that we would run
	numFinders := runtime.NumCPU()
	finders := make([]<-chan interface{}, numFinders)
	fmt.Println("spawnning ", numFinders)

	// as u can see spawwning mutiple goroutines and putting the returned channel in a slice of channel
	for i := 0; i < numFinders; i++ {
		finders[i] = primerFinder(done, randIntStream)
	}

	// fanIn loops through all the finders channel, and spawns seperate goroutine to put
	// those channels in a single multiplexed channel, a seperate goroutines waits for those channels to complete
	// while it returns the multiplexed channel
	for prime := range take(done, fanIn(done, finders...), 10) {
		fmt.Println(prime)
	}

	fmt.Println("since start", time.Since(start))
}
