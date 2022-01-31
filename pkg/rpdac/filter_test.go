package rpdac

import (
	"testing"

	"github.com/b1zzu/reportportal-dashboards-as-code/pkg/reportportal"
	"github.com/google/go-cmp/cmp"
)

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
		Kind:        "Filter",
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
		Kind:        "Filter",
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
		Kind:        "Filter",
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
				Kind:        "Filter",
				Name:        "Test",
				Description: "My test description",
				origin:      &reportportal.Filter{ID: 1},
			},
			right: &Filter{
				Kind:        "Filter",
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
				Kind:        "Filter",
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
				Kind:        "Filter",
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
