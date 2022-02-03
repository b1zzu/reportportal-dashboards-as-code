package rpdac

import (
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

func testEqual(t *testing.T, got, want interface{}) {
	t.Helper()

	if got != want {
		t.Errorf("Want \"%s\" but got \"%s\"", want, got)
	}
}
