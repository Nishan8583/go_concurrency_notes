package main

import (
	"bytes"
	"fmt"
	"sync"
)

var pool = sync.Pool{
	// New creates an object when the pool has nothing available to return.
	// New must return an interface{} to make it flexible. You have to cast
	// your type after getting it.
	New: func() interface{} {
		// Pools often contain things like *bytes.Buffer, which are
		// temporary and re-usable.
		return &bytes.Buffer{}
	},
}

func main() {
	// When getting from a Pool, you need to cast
	s := pool.Get().(*bytes.Buffer)
	// We write to the object
	s.Write([]byte("dirty"))
	// Then put it back
	pool.Put(s)

	// Pools can return dirty results

	// Get 'another' buffer
	s = pool.Get().(*bytes.Buffer)
	// Write to it
	s.Write([]byte("append"))
	// At this point, if GC ran, this buffer *might* exist already, in
	// which case it will contain the bytes of the string "dirtyappend"
	fmt.Println(s)
	// So use pools wisely, and clean up after yourself
	//s.Reset()
	pool.Put(s)

	// When you clean up, your buffer should be empty
	s = pool.Get().(*bytes.Buffer)
	s.Reset()
	s.Write([]byte("reset!"))
	// This prints "reset!", and not "dirtyappendreset!"
	fmt.Println(s)

	// If you do not put it back, new bytes.buffer will be created
	pool.Put(s)

	// to put new object in pool, but if you just call Get(), it will create new one
	pool.Put(pool.New())
	s1 := pool.Get().(*bytes.Buffer)
	s2 := pool.Get().(*bytes.Buffer)
	defer pool.Put(s1)
	defer pool.Put(s1)
	fmt.Println("new values", s1.String(), s2.String())
}
