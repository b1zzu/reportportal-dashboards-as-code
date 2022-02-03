package util

import (
	"testing"
)

type Object struct {
	Name string
}

type ComparableObject struct {
	Name string
}

func (left *ComparableObject) Equals(r Comparable) bool {

	// left can not be nil so if right is nil they are for sure different
	if r == nil {
		return false
	}

	right, ok := r.(*ComparableObject)
	if !ok {
		return false
	}

	return left.Name == right.Name
}

func (o *ComparableObject) Key() string {
	return o.Name
}

func TestCompareSlices(t *testing.T) {

	tests := []*struct {
		description string
		left        []*ComparableObject
		right       []*ComparableObject
		expexct     bool
	}{
		{
			description: "Compare equal slices should return true",
			left:        []*ComparableObject{{Name: "One"}, {Name: "Two"}},
			right:       []*ComparableObject{{Name: "One"}, {Name: "Two"}},
			expexct:     true,
		},
		{
			description: "Compare equal slices but with differnt order should return true",
			left:        []*ComparableObject{{Name: "One"}, {Name: "Two"}},
			right:       []*ComparableObject{{Name: "Two"}, {Name: "One"}},
			expexct:     true,
		},
		{
			description: "Compare different slices should return false",
			left:        []*ComparableObject{{Name: "One"}},
			right:       []*ComparableObject{{Name: "Two"}, {Name: "One"}},
			expexct:     false,
		},
		{
			description: "Compare nil slices should return true",
			left:        nil,
			right:       nil,
			expexct:     true,
		},
		{
			description: "Compare a nil slice with a normal slice should return false",
			left:        []*ComparableObject{{Name: "One"}},
			right:       nil,
			expexct:     false,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {

			left := make([]Comparable, len(test.left))
			for i := range test.left {
				left[i] = test.left[i]
			}

			right := make([]Comparable, len(test.right))
			for i := range test.right {
				right[i] = test.right[i]
			}

			r := CompareSlices(left, right)
			if r != test.expexct {
				t.Errorf("expected '%t' but got '%t'", test.expexct, r)
			}
		})
	}
}

func TestCompareStringSlices(t *testing.T) {

	tests := []*struct {
		description string
		left        []string
		right       []string
		expexct     bool
	}{
		{
			description: "Compare equal slices should return true",
			left:        []string{"One", "Two"},
			right:       []string{"One", "Two"},
			expexct:     true,
		},
		{
			description: "Compare equal slices but with differnt order should return true",
			left:        []string{"One", "Two"},
			right:       []string{"Two", "One"},
			expexct:     true,
		},
		{
			description: "Compare different slices should return false",
			left:        []string{"One", "Two"},
			right:       []string{"Three", "Two"},
			expexct:     false,
		},
		{
			description: "Compare different slices with different number lengths should return false",
			left:        []string{"One", "Two"},
			right:       []string{"Two"},
			expexct:     false,
		},
		{
			description: "Compare nil slices should return true",
			left:        nil,
			right:       nil,
			expexct:     true,
		},
		{
			description: "Compare a nil slice with a normal slice should return false",
			left:        []string{"One"},
			right:       nil,
			expexct:     false,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			r := CompareStringSlices(test.left, test.right)
			if r != test.expexct {
				t.Errorf("expected '%t' but got '%t'", test.expexct, r)
			}
		})
	}
}
