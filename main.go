package main

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/remeh/sizedwaitgroup"
)

var MaxSteps uint64 = 0
var MaxLock sync.Mutex

func checkMaxSteps(seed uint64, i uint64, steps uint64) {
	MaxLock.Lock()
	if steps > MaxSteps {
		MaxSteps = steps
		fmt.Printf("NEW RECORD LENGTH:\nSEED: %v, I: %v, STEPS: %v\n", seed, i, steps)
	}
	MaxLock.Unlock()
}

func collatz(seed uint64, i uint64, steps uint64) {

	if seed%1000000 == 0 {
		checkMaxSteps(seed, i, steps)
	} else if i == 1 {
		checkMaxSteps(seed, i, steps)
		//fmt.Printf("seed:%v steps:%v | ", seed, steps)
	} else if i%2 == 0 {
		i = i / 2
		collatz(seed, i, steps+1)
	} else {
		i = i*3 + 1
		collatz(seed, i, steps+1)
	}
}

func main() {
	numCPU := runtime.NumCPU()
	fmt.Printf("Found %v vCPUs.\n", numCPU)
	swg := sizedwaitgroup.New(numCPU)
	fmt.Println("Running: ")

	var max uint64 = 18446744073709551615
	for n := max; n > 0; n-- {
		swg.Add()
		go func(n uint64) {
			defer swg.Done()
			collatz(n, n, 0)
		}(n)
	}
	swg.Wait()
}
