package main

import (
	"container/list"
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

var makes int
var frees int

// Return a buffer between 5mil and 10mil bytes
func makeBuffer() []byte {
	makes += 1
	return make([]byte, rand.Intn(5000000)+5000000)
}

// A new type to wrap the buffers in
type queued struct {
	when  time.Time // When the buffer was created
	slice []byte    // The slice itself
}

/* Make the recycler. Any buffer in the pool (queue) that is older than
   a minute is considered unlikely to be used again and can be
   discarded */
func makeRecycler() (get, give chan []byte) {
	get = make(chan []byte)  // Channel to get buffers from the pool
	give = make(chan []byte) // Channel to give buffers back to the pool

	// Start function as a goroutine
	go func() {
		q := new(list.List)
		for {
			if q.Len() == 0 {
				// If the queue is empty, create a new buffer and add it to the queue
				q.PushFront(queued{when: time.Now(), slice: makeBuffer()})
			}

			e := q.Front()

			// Start a 1 minute long timer
			timeout := time.NewTimer(time.Minute)
			select {
			// If there is a buffer in give, add it to the queue
			case b := <-give:
				timeout.Stop()
				q.PushFront(queued{when: time.Now(), slice: b})

			// If there is a buffer in queue, send it to the get channel
			case get <- e.Value.(queued).slice:
				timeout.Stop()
				q.Remove(e)

			// Time is up! Any buffer older than a minute is discarded.
			case <-timeout.C:
				// Go through list checking if creation time is more than 1 minute ago
				e := q.Front()
				for e != nil {
					n := e.Next()
					if time.Since(e.Value.(queued).when) > time.Minute {
						q.Remove(e)
						e.Value = nil
					}
					e = n
				}
			}
		}
	}() // Start function as a goroutine with no params

	return // Return the get and give channels (named return values)
}

func main() {
	// Create a pool that can hold 20 buffers
	pool := make([][]byte, 20)

	// Create a recycler in the form of 2 channels
	get, give := makeRecycler()

	var m runtime.MemStats
	for {
		// Get a buffer from the pool
		b := <-get
		i := rand.Intn(len(pool))
		// Give a buffer back to the pool
		if pool[i] != nil {
			give <- pool[i]
		}

		pool[i] = b

		time.Sleep(time.Second)

		// Get the total amount of bytes in the pool
		bytes := 0
		for i := 0; i < len(pool); i++ {
			if pool[i] != nil {
				bytes += len(pool[i])
			}
		}

		/* Report the current memory statistics
		   HeapSys: # of bytes program has asked OS for
		   bytes: current total size of the buffers in the pool
		   HeapAlloc: # of bytes currently allocated on the heap
		   HeapIdle: # of unused bytes in the heap
		   HeapReleased: # of bytes returned to the OS
		   makes: total # of buffers created
		*/
		runtime.ReadMemStats(&m)
		fmt.Printf("%d,%d,%d,%d,%d,%d,%d\n", m.HeapSys, bytes, m.HeapAlloc,
			m.HeapIdle, m.HeapReleased, makes, frees)
	}
}
