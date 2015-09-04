package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

// Return a buffer between 5mil and 10mil bytes
func makeBuffer() []byte {
	return make([]byte, rand.Intn(5000000)+5000000)
}

/* The purpose of main is to replicate a common occurrence:
   Memory is allocated over time, eventually some of it is
   no longer needed. This can easily happen in networked Go
   programs causing a lot of wasted memory. */
func main() {
	// Create a pool of 20 buffers
	pool := make([][]byte, 20)

	var m runtime.MemStats
	makes := 0

	// Every second create a buffer and add it to the pool
	for {
		b := makeBuffer()
		makes += 1
		// Pick a random position in the pool
		i := rand.Intn(len(pool))
		/* If a buffer is already in pool[i] it is replaced
		   by b and becomes garbage */
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
		fmt.Printf("%d,%d,%d,%d,%d,%d\n", m.HeapSys, bytes, m.HeapAlloc,
			m.HeapIdle, m.HeapReleased, makes)

	}
}
