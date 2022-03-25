# Pipelines
-> Series of things(functions or others) that takes input, performs actions on it, and gives output
-> each operation is a stage
```go
// Consider functions multipy and add already declared
func multiply(values []int, n int) []int{}
func add(values []int, n int) []int{}

// now we can do something like
test := []int{1,2,3}
for _,v := range add(multiply(test,5),4) { .. }

// Here funcitons add and multiply are stages, we combine them in for loop to create a pipeline
```
-> Using channels makes them even better, look at main.go for detailed code.