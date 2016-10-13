package main

import (
	"math/rand"
	"sort"
)

// MAXSECS is the maximum number of seconds in a year.
const MAXSECS = 365 * 24 * 3600

// Interval represents an availability interval.
type Interval struct {
	beg, end uint32
	cnt      uint32
	ratio    float32
}

// Overlap checks the overlap between two intervals.
func (i Interval) Overlap(o Interval) bool {
	return i.end > o.beg && i.beg < o.end
}

// Contiguous checks the intervals are Contiguous.
func (i Interval) Contiguous(o Interval) bool {
	return i.end == o.beg || i.beg == o.end
}

// Include returns true if the interval include t.
func (i Interval) Include(t uint32) bool {
	return t >= i.beg && t <= i.end
}

// Normalize ensures the interval is within the range.
func (i *Interval) Normalize() {
	if i.end > MAXSECS {
		i.end = MAXSECS
	}
}

// Intervals is a slice of intervals.
type Intervals []Interval

// Reset cleans so the object can be reused.
func (s *Intervals) Reset() {
	*s = (*s)[:0]
}

// AddFailure adds a failure event avoiding collisions.
// Return false if an overlap check fails.
func (s *Intervals) AddFailure(t uint32, mttr uint32, check bool) bool {
	x := Interval{t, t + mttr, 0, 0.0}
	x.Normalize()
	if check && s.CheckCollision(x) {
		return false
	}
	*s = append(*s, x)
	return true
}

// CheckCollision returns true if an existing interval overlaps.
func (s Intervals) CheckCollision(x Interval) bool {
	for i := range s {
		if s[i].Overlap(x) {
			return true
		}
	}
	return false
}

// CheckCollisionTime returns true if t matches an existing interval.
func (s Intervals) CheckCollisionTime(t uint32) bool {
	for i := range s {
		if s[i].Include(t) {
			return true
		}
	}
	return false
}

// AddFailures adds multiple failures avoiding collisions.
func (s *Intervals) AddFailures(n int, r *rand.Rand, mttr uint32) {
	for n > 0 {
		t := uint32(r.Int31n(MAXSECS))
		if !s.AddFailure(t, mttr, true) {
			continue
		}
		n--
	}
}

// FindNonFailureTime returns a timestamp which does not match an existing interval.
func (s Intervals) FindNonFailureTime(r *rand.Rand) uint32 {
	for {
		t := uint32(r.Int31n(MAXSECS))
		if !s.CheckCollisionTime(t) {
			return t
		}
	}
}

// Normalize puts the list of intervals in canonical form
func (sp *Intervals) Normalize(sorted bool) {

	// Sort intervals
	s := *sp
	if len(s) == 0 {
		return
	}
	if !sorted {
		sort.Sort(s)
	}

	// Merge contiguous or overlapping intervals
	for i := 1; i < len(s); {
		if s[i].Overlap(s[i-1]) {
			if s[i-1].end < s[i].end {
				s[i-1].end = s[i].end
			}
			s = append(s[:i], s[i+1:]...)
			continue
		}
		if s[i].Contiguous(s[i-1]) {
			s[i-1].end = s[i].end
			s = append(s[:i], s[i+1:]...)
			continue
		}
		i++
	}
	*sp = s
}

// Equal returns true if the two objects are identical
func (s Intervals) Equal(other Intervals) bool {
	if len(s) != len(other) {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] != other[i] {
			return false
		}
	}
	return true
}

// MergeNodes merges two interval slices associated to two nodes of the same zone
func (s *Intervals) MergeNodes(a, b Intervals) {
	s.merge(a, b, func(x, y Interval) Interval {
		return Interval{x.beg, x.end, x.cnt, x.ratio + y.ratio}
	})
}

// MergeNodes merges two interval slices associated to two zones of the same cluster
func (s *Intervals) MergeZones(a, b Intervals) {
	s.merge(a, b, func(x, y Interval) Interval {
		return Interval{x.beg, x.end, x.cnt + y.cnt, x.ratio * y.ratio}
	})
}

// merge applies the generic merge algorithm of two interval slices.
// The two slices must be sorted.
func (s *Intervals) merge(a, b Intervals, gen func(x, y Interval) Interval) {
	s.Reset()
	i, j := 0, 0
	var x Interval
	for i < len(a) && j < len(b) {
		ab, ae, bb, be := a[i].beg, a[i].end, b[j].beg, b[j].end
		switch {
		case ae <= bb:
			*s = append(*s, a[i])
			i++
		case be <= ab:
			*s = append(*s, b[j])
			j++
		case ab == bb:
			switch {
			case ae == be:
				*s = append(*s, gen(a[i], b[j]))
				i++
				j++
			case ae < be:
				*s = append(*s, gen(a[i], b[j]))
				b[j].beg = ae
				i++
			default:
				*s = append(*s, gen(b[j], a[i]))
				a[i].beg = be
				j++
			}
		case ab < bb:
			x = a[i]
			x.end, a[i].beg = bb, bb
			*s = append(*s, x)
		default:
			x = b[j]
			x.end, b[j].beg = ab, ab
			*s = append(*s, x)
		}
	}
	switch {
	case i < len(a):
		*s = append(*s, a[i:]...)
	case j < len(b):
		*s = append(*s, b[j:]...)
	}
}

// Implements sort interface.
func (s Intervals) Len() int           { return len(s) }
func (s Intervals) Less(i, j int) bool { return s[i].beg < s[j].beg }
func (s Intervals) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
