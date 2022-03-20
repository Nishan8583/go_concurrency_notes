# Sync Pool
1. Set of temporary objects saved and retreived temporarily. 
2. Purpose, allocated item, not used now but maybe used later, cache it
3. When multiple goroutines using objects
4. If only pool has reference to it, GC might collect it
5. When you GET. make no assumption of its state.
6.  potentially reduce the GC workload and improve performance, by reducing number of allocation and deallocation
7. There are two pools behind the scene, local pool (active) and victim pool (archived). When GC runs, object inside victim pool is removed and object inside local pool will be moved to victim pool.
8. Put puts in local pool
9. Get tries to get from Victim cache first, if victim cache was emppty, it gets from local cahce.
10. Use case pool of db connections

Referenec: 
(info about victim and local cache) https://medium.com/swlh/go-the-idea-behind-sync-pool-32da5089df72


```go
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
    s.Write([]bytes("append"))
    // At this point, if GC ran, this buffer *might* exist already, in
    // which case it will contain the bytes of the string "dirtyappend"
    fmt.Println(s)
    // So use pools wisely, and clean up after yourself
    s.Reset()
    pool.Put(s)

    // When you clean up, your buffer should be empty
    s = pool.Get().(*bytes.Buffer)
    // Defer your Puts to make sure you don't leak!
    defer pool.Put(s)
    s.Write([]byte("reset!"))
    // This prints "reset!", and not "dirtyappendreset!"
    fmt.Println(s)
}
```