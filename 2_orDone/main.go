package main

func orDone(done <-chan interface{}, myChan <-chan interface{}) <-chan interface{} {
	valueStream := make(chan interface{})

	go func() {
		defer close(valueStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-myChan:
				if ok == false {
					return
				}
				select {
				case valueStream <- myChan:
				case <-done:
				}
			}
		}
	}()
	return valueStream
}

func main() {
	done := make(chan interface{})
	defer close(done)

	v := make(chan interface{})
	defer close(v)
	for i := range orDone(done, v) {
		// Do Something
	}
}
