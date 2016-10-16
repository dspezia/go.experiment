package main

// Statistic contains the metric values associated to a single event.
type Statistic struct {
	n, dur int
	rat    float32
}

// Aggregate sums statistics
func (s *Statistic) Aggregate(other *Statistic) {
	s.n += other.n
	s.dur += other.dur
	s.rat += other.rat
	if s.n < 0 || s.dur < 0 {
		panic("Statistic.Aggregate: integer overflow")
	}
}

// Update aggregates results with previous runs.
func (s *Statistic) Update(x Interval) {
	s.n++
	s.dur += int(x.end - x.beg)
	s.rat += x.ratio
}

// Result contains the result of a simulation run.
// It can also aggregate results of multiple runs.
type Result struct {
	n        int
	outages  [N_ZONES]Statistic
	failures [N_ZONES]Statistic
}

// Reset zeroes the result object.
func (r *Result) Reset() {
	r.n = 0
	for i := range r.outages {
		r.outages[i] = Statistic{}
	}
	for i := range r.failures {
		r.failures[i] = Statistic{}
	}
}

// Aggregate sums statistics from other runs.
func (r *Result) Aggregate(other *Result) {
	r.n += other.n
	for i := range r.outages {
		r.outages[i].Aggregate(&(other.outages[i]))
	}
	for i := range r.failures {
		r.failures[i].Aggregate(&(other.failures[i]))
	}
}
