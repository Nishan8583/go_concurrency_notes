# Fan out Fan in pattern
1. Fan out: starting multiple goroutine to handle input from a single pipeline
2. Fan in: combining multiple results into one
3. Use it when: your task is order independent. task takes long time to run

## In practical while writing code
(taken from the book, there should be easier way)
1. Use generator pattern for a heavy task, so it returns channel, another goroutine will be sending data in that channel.
2. Call that generater multiple times by spawwning goroutines, and store that in something like a slice of goroutine.
3. Pass that collection of goroutine into a fan-in funciton, it does the following task:
    -> creates a multiplex channel
    -> spawns multiple goroutine all of which tries to get from each of channel from collection and tries to put it in the multiplex channel.
    -> a goroutine to wait for each channel to fnixish.
    -> returns multiplex channel
4. Read from the multiplex channel from main.
see main.go