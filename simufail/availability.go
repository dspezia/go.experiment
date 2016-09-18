package main

// AvailabilityYear represents an availability for a given year
type AvailabilityYear struct {
	nodes   []Intervals
	zones   []Intervals
	cluster Intervals
}

func NewAvailabilityYear(nodes, zones int) *AvailabilityYear {
	const NINIT = 8
	ret := &AvailabilityYear{
		nodes:   make([]Intervals, nodes),
		zones:   make([]Intervals, zones),
		cluster: make([]Interval, 0, NINIT),
	}
	for i := 0; i < nodes; i++ {
		ret.nodes[i] = make([]Interval, 0, NINIT)
	}
	for i := 0; i < zones; i++ {
		ret.zones[i] = make([]Interval, 0, NINIT)
	}
	return ret
}
