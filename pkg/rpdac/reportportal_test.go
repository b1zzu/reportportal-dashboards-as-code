package rpdac

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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

func writeTmpFile(t *testing.T, pattern string, content string) (name string, clean func()) {
	t.Helper()

	f, c := tmpFile(t, pattern)

	err := ioutil.WriteFile(f, []byte(content), 0644)
	if err != nil {
		t.Errorf("failed to write file '%s': %s", f, err)
	}
	return f, c
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

func TestExport_Dashboard(t *testing.T) {
	file, cleanFile := tmpFile(t, "dashboard")
	defer cleanFile()

	r := NewReportPortal(nil)

	r.Dashboard = &MockService{
		GetM: func(project string, id int) (Object, error) {
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

	err := r.Export(DashboardKind, "test_project", 3, "", file)
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

func TestExport_DashboardByName(t *testing.T) {
	file, cleanFile := tmpFile(t, "dashboard")
	defer cleanFile()

	r := NewReportPortal(nil)

	r.Dashboard = &MockService{
		GetByNameM: func(project string, name string) (Object, error) {
			testEqual(t, project, "test_project")
			testEqual(t, name, "MK E2E Tests Overview")
			return &Dashboard{
				Kind:        DashboardKind,
				Name:        "MK E2E Tests Overview",
				Description: "",
				Widgets:     []*Widget{},
			}, nil
		},
	}

	err := r.Export(DashboardKind, "test_project", -1, "MK E2E Tests Overview", file)
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

func TestExport_Filter(t *testing.T) {
	file, cleanFile := tmpFile(t, "filter")
	defer cleanFile()

	r := NewReportPortal(nil)
	r.Filter = &MockService{
		GetM: func(project string, id int) (Object, error) {
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

	err := r.Export(FilterKind, "test_project", 3, "", file)
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

func TestCreate_Dashboard(t *testing.T) {

	input := `kind: Dashboard
name: MK E2E Tests Overview
description: ""
widgets: []
`

	file, cleanFile := writeTmpFile(t, "dashboard", input)
	defer cleanFile()

	mockDashboard := &MockService{
		CreateM: func(project string, o Object) error {
			testEqual(t, project, "test_project")
			testDeepEqual(t, o.(*Dashboard), &Dashboard{
				Kind:        DashboardKind,
				Name:        "MK E2E Tests Overview",
				Description: "",
				Widgets:     []*Widget{},
			}, cmp.AllowUnexported(Dashboard{}))
			return nil
		},
	}

	r := NewReportPortal(nil)
	r.Dashboard = mockDashboard

	err := r.Create("test_project", file)
	if err != nil {
		t.Errorf("Create retunred error: %s", err)
	}

	testDeepEqual(t, mockDashboard.Counter, MockServiceCounter{Create: 1})
}

func TestCreate_Filter(t *testing.T) {

	input := `kind: Filter
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

	file, cleanFile := writeTmpFile(t, "filter", input)
	defer cleanFile()

	mockFilter := &MockService{
		CreateM: func(project string, o Object) error {
			testEqual(t, project, "test_project")
			testDeepEqual(t, o.(*Filter), &Filter{
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
			}, cmp.AllowUnexported(Filter{}))
			return nil
		},
	}

	r := NewReportPortal(nil)
	r.Filter = mockFilter

	err := r.Create("test_project", file)
	if err != nil {
		t.Errorf("Create retunred error: %s", err)
	}

	testDeepEqual(t, mockFilter.Counter, MockServiceCounter{Create: 1})
}

func TestApply_Skip(t *testing.T) {

	input := `kind: Dashboard
name: MK E2E Tests Overview
description: Test desc
`

	file, cleanFile := writeTmpFile(t, "dashboard", input)
	defer cleanFile()

	mockService := &MockService{
		GetByNameM: func(project, name string) (Object, error) {
			testEqual(t, project, "test_project")
			testEqual(t, name, "MK E2E Tests Overview")

			return &Dashboard{
				Kind:        DashboardKind,
				Name:        "MK E2E Tests Overview",
				Description: "Test desc",
			}, nil
		},
	}

	r := NewReportPortal(nil)
	r.Dashboard = mockService

	err := r.Apply("test_project", file, false)
	if err != nil {
		t.Errorf("Apply retunred error: %s", err)
	}

	testDeepEqual(t, mockService.Counter, MockServiceCounter{GetByName: 1})
}

func TestApply_Update(t *testing.T) {

	input := `kind: Dashboard
name: MK E2E Tests Overview
description: Test new desc
`

	file, cleanFile := writeTmpFile(t, "dashboard", input)
	defer cleanFile()

	mockService := &MockService{
		GetByNameM: func(project, name string) (Object, error) {
			testEqual(t, project, "test_project")
			testEqual(t, name, "MK E2E Tests Overview")

			return &Dashboard{
				Kind:        DashboardKind,
				Name:        "MK E2E Tests Overview",
				Description: "Test old desc",
			}, nil
		},
		UpdateM: func(project string, current, target Object) error {
			testEqual(t, project, "test_project")
			testDeepEqual(t, current, &Dashboard{
				Kind:        DashboardKind,
				Name:        "MK E2E Tests Overview",
				Description: "Test old desc",
			}, cmpopts.IgnoreUnexported(Dashboard{}))

			testDeepEqual(t, target, &Dashboard{
				Kind:        DashboardKind,
				Name:        "MK E2E Tests Overview",
				Description: "Test new desc",
			}, cmpopts.IgnoreUnexported(Dashboard{}))
			return nil
		},
	}

	r := NewReportPortal(nil)
	r.Dashboard = mockService

	err := r.Apply("test_project", file, false)
	if err != nil {
		t.Errorf("Apply retunred error: %s", err)
	}

	testDeepEqual(t, mockService.Counter, MockServiceCounter{GetByName: 1, Update: 1})
}

func TestApply_Create(t *testing.T) {

	input := `kind: Dashboard
name: MK E2E Tests Overview
description: Test desc
`

	file, cleanFile := writeTmpFile(t, "dashboard", input)
	defer cleanFile()

	mockService := &MockService{
		GetByNameM: func(project, name string) (Object, error) {
			testEqual(t, project, "test_project")
			testEqual(t, name, "MK E2E Tests Overview")
			return nil, nil
		},
		CreateM: func(project string, o Object) error {
			testEqual(t, project, "test_project")
			testDeepEqual(t, o, &Dashboard{
				Kind:        DashboardKind,
				Name:        "MK E2E Tests Overview",
				Description: "Test desc",
			}, cmpopts.IgnoreUnexported(Dashboard{}))
			return nil
		},
	}

	r := NewReportPortal(nil)
	r.Dashboard = mockService

	err := r.Apply("test_project", file, false)
	if err != nil {
		t.Errorf("Apply retunred error: %s", err)
	}

	testDeepEqual(t, mockService.Counter, MockServiceCounter{GetByName: 1, Create: 1})
}

func tempDir(t *testing.T) (string, func()) {
	t.Helper()

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("failed to create tmp directory: %s", err)
	}
	return dir, func() { os.RemoveAll(dir) }
}

func mkdir(t *testing.T, name string) {
	t.Helper()

	err := os.Mkdir(name, 0755)
	if err != nil {
		t.Fatalf("failed to create dir '%s': %s", name, err)
	}
}

func writeFile(t *testing.T, filename, data string) {
	t.Helper()

	err := ioutil.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		t.Fatalf("failed to write file '%s': %s", filename, err)
	}
}

func TestApply_Directory(t *testing.T) {

	dir, clean := tempDir(t)
	defer clean()

	writeFile(t, dir+"/dashboard.yml", `kind: Dashboard
name: Test
`)
	writeFile(t, dir+"/skipme.xml", `Some randome stuff`)
	mkdir(t, dir+"/subfolder")
	writeFile(t, dir+"/subfolder/filter.yaml", `kind: Filter
name: Test
`)

	mockDashboardService := &MockService{
		GetByNameM: func(project, name string) (Object, error) {
			testEqual(t, project, "test_project")
			testEqual(t, name, "Test")
			return nil, nil
		},
		CreateM: func(project string, o Object) error {
			testEqual(t, project, "test_project")
			testDeepEqual(t, o, &Dashboard{
				Kind: DashboardKind,
				Name: "Test",
			}, cmpopts.IgnoreUnexported(Dashboard{}))
			return nil
		},
	}
	mockFilterService := &MockService{
		GetByNameM: func(project, name string) (Object, error) {
			testEqual(t, project, "test_project")
			testEqual(t, name, "Test")
			return &Filter{
				Kind: FilterKind,
				Name: "Test",
			}, nil
		},
	}
	r := NewReportPortal(nil)
	r.Dashboard = mockDashboardService
	r.Filter = mockFilterService

	err := r.Apply("test_project", dir, true)
	if err != nil {
		t.Errorf("Apply retunred error: %s", err)
	}

	testDeepEqual(t, mockDashboardService.Counter, MockServiceCounter{GetByName: 1, Create: 1})
	testDeepEqual(t, mockFilterService.Counter, MockServiceCounter{GetByName: 1})
}

func TestApply_DirectoryWithoutRecursive(t *testing.T) {

	dir, clean := tempDir(t)
	defer clean()
	r := NewReportPortal(nil)

	err := r.Apply("test_project", dir, false)
	if err == nil {
		t.Errorf("Want err but got nil")
	} else {
		want := fmt.Sprintf("error '%s' is a directory, use the `-r` option if you want to recursive apply all object in the directory", dir)
		if err.Error() != want {
			t.Errorf("Want err '%s' but got '%s'", want, err)
		}
	}
}

func TestApply_DirectoryWithWrongObject(t *testing.T) {

	dir, clean := tempDir(t)
	defer clean()

	writeFile(t, dir+"/dashboard.yml", `kind: Something
name: Test
`)
	writeFile(t, dir+"/skipme.xml", `Some randome stuff`)
	mkdir(t, dir+"/subfolder")
	writeFile(t, dir+"/subfolder/filter.yaml", `kind: Filter
name: Test
`)

	mockDashboardService := &MockService{}
	mockFilterService := &MockService{
		GetByNameM: func(project, name string) (Object, error) {
			testEqual(t, project, "test_project")
			testEqual(t, name, "Test")
			return &Filter{
				Kind: FilterKind,
				Name: "Test",
			}, nil
		},
	}
	r := NewReportPortal(nil)
	r.Dashboard = mockDashboardService
	r.Filter = mockFilterService

	err := r.Apply("test_project", dir, true)
	if err == nil {
		t.Errorf("Want err but got nil")
	} else {
		want := "error applying one or more objects"
		if err.Error() != want {
			t.Errorf("Want err '%s' but got '%s'", want, err)
		}
	}

	testDeepEqual(t, mockDashboardService.Counter, MockServiceCounter{})
	testDeepEqual(t, mockFilterService.Counter, MockServiceCounter{GetByName: 1})
}
