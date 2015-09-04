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

/* The purpose of main is to use manual memory management to improve on the
   garbage_creator program. Buffers are recycled which allows for much less
   waste and a more efficient use of memory. */
func main() {
	// Create a pool of 20 buffers
	pool := make([][]byte, 20)

	// Create a channel which can hold 5 buffers
	buffer := make(chan []byte, 5)

	var m runtime.MemStats
	makes := 0

	for {
		var b []byte
		select {
		// If there is a buffer in the channel, place the slice in b
		case b = <-buffer:
			// If the channel is empty, create a new buffer
		default:
			makes += 1
			b = makeBuffer()
		}

		// Pick a random position in the pool
		i := rand.Intn(len(pool))

		if pool[i] != nil {
			select {
			// If the channel isn't full, send the buffer at pool[i]
			case buffer <- pool[i]:
				pool[i] = nil
				// If it is full, don't send a buffer
			default:
			}
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
		fmt.Printf("%d,%d,%d,%d,%d,%d\n", m.HeapSys, bytes, m.HeapAlloc,
			m.HeapIdle, m.HeapReleased, makes)

	}
}
