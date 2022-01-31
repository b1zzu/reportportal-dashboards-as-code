package rpdac

import (
	"log"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type Tesso struct {
	M int
}

func testDeepEqual(t *testing.T, got, want interface{}, opts ...cmp.Option) {
	t.Helper()

	if !cmp.Equal(got, want, opts...) {
		t.Errorf("Want (+) but got (-): %s", cmp.Diff(got, want, opts...))
	}
}

func TestDeepEqual(t *testing.T) {

	a := &Tesso{M: 3}
	b := &Tesso{M: 1}

	first := []*Tesso{a, b}
	second := first

	sort.Slice(second, func(a, b int) bool { return a < b })

	log.Printf("first: %+v", first)
	log.Printf("second: %+v", second)

	// left := []string{"test", "bla"}
	// right := []string{"test", "bla"}

	// if !reflect.DeepEqual(left, right) {
	// 	t.Errorf("Failed: left %+v is different from right %v", left, right)
	// }
}

func TestDeepEqual1(t *testing.T) {

	first := []int{3, 1}
	second := make([]int, len(first))
	copy(second, first)

	sort.Slice(second, func(a, b int) bool { return a > b })

	log.Printf("first: %+v", first)
	log.Printf("second: %+v", second)

	// left := []string{"test", "bla"}
	// right := []string{"test", "bla"}

	// if !reflect.DeepEqual(left, right) {
	// 	t.Errorf("Failed: left %+v is different from right %v", left, right)
	// }
}
