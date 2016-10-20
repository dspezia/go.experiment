package main

// Statistic contains the metric values associated to a single event.
type Statistic struct {
	n, dur int
	rat    float64
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
	s.rat += float64(x.ratio)
}

// Result contains the result of a simulation run.
// It can also aggregate results of multiple runs.
type Result struct {
	n        int
	atLeast2 Statistic
	outages  [N_ZONES]Statistic
	failures [N_ZONES]Statistic
}

// Reset zeroes the result object.
func (r *Result) Reset() {
	*r = Result{}
}

// Aggregate sums statistics from other runs.
func (r *Result) Aggregate(other *Result) {
	r.n += other.n
	r.atLeast2.Aggregate(&(other.atLeast2))
	for i := range r.outages {
		r.outages[i].Aggregate(&(other.outages[i]))
	}
	for i := range r.failures {
		r.failures[i].Aggregate(&(other.failures[i]))
	}
}

// FinalResult contains aggregated results with probabilites.
type FinalResult struct {
	Result
	z1Cnt int
	z1Sum int
	proba [N_ZONES][N_HISTO]int
}

// Update aggregates statistics and maintain probability counters.
func (r *FinalResult) Update(other *Result) {

	r.Result.Aggregate(other)

	// Build an histogram for "at least 2" failures.
	// Expected zone shutdown are excluded.
	n := other.atLeast2.n - N_ZONE_SHUTDOWNS*N_ZONES
	r.calculate(0, n)

	// Node failure in a single zone, an histogram is useless.
	// Just calculate an average instead.
	n = other.failures[0].n
	if n > 0 {
		r.z1Cnt++
		r.z1Sum += n
	}

	// Failures in two or more zones, build an histogram.
	for i := 1; i < N_ZONES; i++ {
		n = other.failures[i].n
		r.calculate(i, n)
	}
}

func (r *FinalResult) calculate(i, n int) {
	if n > 0 {
		if n >= N_HISTO {
			n = 0
		}
		r.proba[i][n]++
	}
}

// Aggregate sums results from other runs
func (r *FinalResult) Aggregate(other *FinalResult) {

	r.Result.Aggregate(&(other.Result))

	r.z1Cnt += other.z1Cnt
	r.z1Sum += other.z1Sum

	for i := range r.proba {
		for j := range r.proba[i] {
			r.proba[i][j] += other.proba[i][j]
		}
	}
}
