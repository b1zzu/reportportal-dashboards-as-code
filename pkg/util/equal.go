package util

import (
	"sort"
)

type Comparable interface {
	Equals(c Comparable) bool

	// key is used sorting before compare slices and must be a unique identifier in the slice
	Key() string
}

func copySlice(s []Comparable) []Comparable {
	c := make([]Comparable, len(s))
	copy(c, s)
	return c
}

func sortSlice(s []Comparable) {
	sort.Slice(s, func(i, j int) bool {
		return s[i].Key() < s[j].Key()
	})
}

func CompareSlices(l []Comparable, r []Comparable) bool {

	// if one of the two is nil return true if both are nil otherwise return false
	if l == nil || r == nil {
		return l == nil && r == nil
	}

	// if they are not of the same lenght return false
	if len(l) != len(r) {
		return false
	}

	// copy the slice so we can sort them without changing the original slices
	left := copySlice(l)
	right := copySlice(r)

	sortSlice(left)
	sortSlice(right)

	for i, l := range left {
		r := right[i]

		// if one of the two is nil continue if both are nil otherwise return false
		if l == nil || r == nil {
			if l == nil && r == nil {
				continue
			} else {
				return false
			}
		}

		if !l.Equals(r) {
			return false
		}
	}
	return true
}

func CompareStringSlices(l []string, r []string) bool {

	// if one of the two is nil return true if both are nil otherwise return false
	if l == nil || r == nil {
		return l == nil && r == nil
	}

	// if they are not of the same lenght return false
	if len(l) != len(r) {
		return false
	}

	left := make([]string, len(l))
	copy(left, l)

	right := make([]string, len(r))
	copy(right, r)

	sort.StringSlice(left).Sort()
	sort.StringSlice(right).Sort()

	for i, l := range left {
		r := right[i]

		if l != r {
			return false
		}
	}
	return true
}
