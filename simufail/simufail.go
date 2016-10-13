package main

import (
	"flag"
	"fmt"
	"math/rand"
)

var flagNbIter = flag.Int("n", 10000, "Number of iterations")

type Simulation struct {
	nNodes uint16
	nZones uint16
}

func NewSimulation() *Simulation {
	return &Simulation{nNodes: 9, nZones: 3}
}

func (s *Simulation) RunOnce() {

	ay := NewAvailabilityYear()
	r := rand.New(rand.NewSource(0))
	for i := 0; i < 10; i++ {
		fmt.Println("---------------")
		ay.Reset()
		ay.Build(r)
		ay.Simulate()
	}
}

func main() {

	flag.Parse()
	fmt.Println(*flagNbIter)

	s := NewSimulation()
	s.RunOnce()
}
