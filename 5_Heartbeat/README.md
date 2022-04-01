# Heartbeat
-> A signalt that a concurrent process gives to outside parties
-> 2 types: A time interval, and while starting a unit of work
-> example for time interval in main.go
-> sending in a unit of work helpful in testing involving multiple goroutines, because no need to rely on timeouts
