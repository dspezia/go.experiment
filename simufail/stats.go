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

// FinalResult contains aggregated results with probabilites.
type FinalResult struct {
	Result
	proba [N_ZONES][N_HISTO]int
}

// Update aggregates statistics and maintain probability counters.
func (r *FinalResult) Update(other *Result) {

	r.Result.Aggregate(other)

	for i := 0; i < N_ZONES; i++ {
		n := other.failures[i].n
		if i == 0 {
			// Scale change since there are numerous events in this case
			n /= 10
		}
		if n > 0 {
			if n >= N_HISTO {
				n = 0
			}
			r.proba[i][n]++
		}
	}
}

// Aggregate sums results from other runs
func (r *FinalResult) Aggregate(other *FinalResult) {

	r.Result.Aggregate(&(other.Result))

	for i := range r.proba {
		for j := range r.proba[i] {
			r.proba[i][j] += other.proba[i][j]
		}
	}
}
