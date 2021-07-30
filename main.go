package main

import (
	"fmt"
	"runtime"
	"sync"
)

var MaxSteps uint64 = 0
var MaxLock sync.Mutex

const maxInt = 18446744073709551615 // 2^64-1

func checkMaxSteps(seed uint64, i uint64, steps uint64) {
	MaxLock.Lock()
	if steps > MaxSteps {
		MaxSteps = steps
		fmt.Printf("NEW RECORD LENGTH:\nSEED: %v, I: %v, STEPS: %v\n", seed, i, steps)
	}
	MaxLock.Unlock()
}

func collatz(seed uint64, i uint64, steps uint64) {

	if steps%1000000 == 0 {
		checkMaxSteps(seed, i, steps)
	}

	if i == 1 {
		checkMaxSteps(seed, i, steps)
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
	fmt.Println("Running: ")

	var workSize uint64 = maxInt / uint64(numCPU)

	for cpu := 1; cpu <= numCPU; cpu++ {
		workStart := workSize * uint64(cpu-1)
		workEnd := workStart + workSize - 1
		fmt.Printf("CPU: %v, Work area: %v to %v\n", cpu, workStart, workEnd)
		go func(workStart uint64, workEnd uint64, cpu int) {
			for x := workStart; x < workEnd; x++ {
				collatz(x, x, 0)
			}
			fmt.Printf("\nvCPU %v IS FINISHED\n", cpu)
		}(workStart, workEnd, cpu)
	}
	for {
		//Wait forever
		//TODO: Wait for all vCPUs waitgroup to finish
	}
}
