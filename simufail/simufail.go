package main

import (
	"flag"
	"fmt"
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

}

func main() {

	flag.Parse()
	fmt.Println(*flagNbIter)

	s := NewSimulation()
	s.RunOnce()
}
