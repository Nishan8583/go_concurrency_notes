package main

import "fmt"

// generater takes a slice of integer and starts a seperate goroutine and puts it in a channel
// it then returns the read only channel for it,
// converting discrete values to stream of data i.e. thats what generators do
func generator(done <-chan int, numbers []int) <-chan int {

	// make a channel, we are closing it in then goroutine we create
	// thus ensuring that opening and closing the channel, returning only read capability
	// fulfilling the pattern
	gen := make(chan int)

	go func() {
		defer close(gen)
		for _, v := range numbers {

			// preventing goroutine leaks by using done channel
			select {
			case <-done:
				fmt.Println("generate goroutine exiting")
				return
			case gen <- v:
			}
		}
	}()
	return gen

}

// multiply takes done, and an input stream i.e. read only channel, and a multiplier
func multipy(
	done <-chan int,
	inStream <-chan int,
	n int,
) <-chan int {

	result := make(chan int) // will contain the output of multiplicaiton
	go func() {
		defer close(result)

		// ranging over input stream channel, if its closed or values are finished, it exits
		for i := range inStream {
			select {
			case <-done:
				fmt.Println("closing multiply")
				return
			case result <- n * i:
			}
		}
	}()
	return result
}

// same as multiplier
func add(
	done <-chan int,
	inStream <-chan int,
	n int,
) <-chan int {
	result := make(chan int)

	go func() {
		defer close(result)
		for i := range inStream {
			select {
			case <-done:
				fmt.Println("exiting add goroutine")
			case result <- n + i:
			}
		}
	}()

	return result
}
func main() {
	done := make(chan int)

	numStream := generator(done, []int{1, 2, 3, 4, 5, 6})

	// see how were are passing channels from ones output to anothers input
	for i := range add(done, multipy(done, numStream, 2), 1) {
		fmt.Println(i)
	}

}

// Since done was taken as input, closing one will close others
// Since inputstream are passed around as channels, closing input stream will close others
// using channels in pipelines, makes it concurrent safe
