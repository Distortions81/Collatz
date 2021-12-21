package main

import (
	"fmt"
	"math/big"
	"runtime"
	"sync"

	"github.com/remeh/sizedwaitgroup"
)

// Delay record, 2426
// http://www.ericr.nl/wondrous/delrecs.html
const startNumber = "9781262575275081247"

var MaxSteps *big.Int
var MaxLock sync.Mutex

func checkMaxSteps(seed *big.Int, steps *big.Int) {
	MaxLock.Lock()
	if steps.Cmp(MaxSteps) > 0 {
		MaxSteps.Set(steps)
		fmt.Printf("NEW RECORD: SEED: %v, STEPS: %v\n", seed, steps)
	}
	MaxLock.Unlock()
}

func collatz(seed *big.Int, i *big.Int, steps *big.Int) {

	if i.Cmp(big.NewInt(1)) < 1 {
		checkMaxSteps(seed, steps)
	} else if is_even(i) {
		i.Div(i, big.NewInt(2))
		collatz(seed, i, steps.Add(steps, big.NewInt(1)))
	} else {
		i.Mul(i, big.NewInt(3))
		i.Add(i, big.NewInt(1))
		collatz(seed, i, steps.Add(steps, big.NewInt(1)))
	}
}

func main() {
	MaxSteps = big.NewInt(0)
	numCPU := runtime.NumCPU()

	fmt.Printf("Found %v vCPUs.\n", numCPU)
	swg := sizedwaitgroup.New(numCPU)

	i := big.NewInt(0)
	i.SetString(startNumber, 10)
	i.Sub(i, big.NewInt(int64(numCPU)))
	for ; true; i.Add(i, big.NewInt(1)) {
		swg.Add()
		go func() {
			newInt := big.NewInt(0)
			newInt.Set(i)
			newSeed := big.NewInt(0)
			newSeed.Set(i)

			collatz(newSeed, newInt, big.NewInt(0))
			swg.Done()
		}()
	}
	swg.Wait()
}

func is_even(i *big.Int) bool {
	z := big.NewInt(0)
	z.Mod(i, big.NewInt(2))
	return z.Cmp(big.NewInt(0)) == 0
}
