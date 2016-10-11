package main

// AvailabilityYear represents an availability for a given year
type AvailabilityYear struct {
	nodes   []Intervals
	zones   []Intervals
	cluster Intervals
}

// NewAvailabilityYear creates a new object representing a full year of availability
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

// Reset re-initializes the object
func (ay *AvailabilityYear) Reset() {
	for i := 0; i < len(ay.nodes); i++ {
		ay.nodes[i].Reset()
	}
	for i := 0; i < len(ay.zones); i++ {
		ay.zones[i].Reset()
	}
	ay.cluster.Reset()
}

// Build generates a simulation
func (ay *AvailabilityYear) Build() {
	// Each node suffers from 2 unexpected reboot a year, resulting in 5 min outage per node.
	// Each node has a 10% chance a year to get an hardware failure resulting in 3 days outage.
	// Each availability zone is brought down once a year resulting in 3 hours outage per zone.
	// The couchbase cluster is upgraded once a year using rolling upgrade, resulting in 2 min outage per node in sequence.
	// Availability zone shutdowns and couchbase rolling upgrades are scheduled events, so:
	//   - they are mutually exclusive
	//   - they are not scheduled when there is already a node down for any reason
}
