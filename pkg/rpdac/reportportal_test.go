package rpdac

import (
	"io/ioutil"
	"os"
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

func testFileContains(t *testing.T, file string, want string) {
	t.Helper()

	b, err := ioutil.ReadFile(file)
	if err != nil {
		t.Errorf("failed to read file '%s': %s", file, err)
	}
	got := string(b)

	if got != want {
		t.Errorf("Want '''%s''' from file '%s' but got '''%s'''", want, file, got)
	}
}

func tmpFile(t *testing.T, patter string) (name string, clean func()) {
	t.Helper()

	file, err := ioutil.TempFile("", "dashboard")
	if err != nil {
		t.Errorf("failed to create tmp file: %s", err)
		return "", func() {}
	}

	name = file.Name()
	clean = func() {
		err := os.Remove(name)
		if err != nil {
			t.Errorf("failed to clean tm file '%s': %s", name, err)
		}
	}

	return name, clean
}

func TestExportDashboard(t *testing.T) {
	file, cleanFile := tmpFile(t, "dashboard")
	defer cleanFile()

	r := NewReportPortal(nil)

	r.Dashboard = &MockDashboardService{
		GetDashboardM: func(project string, id int) (*Dashboard, error) {
			testEqual(t, project, "test_project")
			testEqual(t, id, 3)
			return &Dashboard{
				Kind:        DashboardKind,
				Name:        "MK E2E Tests Overview",
				Description: "",
				Widgets:     []*Widget{},
			}, nil
		},
	}

	err := r.Export(DashboardKind, "test_project", 3, file)
	if err != nil {
		t.Errorf("Export returned error: %s", err)
	}

	want := `kind: Dashboard
name: MK E2E Tests Overview
description: ""
widgets: []
`

	testFileContains(t, file, want)
}

func TestExportFilter(t *testing.T) {
	file, cleanFile := tmpFile(t, "filter")
	defer cleanFile()

	r := NewReportPortal(nil)

	r.Filter = &MockFilterService{
		GetFilterM: func(project string, id int) (*Filter, error) {
			testEqual(t, project, "test_project")
			testEqual(t, id, 3)
			return &Filter{
				Kind:        FilterKind,
				Name:        "mk-e2e-test-suite",
				Type:        "Launch",
				Description: "",
				Conditions: []FilterCondition{
					{FilteringField: "name", Condition: "eq", Value: "mk-e2e-test-suite"},
				},
				Orders: []FilterOrder{
					{SortingColumn: "startTime", IsAsc: false},
					{SortingColumn: "number", IsAsc: false},
				},
			}, nil
		},
	}

	err := r.Export(FilterKind, "test_project", 3, file)
	if err != nil {
		t.Errorf("Export returned error: %s", err)
	}

	want := `kind: Filter
name: mk-e2e-test-suite
type: Launch
description: ""
conditions:
- filteringfield: name
  condition: eq
  value: mk-e2e-test-suite
orders:
- sortingcolumn: startTime
  isasc: false
- sortingcolumn: number
  isasc: false
`

	testFileContains(t, file, want)
}
