package main

import (
	"flag"
	"fmt"
	"math/rand"
	"runtime"
)

var (
	flagNbIter   = flag.Int("n", 10, "Number of iterations")
	flagBatch    = flag.Int("b", 1000000, "Size of batch")
	flagSeed     = flag.Int("s", 0, "Seed for random numbers")
	flagParallel = flag.Int("p", runtime.NumCPU(), "Parallelism level")
)

type Simulation struct {
	id          int
	r           *rand.Rand
	ay          *AvailabilityYear
	r2, r2x, r3 Result
}

func NewSimulation(n int) *Simulation {
	return &Simulation{
		id: n,
		r:  rand.New(rand.NewSource(int64(n + *flagSeed))),
		ay: NewAvailabilityYear(),
	}
}

func (s *Simulation) Run(n int) {
	for i := 0; i < *flagBatch; i++ {
		s.ay.Reset()
		s.ay.Build(s.r)
		s.ay.Simulate()
		s.ay.Evaluate()
		s.r2.Update(s.ay.r2)
		s.r2x.Update(s.ay.r2x)
		s.r3.Update(s.ay.r3)
	}
}

func handleWorker(n int, input chan int, done chan bool) {
	simu := NewSimulation(n)
	for x := range input {
		simu.Run(x)
	}
	fmt.Println(n, simu.r2, simu.r2x, simu.r3)
	done <- true
}

func main() {

	flag.Parse()
	fmt.Printf("Starting %d iterations over %d threads with %d batch size\n", *flagNbIter, *flagParallel, *flagBatch)

	input := make(chan int, 16)
	done := make(chan bool)
	for i := 0; i < *flagParallel; i++ {
		go handleWorker(i, input, done)
	}
	for i := 0; i < *flagNbIter; i++ {
		input <- i
	}
	close(input)
	for i := 0; i < *flagParallel; i++ {
		_ = <-done
	}

	fmt.Println("Done.")
}
