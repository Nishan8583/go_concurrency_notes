# go_concurrency_notes
Contains notes i made for myself while reading Oreillys Concurrency in Go

Common issues in concurrency:
1. Race conditions: 2 or more threads trying to access data in the same time.
2. Atomicity: In a given context certain stuff is indivisible. Ex; get i, incerase i, store i -> within same function can be said to be atomic, within different goroutine but local variable still atmoic, but  within different goroutine but shared variable is not atmoic.
3. Memory Access Synchronization: 2 different threads accessing memory in non atmoic way, they are not accessing in the order they are supposed to.
4. Deadlock: all threads are paused, usually waiting for one another.
5. Livelock: All threads running, but state of program is not moving forward.
6. Starvation: Single thread consuming huge resource.

Concurrency is property of code, parallelism is property of running code.

Communicating Sequential Process (CSP): Basically, I/O is sequantial. An operator gives input to process and another operator gets output. Inspiraiton for go channels.

Questions to ask for choosing concurrency primitives:
1. is it performance critical? (you have already profiled an realized that channels were bottle neck) use primitives.
2. If transferring ownership use channels.
3. If trying to guard internat state of struct use primitives.

Goroutine shares same address space as its executor, so hass access to its vars.

go keyword schedules the goroutine, no guarantee it runs immediately.

M:N model, M number if goroutines mapped to N os thread.

Primitives:
1. waitgroup.
2. Mutex: use RWMutex for performance boost. 
3. Cond: Similar to channel, faster, but difficult to read IMO.
3. sync pool: A seperate section on it.
4. Channels:
    -> Even bidirectional channel, when function signature takes uni directional becomes unidirectional inside the funciton,
    -> u can range over channels, when channel closes, loop exists
    ->i,ok := <-ch returns false if channels closed
    -> Owner should :
        -> instantiate channel
        -> perform ownership transfer
        -> close channel
        -> encapsulate all of this and return reader only


## Preventing goroutine leaks
-> a groutine exits when
    1. its done its job
    2. panic
    3. when its told
-> 1 and 2 are pretty clear, but it might not always happen, for that we need 3.
-> use of done channel
```go
func(done <-chan int) {
    for {
        select {
            case <-done:
                return
            case ...
        }
    }
}

// the caller of the above goroutine funciton is responsible for closing the channel
// or time.AfterFunc(DurationOfTime, close(done))
```