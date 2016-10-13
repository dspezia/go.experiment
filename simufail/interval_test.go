package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"testing"
)

func TestOverlap(t *testing.T) {
	tests := []struct {
		a, b     Interval
		expected bool
	}{
		{Interval{beg: 1, end: 10}, Interval{beg: 20, end: 30}, false},
		{Interval{beg: 20, end: 30}, Interval{beg: 0, end: 10}, false},
		{Interval{beg: 1, end: 10}, Interval{beg: 5, end: 30}, true},
		{Interval{beg: 20, end: 30}, Interval{beg: 0, end: 25}, true},
		{Interval{beg: 10, end: 20}, Interval{beg: 10, end: 20}, true},
		{Interval{beg: 0, end: 30}, Interval{beg: 10, end: 20}, true},
		{Interval{beg: 0, end: 10}, Interval{beg: 10, end: 20}, false},
		{Interval{beg: 20, end: 30}, Interval{beg: 10, end: 20}, false},
	}

	for _, x := range tests {
		if b := x.a.Overlap(x.b); b != x.expected {
			t.Errorf("Expected %t, got %t for %v %v", x.expected, b, x.a, x.b)
		}
	}
}

func TestContiguous(t *testing.T) {
	tests := []struct {
		a, b     Interval
		expected bool
	}{
		{Interval{beg: 1, end: 10}, Interval{beg: 10, end: 30}, true},
		{Interval{beg: 1, end: 10}, Interval{beg: 11, end: 30}, false},
		{Interval{beg: 10, end: 20}, Interval{beg: 0, end: 10}, true},
		{Interval{beg: 10, end: 20}, Interval{beg: 0, end: 9}, false},
		{Interval{beg: 10, end: 20}, Interval{beg: 10, end: 20}, false},
	}

	for _, x := range tests {
		if b := x.a.Contiguous(x.b); b != x.expected {
			t.Errorf("Expected %t, got %t for %v %v", x.expected, b, x.a, x.b)
		}
	}
}

func TestInclude(t *testing.T) {
	tests := []struct {
		a        Interval
		b        uint32
		expected bool
	}{
		{Interval{beg: 10, end: 20}, 1, false},
		{Interval{beg: 10, end: 20}, 9, false},
		{Interval{beg: 10, end: 20}, 10, true},
		{Interval{beg: 10, end: 20}, 15, true},
		{Interval{beg: 10, end: 20}, 20, true},
		{Interval{beg: 10, end: 20}, 21, false},
		{Interval{beg: 10, end: 20}, 30, false},
	}

	for _, x := range tests {
		if b := x.a.Include(x.b); b != x.expected {
			t.Errorf("Expected %t, got %t for %v %v", x.expected, b, x.a, x.b)
		}
	}
}

func TestIntervals(t *testing.T) {
	tests := []struct {
		t, mttr  uint32
		size     int
		expected bool
	}{
		{0, 100, 1, true},
		{200, 100, 2, true},
		{210, 150, 2, false},
		{50, 100, 2, false},
		{500, 200, 3, true},
		{120, 50, 4, true},
		{400, 50, 5, true},
		{800, math.MaxInt32, 6, true},
		{750, 100, 6, false},
	}

	// Check AddFailure
	col := Intervals{}
	for _, x := range tests {
		if b := col.AddFailure(x.t, x.mttr, true); b != x.expected {
			t.Errorf("Expected %t, got %t for %v", x.expected, b, x)
		}
		if len(col) != x.size {
			t.Errorf("Expected %d, got %d for %v", x.size, len(col), x)
		}
	}

	// Check sorting
	sort.Sort(col)
	prev := col[0]
	for _, x := range col[1:] {
		if prev.beg > x.beg {
			t.Errorf("%v not correctly ordered", x)
		}
		prev = x
	}
	last := col[len(col)-1]
	if last.end != MAXSECS {
		t.Errorf("Expected %d, got %d for %v", MAXSECS, last.end, last)
	}

	// Check time collision
	ttests := []struct {
		t        uint32
		expected bool
	}{
		{1, true},
		{100, true},
		{101, false},
		{125, true},
		{475, false},
		{550, true},
		{900, true},
	}
	for _, x := range ttests {
		if b := col.CheckCollisionTime(x.t); b != x.expected {
			t.Errorf("Expected %t, got %t for %v", x.expected, b, x)
		}
	}
}

func TestMultipleFailures(t *testing.T) {
	r := rand.New(rand.NewSource(0))
	col := Intervals{}
	col.AddFailures(10, r, 100)
	if len(col) != 10 {
		t.Errorf("Expected %d, got %d", 10, len(col))
	}
}

func TestFindNonFailureTime(t *testing.T) {
	r := rand.New(rand.NewSource(0))
	col := Intervals{}
	col.AddFailures(20, r, 3600)
	sort.Sort(col)
	for i := 0; i < 100; i++ {
		ts := col.FindNonFailureTime(r)
		if col.CheckCollisionTime(ts) {
			t.Errorf("Got unexpected collision for %v", ts)
		}
	}
}

func parseIntervals(s string) Intervals {
	out := Intervals{}
	for _, x := range strings.Fields(s) {
		var val Interval
		fmt.Sscanf(x, "%d,%d,%d,%f", &val.beg, &val.end, &val.cnt, &val.ratio)
		out = append(out, val)
	}
	return out
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		in, expected string
	}{
		{"", "[]"},
		{"1,10 20,30", "[{1 10 0 0} {20 30 0 0}]"},
		{"1,10 10,30", "[{1 30 0 0}]"},
		{"1,10 5,15 14,20", "[{1 20 0 0}]"},
		{"1,10 15,25 25,30 30,40 60,100", "[{1 10 0 0} {15 40 0 0} {60 100 0 0}]"},
		{"25,32 15,25 1,10 50,60 30,40", "[{1 10 0 0} {15 40 0 0} {50 60 0 0}]"},
		{"10,20 9,21 8,22 22,30 22,40 1,10", "[{1 40 0 0}]"},
	}
	for _, x := range tests {
		in := parseIntervals(x.in)
		in.Normalize(false)
		out := fmt.Sprintf("%v", in)
		if out != x.expected {
			t.Errorf("Expected %v, got %v", x.expected, out)
		}
	}
}

func TestMergeNodes(t *testing.T) {
	tests := []struct {
		a, b, expected string
	}{
		{"", "1,10,1,0.333 30,100,1,0.333", "[{1 10 1 0.333} {30 100 1 0.333}]"},
		{"1,10,1,0.333 30,100,1,0.333", "", "[{1 10 1 0.333} {30 100 1 0.333}]"},
		{"1,10,1,0.333 20,30,1,0.333", "1,50,1,0.333", "[{1 10 1 0.666} {10 20 1 0.333} {20 30 1 0.666} {30 50 1 0.333}]"},
		{"10,20,1,0.333 30,40,1,0.333", "20,30,1,0.333 35,45,1,0.333", "[{10 20 1 0.333} {20 30 1 0.333} {30 35 1 0.333} {35 40 1 0.666} {40 45 1 0.333}]"},
		{"20,30,1,0.333 35,45,1,0.333", "10,20,1,0.333 30,40,1,0.333", "[{10 20 1 0.333} {20 30 1 0.333} {30 35 1 0.333} {35 40 1 0.666} {40 45 1 0.333}]"},
		{"30,40,1,0.333 40,50,1,0.333", "10,35,1,0.333 38,45,1,0.333 48,60,1,0.333", "[{10 30 1 0.333} {30 35 1 0.666} {35 38 1 0.333} {38 40 1 0.666} {40 45 1 0.666} {45 48 1 0.333} {48 50 1 0.666} {50 60 1 0.333}]"},
		{"10,22,1,0.333", "15,22,1,0.333 30,40,1,0.333", "[{10 15 1 0.333} {15 22 1 0.666} {30 40 1 0.333}]"},
		{"10,22,1,0.333", "15,22,1,0.333", "[{10 15 1 0.333} {15 22 1 0.666}]"},
	}

	for _, x := range tests {
		a, b := parseIntervals(x.a), parseIntervals(x.b)
		var res Intervals
		res.MergeNodes(a, b)
		out := fmt.Sprintf("%v", res)
		if out != x.expected {
			t.Errorf("Expected %v, got %v", x.expected, out)
		}
	}
}
