package rpdac

import (
	"errors"
	"testing"

	"github.com/b1zzu/reportportal-dashboards-as-code/pkg/reportportal"
	"github.com/google/go-cmp/cmp"
)

func TestGetFilter(t *testing.T) {

	filter := &reportportal.Filter{
		Owner: "dbizzarr",
		Share: true,
		ID:    2,
		Name:  "mk-e2e-test-suite",
		Conditions: []reportportal.FilterCondition{
			{
				FilteringField: "name",
				Condition:      "eq",
				Value:          "mk-e2e-test-suite",
			},
		},
		Orders: []reportportal.FilterOrder{
			{
				SortingColumn: "startTime",
				IsAsc:         false,
			},
			{
				SortingColumn: "number",
				IsAsc:         false,
			},
		},
		Type: "Launch",
	}

	mockFilter := &reportportal.MockFilterService{
		GetByIDM: func(projectName string, id int) (*reportportal.Filter, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, id, 2)
			return filter, nil, nil
		},
	}

	r := NewReportPortal(&reportportal.Client{
		Filter: mockFilter,
	})

	got, err := r.Filter.GetFilter("test_project", 2)
	if err != nil {
		t.Errorf("ReportPortal.GetFilter returned error: %v", err)
	}

	want := &Filter{
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
		origin: filter,
	}

	testDeepEqual(t, got, want, cmp.AllowUnexported(Filter{}))
}

func TestGetFilterByName(t *testing.T) {

	filter := &reportportal.Filter{
		Owner: "dbizzarr",
		Share: true,
		ID:    2,
		Name:  "mk-e2e-test-suite",
		Conditions: []reportportal.FilterCondition{
			{
				FilteringField: "name",
				Condition:      "eq",
				Value:          "mk-e2e-test-suite",
			},
		},
		Orders: []reportportal.FilterOrder{
			{
				SortingColumn: "startTime",
				IsAsc:         false,
			},
			{
				SortingColumn: "number",
				IsAsc:         false,
			},
		},
		Type: "Launch",
	}

	mockFilter := &reportportal.MockFilterService{
		GetByNameM: func(projectName string, name string) (*reportportal.Filter, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, name, "mk-e2e-test-suite")
			return filter, nil, nil
		},
	}

	r := NewReportPortal(&reportportal.Client{
		Filter: mockFilter,
	})

	got, err := r.Filter.GetFilterByName("test_project", "mk-e2e-test-suite")
	if err != nil {
		t.Errorf("ReportPortal.GetFilterByName returned error: %v", err)
	}

	want := &Filter{
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
		origin: filter,
	}

	testDeepEqual(t, got, want, cmp.AllowUnexported(Filter{}))
}

func TestGetFilterByName_NotFound(t *testing.T) {

	mockFilter := &reportportal.MockFilterService{
		GetByNameM: func(projectName string, name string) (*reportportal.Filter, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, name, "mk-e2e-test-suite")
			return nil, nil, reportportal.NewFilterNotFoundError(projectName, name)
		},
	}

	r := NewReportPortal(&reportportal.Client{
		Filter: mockFilter,
	})

	got, err := r.Filter.GetFilterByName("test_project", "mk-e2e-test-suite")
	if err != nil {
		t.Errorf("ReportPortal.GetFilterByName returned error: %v", err)
	}

	if got != nil {
		t.Errorf("ReportPortal.GetFilterByName want nil but got %+v", got)
	}
}

func TestGetFilterByName_Error(t *testing.T) {

	mockFilter := &reportportal.MockFilterService{
		GetByNameM: func(projectName string, name string) (*reportportal.Filter, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, name, "mk-e2e-test-suite")
			return nil, nil, errors.New("unexpected error")
		},
	}

	r := NewReportPortal(&reportportal.Client{
		Filter: mockFilter,
	})

	_, err := r.Filter.GetFilterByName("test_project", "mk-e2e-test-suite")
	if err == nil {
		t.Errorf("ReportPortal.GetFilterByName did not return the error")
	}
}

func TestCreateFilter(t *testing.T) {

	mockFilter := &reportportal.MockFilterService{
		CreateM: func(projectName string, f *reportportal.NewFilter) (int, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")

			want := &reportportal.NewFilter{
				Share: true,
				Name:  "mk-e2e-test-suite",
				Conditions: []reportportal.FilterCondition{
					{
						FilteringField: "name",
						Condition:      "eq",
						Value:          "mk-e2e-test-suite",
					},
				},
				Orders: []reportportal.FilterOrder{
					{
						SortingColumn: "startTime",
						IsAsc:         false,
					},
					{
						SortingColumn: "number",
						IsAsc:         false,
					},
				},
				Type: "Launch",
			}

			testDeepEqual(t, f, want)
			return 2, nil, nil
		},
	}

	r := NewReportPortal(&reportportal.Client{
		Filter: mockFilter,
	})

	inputFilter := &Filter{
		Kind:        FilterKind,
		Name:        "mk-e2e-test-suite",
		Description: "",
		Type:        "Launch",
		Conditions: []FilterCondition{
			{
				FilteringField: "name",
				Condition:      "eq",
				Value:          "mk-e2e-test-suite",
			},
		},
		Orders: []FilterOrder{
			{
				SortingColumn: "startTime",
				IsAsc:         false,
			},
			{
				SortingColumn: "number",
				IsAsc:         false,
			},
		},
	}

	err := r.Filter.CreateFilter("test_project", inputFilter)
	if err != nil {
		t.Errorf("ReportPortal.CreateFilter returned error: %v", err)
	}

	testDeepEqual(t, mockFilter.Counter, reportportal.MockFilterServiceCounter{Create: 1})
}

func TestApplyFilter_Create(t *testing.T) {

	mockFilter := &reportportal.MockFilterService{
		GetByNameM: func(projectName, name string) (*reportportal.Filter, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, name, "mk-e2e-test-suite")
			return nil, nil, reportportal.NewFilterNotFoundError(projectName, name)
		},
		CreateM: func(projectName string, f *reportportal.NewFilter) (int, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")

			want := &reportportal.NewFilter{
				Share: true,
				Name:  "mk-e2e-test-suite",
				Conditions: []reportportal.FilterCondition{
					{
						FilteringField: "name",
						Condition:      "eq",
						Value:          "mk-e2e-test-suite",
					},
				},
				Orders: []reportportal.FilterOrder{
					{
						SortingColumn: "startTime",
						IsAsc:         false,
					},
					{
						SortingColumn: "number",
						IsAsc:         false,
					},
				},
				Type: "Launch",
			}

			testDeepEqual(t, f, want)
			return 2, nil, nil
		},
	}

	r := NewReportPortal(&reportportal.Client{
		Filter: mockFilter,
	})

	inputFilter := &Filter{
		Kind:        FilterKind,
		Name:        "mk-e2e-test-suite",
		Description: "",
		Type:        "Launch",
		Conditions: []FilterCondition{
			{
				FilteringField: "name",
				Condition:      "eq",
				Value:          "mk-e2e-test-suite",
			},
		},
		Orders: []FilterOrder{
			{
				SortingColumn: "startTime",
				IsAsc:         false,
			},
			{
				SortingColumn: "number",
				IsAsc:         false,
			},
		},
	}

	err := r.Filter.ApplyFilter("test_project", inputFilter)
	if err != nil {
		t.Errorf("ReportPortal.ApplyFilter returned error: %v", err)
	}

	testDeepEqual(t, mockFilter.Counter, reportportal.MockFilterServiceCounter{GetByName: 1, Create: 1})
}

func TestApplyFilter_Update(t *testing.T) {
	mockFilter := &reportportal.MockFilterService{
		GetByNameM: func(projectName, name string) (*reportportal.Filter, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, name, "mk-e2e-test-suite")
			return &reportportal.Filter{
				Owner: "dbizzarr",
				Share: true,
				ID:    2,
				Name:  "mk-e2e-test-suite",
				Conditions: []reportportal.FilterCondition{
					{
						FilteringField: "name",
						Condition:      "eq",
						Value:          "mk-e2e-test",
					},
				},
				Orders: []reportportal.FilterOrder{
					{
						SortingColumn: "startTime",
						IsAsc:         false,
					},
					{
						SortingColumn: "number",
						IsAsc:         false,
					},
				},
				Type: "Launch",
			}, nil, nil
		},
		UpdateM: func(projectName string, id int, f *reportportal.UpdateFilter) (string, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, id, 2)

			want := &reportportal.UpdateFilter{
				Share: true,
				Name:  "mk-e2e-test-suite",
				Conditions: []reportportal.FilterCondition{
					{
						FilteringField: "name",
						Condition:      "eq",
						Value:          "mk-e2e-test-suite",
					},
				},
				Orders: []reportportal.FilterOrder{
					{
						SortingColumn: "startTime",
						IsAsc:         false,
					},
					{
						SortingColumn: "number",
						IsAsc:         false,
					},
				},
				Type: "Launch",
			}

			testDeepEqual(t, f, want)
			return "", nil, nil
		},
	}

	r := NewReportPortal(&reportportal.Client{
		Filter: mockFilter,
	})

	inputFilter := &Filter{
		Kind:        FilterKind,
		Name:        "mk-e2e-test-suite",
		Description: "",
		Type:        "Launch",
		Conditions: []FilterCondition{
			{
				FilteringField: "name",
				Condition:      "eq",
				Value:          "mk-e2e-test-suite",
			},
		},
		Orders: []FilterOrder{
			{
				SortingColumn: "startTime",
				IsAsc:         false,
			},
			{
				SortingColumn: "number",
				IsAsc:         false,
			},
		},
	}

	err := r.Filter.ApplyFilter("test_project", inputFilter)
	if err != nil {
		t.Errorf("ReportPortal.ApplyFilter returned error: %v", err)
	}

	testDeepEqual(t, mockFilter.Counter, reportportal.MockFilterServiceCounter{GetByName: 1, Update: 1})

}

func TestApplyFilter_Skip(t *testing.T) {
	mockFilter := &reportportal.MockFilterService{
		GetByNameM: func(projectName, name string) (*reportportal.Filter, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, name, "mk-e2e-test-suite")
			return &reportportal.Filter{
				Owner: "dbizzarr",
				Share: true,
				ID:    2,
				Name:  "mk-e2e-test-suite",
				Conditions: []reportportal.FilterCondition{
					{
						FilteringField: "name",
						Condition:      "eq",
						Value:          "mk-e2e-test-suite",
					},
				},
				Orders: []reportportal.FilterOrder{
					{
						SortingColumn: "startTime",
						IsAsc:         false,
					},
					{
						SortingColumn: "number",
						IsAsc:         false,
					},
				},
				Type: "Launch",
			}, nil, nil
		},
	}

	r := NewReportPortal(&reportportal.Client{
		Filter: mockFilter,
	})

	inputFilter := &Filter{
		Kind:        FilterKind,
		Name:        "mk-e2e-test-suite",
		Description: "",
		Type:        "Launch",
		Conditions: []FilterCondition{
			{
				FilteringField: "name",
				Condition:      "eq",
				Value:          "mk-e2e-test-suite",
			},
		},
		Orders: []FilterOrder{
			{
				SortingColumn: "startTime",
				IsAsc:         false,
			},
			{
				SortingColumn: "number",
				IsAsc:         false,
			},
		},
	}

	err := r.Filter.ApplyFilter("test_project", inputFilter)
	if err != nil {
		t.Errorf("ReportPortal.ApplyFilter returned error: %v", err)
	}

	testDeepEqual(t, mockFilter.Counter, reportportal.MockFilterServiceCounter{GetByName: 1})
}

func TestToFilter(t *testing.T) {

	inputFilter := &reportportal.Filter{
		Owner: "dbizzarr",
		Share: true,
		ID:    2,
		Name:  "mk-e2e-test-suite",
		Conditions: []reportportal.FilterCondition{
			{
				FilteringField: "name",
				Condition:      "eq",
				Value:          "mk-e2e-test-suite",
			},
		},
		Orders: []reportportal.FilterOrder{
			{
				SortingColumn: "startTime",
				IsAsc:         false,
			},
			{
				SortingColumn: "number",
				IsAsc:         false,
			},
		},
		Type: "Launch",
	}

	got := ToFilter(inputFilter)

	want := &Filter{
		Kind:        FilterKind,
		Name:        "mk-e2e-test-suite",
		Description: "",
		Type:        "Launch",
		Conditions: []FilterCondition{
			{
				FilteringField: "name",
				Condition:      "eq",
				Value:          "mk-e2e-test-suite",
			},
		},
		Orders: []FilterOrder{
			{
				SortingColumn: "startTime",
				IsAsc:         false,
			},
			{
				SortingColumn: "number",
				IsAsc:         false,
			},
		},
		origin: inputFilter,
	}

	testDeepEqual(t, got, want, cmp.AllowUnexported(Filter{}))
}

func TestFilterToNewFilter(t *testing.T) {

	inputFilter := &Filter{
		Kind:        FilterKind,
		Name:        "mk-e2e-test-suite",
		Description: "",
		Type:        "Launch",
		Conditions: []FilterCondition{
			{
				FilteringField: "name",
				Condition:      "eq",
				Value:          "mk-e2e-test-suite",
			},
		},
		Orders: []FilterOrder{
			{
				SortingColumn: "startTime",
				IsAsc:         false,
			},
			{
				SortingColumn: "number",
				IsAsc:         false,
			},
		},
	}

	got := FilterToNewFilter(inputFilter)

	want := &reportportal.NewFilter{
		Share: true,
		Name:  "mk-e2e-test-suite",
		Conditions: []reportportal.FilterCondition{
			{
				FilteringField: "name",
				Condition:      "eq",
				Value:          "mk-e2e-test-suite",
			},
		},
		Orders: []reportportal.FilterOrder{
			{
				SortingColumn: "startTime",
				IsAsc:         false,
			},
			{
				SortingColumn: "number",
				IsAsc:         false,
			},
		},
		Type: "Launch",
	}

	testDeepEqual(t, got, want)
}

func TestFilterToUpdateFilter(t *testing.T) {

	inputFilter := &Filter{
		Kind:        FilterKind,
		Name:        "mk-e2e-test-suite",
		Description: "",
		Type:        "Launch",
		Conditions: []FilterCondition{
			{
				FilteringField: "name",
				Condition:      "eq",
				Value:          "mk-e2e-test-suite",
			},
		},
		Orders: []FilterOrder{
			{
				SortingColumn: "startTime",
				IsAsc:         false,
			},
			{
				SortingColumn: "number",
				IsAsc:         false,
			},
		},
	}

	got := FilterToUpdateFilter(inputFilter)

	want := &reportportal.UpdateFilter{
		Share: true,
		Name:  "mk-e2e-test-suite",
		Conditions: []reportportal.FilterCondition{
			{
				FilteringField: "name",
				Condition:      "eq",
				Value:          "mk-e2e-test-suite",
			},
		},
		Orders: []reportportal.FilterOrder{
			{
				SortingColumn: "startTime",
				IsAsc:         false,
			},
			{
				SortingColumn: "number",
				IsAsc:         false,
			},
		},
		Type: "Launch",
	}

	testDeepEqual(t, got, want)
}

func TestFilterEquals(t *testing.T) {

	tests := []*struct {
		description string
		left        *Filter
		right       *Filter
		expexct     bool
	}{
		{
			description: "Compare empty filters should return true",
			left:        &Filter{},
			right:       &Filter{},
			expexct:     true,
		},
		{
			description: "Compare equal filters but only one with the origin filed should return true",
			left: &Filter{
				Kind:        FilterKind,
				Name:        "Test",
				Description: "My test description",
				origin:      &reportportal.Filter{ID: 1},
			},
			right: &Filter{
				Kind:        FilterKind,
				Name:        "Test",
				Description: "My test description",
			},
			expexct: true,
		},
		{
			description: "Compare filters with differt names should return false",
			left: &Filter{
				Name:        "Test One",
				Description: "My test description",
			},
			right: &Filter{
				Name:        "Test",
				Description: "My test description",
			},
			expexct: false,
		},
		{
			description: "Compare filters with differt description should return false",
			left: &Filter{
				Name:        "Test",
				Description: "My test description",
			},
			right: &Filter{
				Name:        "Test",
				Description: "My updated test description",
			},
			expexct: false,
		},
		{
			description: "Compare filters where one has nil conditions should return false",
			left: &Filter{
				Conditions: []FilterCondition{},
			},
			right: &Filter{
				Conditions: nil,
			},
			expexct: false,
		},
		{
			description: "Compare filters where one has nil orders should return false",
			left: &Filter{
				Orders: []FilterOrder{},
			},
			right: &Filter{
				Orders: nil,
			},
			expexct: false,
		},
		{
			description: "Compare filters with different condtions should return false",
			left: &Filter{
				Conditions: []FilterCondition{
					{FilteringField: "description", Condition: "equal", Value: "test"},
					{FilteringField: "name", Condition: "equal", Value: "test"},
				},
			},
			right: &Filter{
				Conditions: []FilterCondition{
					{FilteringField: "description", Condition: "equal", Value: "test"},
					{FilteringField: "name", Condition: "equal", Value: "test2"},
				},
			},
			expexct: false,
		},
		{
			description: "Compare filters with different number of condtions should return false",
			left: &Filter{
				Conditions: []FilterCondition{
					{FilteringField: "description", Condition: "equal", Value: "test"},
					{FilteringField: "name", Condition: "equal", Value: "test"},
				},
			},
			right: &Filter{
				Conditions: []FilterCondition{
					{FilteringField: "description", Condition: "equal", Value: "test"},
				},
			},
			expexct: false,
		},
		{
			description: "Compare filters with same conditions and same order should return true",
			left: &Filter{
				Conditions: []FilterCondition{
					{FilteringField: "description", Condition: "equal", Value: "test"},
					{FilteringField: "name", Condition: "equal", Value: "test"},
				},
			},
			right: &Filter{
				Conditions: []FilterCondition{
					{FilteringField: "description", Condition: "equal", Value: "test"},
					{FilteringField: "name", Condition: "equal", Value: "test"},
				},
			},
			expexct: true,
		},
		{
			description: "Compare filters with same conditions should return true",
			left: &Filter{
				Conditions: []FilterCondition{
					{FilteringField: "description", Condition: "equal", Value: "test"},
					{FilteringField: "name", Condition: "equal", Value: "test"},
				},
			},
			right: &Filter{
				Conditions: []FilterCondition{
					{FilteringField: "name", Condition: "equal", Value: "test"},
					{FilteringField: "description", Condition: "equal", Value: "test"},
				},
			},
			expexct: true,
		},
		{
			description: "Compare filters with different orders should return false",
			left: &Filter{
				Orders: []FilterOrder{
					{SortingColumn: "description", IsAsc: true},
					{SortingColumn: "name", IsAsc: false},
				},
			},
			right: &Filter{
				Orders: []FilterOrder{
					{SortingColumn: "description", IsAsc: true},
					{SortingColumn: "name", IsAsc: true},
				},
			},
			expexct: false,
		},
		{
			description: "Compare filters with different number of orders should return false",
			left: &Filter{
				Orders: []FilterOrder{
					{SortingColumn: "description", IsAsc: true},
					{SortingColumn: "name", IsAsc: false},
				},
			},
			right: &Filter{
				Orders: []FilterOrder{
					{SortingColumn: "description", IsAsc: true},
				},
			},
			expexct: false,
		},
		{
			description: "Compare filters with same orders and same order should return true",
			left: &Filter{
				Orders: []FilterOrder{
					{SortingColumn: "description", IsAsc: true},
					{SortingColumn: "name", IsAsc: false},
				},
			},
			right: &Filter{
				Orders: []FilterOrder{
					{SortingColumn: "description", IsAsc: true},
					{SortingColumn: "name", IsAsc: false},
				},
			},
			expexct: true,
		},
		{
			description: "Compare filters with same orders should return true",
			left: &Filter{
				Orders: []FilterOrder{
					{SortingColumn: "description", IsAsc: true},
					{SortingColumn: "name", IsAsc: false},
				},
			},
			right: &Filter{
				Orders: []FilterOrder{
					{SortingColumn: "name", IsAsc: false},
					{SortingColumn: "description", IsAsc: true},
				},
			},
			expexct: true,
		},
		{
			description: "Compare equal filters should return true",
			left: &Filter{
				Kind:        FilterKind,
				Name:        "Test",
				Description: "My test description",
				Conditions: []FilterCondition{
					{FilteringField: "description", Condition: "equal", Value: "test"},
					{FilteringField: "name", Condition: "equal", Value: "test"},
				},
				Orders: []FilterOrder{
					{SortingColumn: "description", IsAsc: true},
					{SortingColumn: "name", IsAsc: false},
				},
			},
			right: &Filter{
				Kind:        FilterKind,
				Name:        "Test",
				Description: "My test description",
				Conditions: []FilterCondition{
					{FilteringField: "description", Condition: "equal", Value: "test"},
					{FilteringField: "name", Condition: "equal", Value: "test"},
				},
				Orders: []FilterOrder{
					{SortingColumn: "name", IsAsc: false},
					{SortingColumn: "description", IsAsc: true},
				},
			},
			expexct: true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			r := test.left.Equals(test.right)
			if r != test.expexct {
				t.Errorf("expected '%t' but got '%t'", test.expexct, r)
			}
		})
	}
}
