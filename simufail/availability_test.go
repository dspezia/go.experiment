package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
)

func BenchmarkAddFailure(b *testing.B) {
	ay := NewAvailabilityYear()
	for n := 0; n < b.N; n++ {
		ay.Reset()
		for i := 0; i < len(ay.nodes); i++ {
			ay.nodes[i].AddFailure(0, 200, true)
			ay.nodes[i].AddFailure(1000, 200, true)
			ay.nodes[i].AddFailure(2000, 200, true)
		}
	}
}

func BenchmarkMultipleFailures(b *testing.B) {
	ay := NewAvailabilityYear()
	r := rand.New(rand.NewSource(0))
	for n := 0; n < b.N; n++ {
		ay.Reset()
		for i := 0; i < len(ay.nodes); i++ {
			ay.nodes[i].AddFailures(3, r, 600)
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

func BenchmarkSimulate(b *testing.B) {
	ay := NewAvailabilityYear()
	r := rand.New(rand.NewSource(0))
	for n := 0; n < b.N; n++ {
		ay.Reset()
		ay.Build(r)
		ay.Simulate()
	}
}

func TestSimulation(t *testing.T) {
	tests := []struct {
		geninput  func(i int) string
		genoutput func(i int) string
	}{
		{
			func(i int) string { return "100,200 1000,2000 10000,30000" },
			func(i int) string { return "100,200,3 1000,2000,3 10000,30000,3" },
		},
		{
			func(i int) string {
				switch {
				case i%ZONE_SIZE == 0:
					return "100,200 1000,2000 10000,30000"
				default:
					return "1000,3000"
				}
			},
			func(i int) string {
				return "100,200,3 1000,2000,3 2000,3000,3 10000,30000,3"
			},
		},
		{
			func(i int) string {
				switch {
				case i%ZONE_SIZE == 0 && i < 2*ZONE_SIZE:
					return "100,200 1000,2000 10000,30000"
				case i/ZONE_SIZE > 0:
					return "1000,3000"
				default:
					return ""
				}
			},
			func(i int) string {
				return "100,200,2 1000,2000,3 2000,3000,2 10000,30000,2"
			},
		},
		{
			func(i int) string {
				switch {
				case i == 0:
					return "100,200 1000,2000 10000,30000"
				case i%ZONE_SIZE == 1 && i < ZONE_SIZE*2:
					return "800,3000"
				case i > ZONE_SIZE:
					return "20000,40000"
				default:
					return ""
				}
			},
			func(i int) string {
				return "100,200,1 800,1000,2 1000,2000,2 2000,3000,2 10000,20000,1 20000,30000,3 30000,40000,2"
			},
		},
	}
	ay := NewAvailabilityYear()
	for iTest := 0; iTest < len(tests); iTest++ {
		ay.Reset()
		for i := 0; i < N_NODES; i++ {
			input := tests[iTest].geninput(i)
			ay.nodes[i] = parseIntervals(input)
		}
		ay.Simulate()
		var b bytes.Buffer
		for i, x := range ay.cluster {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d,%d,%d", x.beg, x.end, x.cnt)
		}
		output := tests[iTest].genoutput(iTest)
		if b.String() != output {
			t.Errorf("Expected %v, got %v for %d with %v", output, b.String(), iTest, ay.cluster)
		}
	}
}
