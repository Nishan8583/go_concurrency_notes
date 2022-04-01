# go_concurrency_notes
Contains notes i made for myself while reading Oreillys Concurrency in Go
NOTE: topic i found interesting and helpful and which required more explaination has a seperate sections on them.

## Common issues in concurrency:
1. Race conditions: 2 or more threads trying to access data in the same time.
2. Atomicity: In a given context certain stuff is indivisible. Ex; get i, incerase i, store i -> within same function can be said to be atomic, within different goroutine but local variable still atmoic, but  within different goroutine but shared variable is not atmoic.
3. Memory Access Synchronization: 2 different threads accessing memory in non atmoic way, they are not accessing in the order they are supposed to.
4. Deadlock: all threads are paused, usually waiting for one another.
5. Livelock: All threads running, but state of program is not moving forward.
6. Starvation: Single thread consuming huge resource.

Concurrency is property of code, parallelism is property of running code.

Communicating Sequential Process (CSP): Basically, I/O is sequantial. An operator gives input to process and another operator gets output. Inspiraiton for go channels.

## Questions to ask for choosing concurrency primitives:
1. is it performance critical? (you have already profiled an realized that channels were bottle neck) use primitives.
2. If transferring ownership use channels.
3. If trying to guard internat state of struct use primitives.

Goroutine shares same address space as its executor, so hass access to its vars.

go keyword schedules the goroutine, no guarantee it runs immediately.

M:N model, M number if goroutines mapped to N os thread.

## Primitives:
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

## tee channel
-> take input from one channel and send them to 2 seperate channel (like 2 seperate area of code base)

## Bride channel
-> take data from multiple channels, like slice of channels, and put it in a single channel

## Queuing
-> using queue
-> Use it as last method of optimization, cause using it before hand may hide deadlock bugs

## context package
-> done pattern is ok, but its difficutl to send message and timeout stuff
-> context helps with this
```go
ctx,cancel := context.WithCancel(context.Background())
defer cancel()

select {
    case <- ctx.Done():
        fmt.Println("stuff",cts.Err())
    // else do stuff here
}

// to cancel goroutine after a certain time reaches, time.Now() used here just as an example, makes no sense obviously
ctx,cancel := context.WithDeadline(context.Background(),time.Now())
defer cancel()

// cancel after 5 second
ctx,cancel := context.WithTimeout(context.Background(),5 *time.Second())
defer cancel()
```
-> Do not pass reference to context, even though context does not seem to change in the code, it changes in the background. having context of N call stack up is bad.

-> when we build on top of parents context, parent is not affected
```go
func child(ctx, context.Context) {
    ctx,cancel := context.WithTimeout(ctx,1 * time.Second)
    // parents ctx is not affected
}
```

## Error propagation:
1. Error should contain
    -> what happened: raw error, usually implicitly generated
    -> when and where it occured: stack trace, timestamp in UTC, context it was running on, 
    -> A friendly user message
    -> How user can get more info: log ID may be ?

Wrap error when passing between boundareis(packages)

2. Timeout can cancellation
-> prevent deadlocks
-> remember to rollback changes, maybe implement in memory first?