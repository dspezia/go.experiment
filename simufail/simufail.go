package main

import (
	"flag"
	"fmt"
	"math/rand"
)

var (
	flagNbIter = flag.Int("n", 10000, "Number of iterations")
	flagSeed   = flag.Int("s", 0, "Seed for random numbers")
)

type Simulation struct {
	ay          *AvailabilityYear
	r2, r2x, r3 Result
}

func NewSimulation() *Simulation {
	return &Simulation{ay: NewAvailabilityYear()}
}

func (s *Simulation) RunOnce() {

	r := rand.New(rand.NewSource(int64(*flagSeed)))
	for i := 0; i < *flagNbIter; i++ {
		s.ay.Reset()
		s.ay.Build(r)
		s.ay.Simulate()
		s.ay.Evaluate()
		s.r2.Update(s.ay.r2)
		s.r2x.Update(s.ay.r2x)
		s.r3.Update(s.ay.r3)
	}
	fmt.Println(s.r2, s.r2x, s.r3)
}

func main() {

	flag.Parse()

	s := NewSimulation()
	s.RunOnce()
}
