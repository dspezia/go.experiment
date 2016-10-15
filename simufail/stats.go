package main

// Statistic contains the metric values associated to a single event.
type Statistic struct {
	n, dur int
	rat    float32
}

// Aggregate sums statistics
func (s *Statistic) Aggregate(other Statistic) {
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
	n   int
	r2  Statistic
	r2x Statistic
	r3  Statistic
}

// Reset zeroes the result object.
func (r *Result) Reset() {
	r.n = 0
	r.r2, r.r3, r.r2x = Statistic{}, Statistic{}, Statistic{}
}

// Aggregate sums statistics from other runs.
func (r *Result) Aggregate(other *Result) {
	r.n += other.n
	r.r2.Aggregate(other.r2)
	r.r2x.Aggregate(other.r2x)
	r.r3.Aggregate(other.r3)
}
