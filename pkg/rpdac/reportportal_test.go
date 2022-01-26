package rpdac

import (
	"log"
	"sort"
	"testing"
)

type Tesso struct {
	M int
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
