package main

import "fmt"

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

// Format generates a human-readable output
func (r *FinalResult) Format(f fmt.State, c rune) {

	fmt.Fprintf(f, "Number of simulations:            %d\n", r.Result.n)
	fmt.Fprintf(f, "Average node failures per year:   %.2f\n", float64(r.z1Sum)/float64(r.z1Cnt))
	fmt.Fprintln(f)

	r.displayRangeProb(f, N_ZONES)
	fmt.Fprintln(f)
	r.displayRangeProb(f, N_ZONES-1)
	fmt.Fprintln(f)
	r.displayRangeProb(f, 1)
}

func (r *FinalResult) displayRangeProb(f fmt.State, nz int) {

	if nz == 1 {
		fmt.Fprintf(f, "Failures involving at least 2 nodes\n\n")
		displayAverages(f, nz, r.Result.atLeast2)
	} else {
		fmt.Fprintf(f, "Failures on %d zones\n\n", nz)
		displayAverages(f, nz, r.Result.outages[nz-1])
	}

	fmt.Fprintf(f, "\nHistogram of probability\n")

	t := r.proba[nz-1][:]
	count := float64(r.Result.n)
	more := true

	for i, x := range t[1:] {
		if x == 0 {
			more = false
			break
		}
		fmt.Fprintf(f, "At least %2d occurences:  %8.4f %%\n", i+1, 100.0*float64(sum(i+1, t))/count)
	}
	if more && t[0] != 0 {
		fmt.Fprintf(f, "More occurences:         %8.4f %%\n", 100.0*float64(t[0])/count)
	}
}

func displayAverages(f fmt.State, nz int, s Statistic) {

	imp := s.rat / float64(s.n)
	dur := float64(s.dur) / float64(s.n)

	fmt.Fprintf(f, "Average %% of keys impacted:      %3.4f %%\n", 100.0*imp)
	fmt.Fprintf(f, "Average duration of the event:   %.0f secs\n", dur)

	if nz == N_ZONES {
		fmt.Fprintf(f, "Average number of records lost:  %3.4f\n", float64(THROUGHPUT)*imp)
	} else {
		fmt.Fprintf(f, "Average number of transactions:  %3.4f\n", float64(THROUGHPUT)*dur*imp)
	}
}

func sum(i int, t []int) int {
	n := t[0]
	for _, x := range t[i:] {
		n += x
	}
	return n
}
