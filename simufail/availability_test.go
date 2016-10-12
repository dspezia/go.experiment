package main

import (
	"math/rand"
	"sort"
	"testing"
)

func BenchmarkAvailAlloc(b *testing.B) {
	ay := NewAvailabilityYear()
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

func BenchmarkFailures(b *testing.B) {
	ay := NewAvailabilityYear()
	r := rand.New(rand.NewSource(0))
	for n := 0; n < b.N; n++ {
		ay.Reset()
		for i := 0; i < len(ay.nodes); i++ {
			ay.nodes[i].AddFailures(3, r, 600)
			sort.Sort(ay.nodes[i])
		}
	}
}

func BenchmarkBuild(b *testing.B) {
	ay := NewAvailabilityYear()
	r := rand.New(rand.NewSource(0))
	for n := 0; n < b.N; n++ {
		ay.Reset()
		ay.Build(r)
	}
}
