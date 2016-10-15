package main

import (
	"flag"
	"fmt"
	"math/rand"
	"runtime"
)

var (
	flagNbIter   = flag.Int("n", 32, "Number of iterations")
	flagBatch    = flag.Int("b", 100000, "Size of batch")
	flagSeed     = flag.Int("s", 0, "Seed for random numbers")
	flagParallel = flag.Int("p", runtime.NumCPU(), "Parallelism level")
)

type Simulator struct {
	id  int
	r   *rand.Rand
	ay  *AvailabilityYear
	res Result
}

func NewSimulator(n int) *Simulator {
	return &Simulator{
		id: n,
		r:  rand.New(rand.NewSource(int64(n*17 + *flagSeed))),
		ay: NewAvailabilityYear(),
	}
}

func (s *Simulator) Run(n int) {
	for i := 0; i < *flagBatch; i++ {
		s.ay.Reset()
		s.ay.Build(s.r)
		s.ay.Simulate()
		s.ay.Evaluate()
		s.res.Aggregate(&(s.ay.res))
	}
}

func handleWorker(n int, input chan int, done chan Result) {
	simu := NewSimulator(n)
	for x := range input {
		simu.Run(x)
	}
	done <- simu.res
}

func main() {

	flag.Parse()
	fmt.Printf("Starting %d iterations over %d threads with %d batch size\n", *flagNbIter, *flagParallel, *flagBatch)

	input := make(chan int, 16)
	done := make(chan Result)
	for i := 0; i < *flagParallel; i++ {
		go handleWorker(i, input, done)
	}

	for i := 0; i < *flagNbIter; i++ {
		input <- i
	}
	close(input)

	var result Result
	for i := 0; i < *flagParallel; i++ {
		r := <-done
		result.Aggregate(&r)
	}

	fmt.Println("Result:", result)
	fmt.Println("Done.")
}
