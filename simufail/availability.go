package main

import "math/rand"

// AvailabilityYear represents an availability for a given year.
type AvailabilityYear struct {
	nodes   []Intervals
	zones   []Intervals
	cluster Intervals
	tmp     Intervals
}

// Each node suffers from N_REBOOTS unexpected reboot a year, resulting in MTBR_REBOOT secs outage per node.
// Each node has a PROB_HW_FAILURE% chance a year to get a hardware failure resulting in MTBR_HW_FAILURE secs outage.
// Each availability zone is brought down N_ZONE_SHUTDOWNS times a year resulting in MTBR_ZONE_SHUTDOWN secs outage per zone.
// The cluster is upgraded N_ROLLING_UPGRADES times a year using rolling upgrade, resulting in MTBR_ROLLING_UPGRADE secs outage
// per node in sequence, separated by IDLE_ROLLING_UPGRADE secs idle periods.
// Availability zone shutdowns and couchbase rolling upgrades are scheduled events, so:
//   - they are mutually exclusive
//   - they are not scheduled when there is already a node down for any reason
const (
	N_NODES              = 12
	N_ZONES              = 3
	N_ROLLING_UPGRADES   = 2
	MTBR_ROLLING_UPGRADE = 2 * 60
	IDLE_ROLLING_UPGRADE = 30
	N_ZONE_SHUTDOWNS     = 1
	MTBR_ZONE_SHUTDOWN   = 3 * 3600
	N_REBOOTS            = 2
	MTBR_REBOOT          = 5 * 60
	PROB_HW_FAILURE      = 6
	MTBR_HW_FAILURE      = 3 * 24 * 3600
	ZONE_SIZE            = N_NODES / N_ZONES
)

// NewAvailabilityYear creates a new object representing a full year of availability.
func NewAvailabilityYear() *AvailabilityYear {
	const NINIT = 8
	ret := &AvailabilityYear{
		nodes:   make([]Intervals, N_NODES),
		zones:   make([]Intervals, N_ZONES),
		cluster: make([]Interval, 0, NINIT),
		tmp:     make([]Interval, 0, NINIT),
	}
	for i := 0; i < N_NODES; i++ {
		ret.nodes[i] = make([]Interval, 0, NINIT)
	}
	for i := 0; i < N_ZONES; i++ {
		ret.zones[i] = make([]Interval, 0, NINIT)
	}
	return ret
}

// Reset re-initializes the object.
func (ay *AvailabilityYear) Reset() {
	for i := 0; i < N_NODES; i++ {
		ay.nodes[i].Reset()
	}
	for i := 0; i < N_ZONES; i++ {
		ay.zones[i].Reset()
	}
	ay.cluster.Reset()
	ay.tmp.Reset()
}

// Build generates a simulation.
func (ay *AvailabilityYear) Build(r *rand.Rand) {
	ay.buildNodes(r)
	ay.buildGlobalEvents(r)
	ay.retrofitGlobalEvents(r)
	ay.tmp.Reset()
	ay.cluster.Reset()
}

// BuildNodes populates the initial node views.
func (ay *AvailabilityYear) buildNodes(r *rand.Rand) {
	// Each node suffers from N_REBOOTS unexpected reboot a year, resulting in MTBR_REBOOT secs outage per node.
	// Each node has a PROB_HW_FAILURE% chance a year to get a hardware failure resulting in MTBR_HW_FAILURE secs outage.
	for i := 0; i < N_NODES; i++ {
		ay.nodes[i].AddFailures(N_REBOOTS, r, MTBR_REBOOT)
		if r.Intn(100) < PROB_HW_FAILURE {
			ay.nodes[i].AddFailures(1, r, MTBR_HW_FAILURE)
		}
		ay.cluster = append(ay.cluster, ay.nodes[i]...)
	}
}

// BuildGlobalEvents generate global events such as zone shutdown and rolling upgrades.
func (ay *AvailabilityYear) buildGlobalEvents(r *rand.Rand) {
	// Each availability zone is brought down N_ZONE_SHUTDOWNS times a year resulting in MTBR_ZONE_SHUTDOWN secs outage per zone.
	// The cluster is upgraded N_ROLLING_UPGRADES times a year using rolling upgrade, resulting in MTBR_ROLLING_UPGRADE secs outage
	// per node in sequence, separated by IDLE_ROLLING_UPGRADE secs idle periods.
	// Availability zone shutdowns and couchbase rolling upgrades are scheduled events, so:
	//   - they are mutually exclusive
	//   - they are not scheduled when there is already a node down for any reason
	for i := 0; i < N_ROLLING_UPGRADES; i++ {
		for {
			tRU := ay.cluster.FindNonFailureTime(r)
			dur := uint32((N_NODES-1)*(MTBR_ROLLING_UPGRADE+IDLE_ROLLING_UPGRADE) + MTBR_ROLLING_UPGRADE)
			if ay.tmp.AddFailure(tRU, dur, true) {
				break
			}
		}
	}
	for i := 0; i < N_ZONE_SHUTDOWNS*N_ZONES; i++ {
		for {
			tZS := ay.cluster.FindNonFailureTime(r)
			if ay.tmp.AddFailure(tZS, MTBR_ZONE_SHUTDOWN, true) {
				break
			}
		}
	}
}

// BuildNodes update the node views according to the global events.
func (ay *AvailabilityYear) retrofitGlobalEvents(r *rand.Rand) {
	// Retrofit global events in the node views.
	idx := 0
	for i := 0; i < N_ROLLING_UPGRADES; i++ {
		tRU := ay.tmp[idx].beg
		for n := 0; n < N_NODES; n++ {
			ay.nodes[n].AddFailure(tRU, MTBR_ROLLING_UPGRADE, false)
			tRU += MTBR_ROLLING_UPGRADE + IDLE_ROLLING_UPGRADE
		}
		idx++
	}
	for i := 0; i < N_ZONE_SHUTDOWNS*N_ZONES; i++ {
		tZS := ay.tmp[idx].beg
		for n := 0; n < ZONE_SIZE; n++ {
			ay.nodes[n+N_ZONES*(i%N_ZONES)].AddFailure(tZS, MTBR_ZONE_SHUTDOWN, false)
		}
		idx++
	}
}

// Simulate calculate the result of the simulation
func (ay *AvailabilityYear) Simulate() {
	ay.normalize()
	ay.simulateZones()
	ay.simulateCluster()
}

// BuildNodes update the node views according to the global events.
func (ay *AvailabilityYear) normalize() {
	for i := 0; i < N_NODES; i++ {
		ay.nodes[i].Normalize(false)
		for j := range ay.nodes[i] {
			p := &(ay.nodes[i][j])
			p.cnt = 1
			p.ratio = 1.0 / float32(ZONE_SIZE)
		}
	}
}

// simulateZone calculates the result of the simulation for the zones
func (ay *AvailabilityYear) simulateZones() {
	for iZ := 0; iZ < N_ZONES; iZ++ {
		nodes := ay.nodes[iZ*ZONE_SIZE : (iZ+1)*ZONE_SIZE]
		//fmt.Println("nodes", nodes)
		z := ay.zones[iZ]
		z = append(z, nodes[0]...)
		for i := 1; i < len(nodes); i++ {
			ay.tmp.MergeNodes(z, nodes[i])
			z, ay.tmp = ay.tmp, z
		}
		ay.zones[iZ] = z
		//fmt.Println("zone", z)
	}
}

// simulateCluster calculates the result of the simulation for the cluster
func (ay *AvailabilityYear) simulateCluster() {
	c := ay.cluster
	c = append(c, ay.zones[0]...)
	for iZ := 1; iZ < len(ay.zones); iZ++ {
		ay.tmp.MergeZones(c, ay.zones[iZ])
		c, ay.tmp = ay.tmp, c
	}
	ay.cluster = c
	//fmt.Println("cluster", c)
}
