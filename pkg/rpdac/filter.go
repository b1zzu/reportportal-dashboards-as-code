package rpdac

import (
	"github.com/b1zzu/reportportal-dashboards-as-code/pkg/reportportal"
	"github.com/b1zzu/reportportal-dashboards-as-code/pkg/util"
	"gopkg.in/yaml.v2"
)

const FilterKind = "Filter"

type Filter struct {
	Kind        string            `json:"kind"`
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Description string            `json:"description"`
	Conditions  []FilterCondition `json:"conditions"`
	Orders      []FilterOrder     `json:"orders"`

	origin *reportportal.Filter
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

	conditions := make([]FilterCondition, len(f.Conditions))
	for i, c := range f.Conditions {
		conditions[i] = FilterCondition{Condition: c.Condition, FilteringField: c.FilteringField, Value: c.Value}
	}

	orders := make([]FilterOrder, len(f.Orders))
	for i, o := range f.Orders {
		orders[i] = FilterOrder{IsAsc: o.IsAsc, SortingColumn: o.SortingColumn}
	}

	return &Filter{
		Name:        f.Name,
		Kind:        FilterKind,
		Type:        f.Type,
		Description: f.Description,
		Conditions:  conditions,
		Orders:      orders,
		origin:      f,
	}
}

func toFilterConditions(conditions []FilterCondition) []reportportal.FilterCondition {

	r := make([]reportportal.FilterCondition, len(conditions))
	for i, c := range conditions {
		r[i] = reportportal.FilterCondition{Condition: c.Condition, FilteringField: c.FilteringField, Value: c.Value}
	}

	return r
}

func toFilterOrders(orders []FilterOrder) []reportportal.FilterOrder {

	r := make([]reportportal.FilterOrder, len(orders))
	for i, o := range orders {
		r[i] = reportportal.FilterOrder{IsAsc: o.IsAsc, SortingColumn: o.SortingColumn}
	}

	return r
}

func FilterToNewFilter(f *Filter) *reportportal.NewFilter {

	return &reportportal.NewFilter{
		Name:        f.Name,
		Type:        f.Type,
		Description: f.Description,
		Share:       true,
		Conditions:  toFilterConditions(f.Conditions),
		Orders:      toFilterOrders(f.Orders),
	}
}

func FilterToUpdateFilter(f *Filter) *reportportal.UpdateFilter {

	return &reportportal.UpdateFilter{
		Name:        f.Name,
		Type:        f.Type,
		Description: f.Description,
		Share:       true,
		Conditions:  toFilterConditions(f.Conditions),
		Orders:      toFilterOrders(f.Orders),
	}
}

func LoadFilterFromFile(file []byte) (*Filter, error) {

	f := new(Filter)
	err := yaml.Unmarshal(file, f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (f *Filter) GetName() string {
	return f.Name
}

func (f *Filter) GetKind() string {
	return f.Kind
}

func conditionsToStringSlice(conditions []FilterCondition) []string {
	s := make([]string, len(conditions))
	for i, c := range conditions {
		s[i] = c.Condition + c.FilteringField + c.Value
	}
	return s
}

func ordersToStringSlice(orders []FilterOrder) []string {
	s := make([]string, len(orders))
	for i, o := range orders {
		sort := "dsc"
		if o.IsAsc {
			sort = "asc"
		}

		s[i] = o.SortingColumn + sort
	}
	return s
}

func (left *Filter) Equals(right *Filter) bool {
	if left == nil || right == nil {
		return left == right
	}

	if left.Name != right.Name ||
		left.Type != right.Type ||
		left.Description != right.Description {

		return false
	}

	// compare conditions
	if left.Conditions != nil && right.Conditions != nil {

		leftConditions := conditionsToStringSlice(left.Conditions)
		rightCondition := conditionsToStringSlice(right.Conditions)

		if !util.CompareStringSlices(leftConditions, rightCondition) {
			return false
		}

	} else if !(left.Conditions == nil && right.Conditions == nil) {
		return false
	}

	// compare orders
	if left.Orders == nil || right.Orders == nil {
		return left.Orders == nil && right.Orders == nil
	}

	leftOrders := ordersToStringSlice(left.Orders)
	rightOrders := ordersToStringSlice(right.Orders)

	return util.CompareStringSlices(leftOrders, rightOrders)
}
