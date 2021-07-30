package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/remeh/sizedwaitgroup"
)

var MaxSteps uint64 = 0
var MaxLock sync.Mutex

var N uint64 = 0

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

	go func() {
		for {
			oldn := N
			time.Sleep(15 * time.Second)
			diff := N - oldn
			fmt.Printf("cps: %v, ", diff/15.0)
		}
	}()

	var max uint64 = 18446744073709551615 // 2^64-1
	for N = max; N > 0; N-- {
		swg.Add()
		go func(n uint64) {
			defer swg.Done()
			collatz(N, N, 0)
		}(N)
	}
	swg.Wait()
}
