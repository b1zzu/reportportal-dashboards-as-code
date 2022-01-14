package rpdac

import "github.com/b1zzu/reportportal-dashboards-as-code/pkg/reportportal"

type Filter struct {
	Name        string             `json:"name"`
	Kind        string             `json:"kind"`
	Type        string             `json:"type"`
	Description string             `json:"description"`
	Conditions  []*FilterCondition `json:"conditions"`
	Orders      []*FilterOrder     `json:"orders"`
}

type FilterCondition struct {
	FilteringField string `json:"filteringField"`
	Condition      string `json:"condition"`
	Value          string `json:"value"`
}

type FilterOrder struct {
	SortingColumn string `json:"sortingColumn"`
	IsAsc         bool   `json:"isAsc"`
}

func ToFilter(f *reportportal.Filter) *Filter {

	conditions := make([]*FilterCondition, len(f.Conditions))
	for i, c := range f.Conditions {
		conditions[i] = &FilterCondition{Condition: c.Condition, FilteringField: c.FilteringField, Value: c.Value}
	}

	orders := make([]*FilterOrder, len(f.Orders))
	for i, o := range f.Orders {
		orders[i] = &FilterOrder{IsAsc: o.IsAsc, SortingColumn: o.SortingColumn}
	}

	return &Filter{Name: f.Name, Kind: "Filter", Type: f.Type, Description: f.Description, Conditions: conditions, Orders: orders}
}

func (f *Filter) GetName() string {
	return f.Name
}

func (f *Filter) GetKind() string {
	return f.Kind
}
