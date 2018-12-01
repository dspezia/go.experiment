package main

import "fmt"

// Result contains the result of a simulation run.
// It can also aggregate results of multiple runs.
type Result struct {
	n int
}

// Reset zeroes the result object.
func (r *Result) Reset() {
	*r = Result{}
}

// Aggregate sums statistics from other runs.
func (r *Result) Aggregate(other *Result) {
	r.n += other.n
}

// FinalResult contains aggregated results with probabilites.
type FinalResult struct {
	Result
	z1Cnt int
	z1Sum int
}

// Update aggregates statistics and maintain probability counters.
func (r *FinalResult) Update(other *Result) {

	r.Result.Aggregate(other)

}

// Aggregate sums results from other runs
func (r *FinalResult) Aggregate(other *FinalResult) {

	r.Result.Aggregate(&(other.Result))

	r.z1Cnt += other.z1Cnt
	r.z1Sum += other.z1Sum
}

// Format generates a human-readable output
func (r *FinalResult) Format(f fmt.State, c rune) {

}
