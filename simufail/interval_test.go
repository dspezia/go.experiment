package main

import (
	"math"
	"math/rand"
	"sort"
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
	col := Intervals{}
	for _, x := range tests {
		if b := col.AddFailure(x.t, x.mttr, true); b != x.expected {
			t.Errorf("Expected %t, got %t for %v", x.expected, b, x)
		}
		if len(col) != x.size {
			t.Errorf("Expected %d, got %d for %v", x.size, len(col), x)
		}
	}
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
}

func TestMutipleFailures(t *testing.T) {
	r := rand.New(rand.NewSource(0))
	col := Intervals{}
	col.AddFailures(10, r, 100)
	if len(col) != 10 {
		t.Errorf("Expected %d, got %d", 10, len(col))
	}
}
