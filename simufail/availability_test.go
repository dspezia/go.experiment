package main

import (
	"sort"
	"testing"
)

func BenchmarkAvailAlloc(b *testing.B) {
	ay := NewAvailabilityYear(9, 3)
	for n := 0; n < b.N; n++ {
		ay.Reset()
		for i := 0; i < len(ay.nodes); i++ {
			ay.nodes[i].AddFailure(0, 200, true)
			ay.nodes[i].AddFailure(1000, 200, true)
			ay.nodes[i].AddFailure(2000, 200, true)
			sort.Sort(ay.nodes[i])
		}
	}
}
