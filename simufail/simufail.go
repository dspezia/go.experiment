package main

import (
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

var (
	flagNbIter   = flag.Int("n", 32, "Number of iterations")
	flagBatch    = flag.Int("b", 100000, "Size of batch")
	flagSeed     = flag.Int("s", 1, "Seed for random numbers")
	flagParallel = flag.Int("p", runtime.NumCPU(), "Parallelism level")
)

var seed int

type Simulator struct {
	id  int
	r   *rand.Rand
	ay  *AvailabilityYear
	res FinalResult
}

func NewSimulator(n int) *Simulator {

	return &Simulator{
		id: n,
		r:  rand.New(rand.NewSource(int64(n*17 + seed))),
		ay: NewAvailabilityYear(),
	}
}

func (s *Simulator) Run(n int) {
	for i := 0; i < *flagBatch; i++ {
		s.ay.Reset()
		s.ay.Build(s.r)
		s.ay.Simulate()
		s.ay.Evaluate()
		s.res.Update(&(s.ay.res))
	}
}

func handleWorker(n int, input chan int, done chan FinalResult) {
	simu := NewSimulator(n)
	for x := range input {
		simu.Run(x)
	}
	done <- simu.res
}

func main() {

	flag.Parse()
	fmt.Printf("Starting %d iterations over %d threads with %d batch size\n", *flagNbIter, *flagParallel, *flagBatch)

	seed = *flagSeed
	if seed == 0 {
		seed = int(time.Now().UnixNano())
	}

	input := make(chan int, *flagParallel*2)
	done := make(chan FinalResult, *flagParallel*2)
	for i := 0; i < *flagParallel; i++ {
		go handleWorker(i, input, done)
	}

	for i := 0; i < *flagNbIter; i++ {
		input <- i
		fmt.Printf("Iteration %d\r", i)
	}
	close(input)

	var result FinalResult
	for i := 0; i < *flagParallel; i++ {
		r := <-done
		result.Aggregate(&r)
	}

	fmt.Println("\nDone\n")
	fmt.Println(&result)
}
