package rpdac

import (
	"errors"
	"fmt"
	"log"
	"sort"

	"github.com/b1zzu/reportportal-dashboards-as-code/pkg/reportportal"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gopkg.in/yaml.v2"
)

type FilterService service

type Filter struct {
	Kind        ObjectKind        `json:"kind"`
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

func (s *FilterService) Get(project string, id int) (Object, error) {

	// retireve the filter defintion
	f, _, err := s.client.Filter.GetByID(project, id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving filter %d from project %s: %w", id, project, err)
	}

	return ToFilter(f), nil
}

func (s *FilterService) GetByName(project, name string) (Object, error) {

	f, _, err := s.client.Filter.GetByName(project, name)
	if err != nil {
		if _, ok := err.(*reportportal.FilterNotFoundError); ok {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return ToFilter(f), nil
}

func (s *FilterService) Create(project string, o Object) error {
	f := o.(*Filter)

	filterID, _, err := s.client.Filter.Create(project, FilterToNewFilter(f))
	if err != nil {
		return fmt.Errorf("error creating filter %s: %w", f.Name, err)
	}

	log.Printf("filter %s created with id: %d", f.Name, filterID)
	return nil
}

func (s *FilterService) Update(project string, current, target Object) error {
	currentFilter, targetFilter := current.(*Filter), target.(*Filter)

	_, _, err := s.client.Filter.Update(project, currentFilter.origin.ID, FilterToUpdateFilter(targetFilter))
	if err != nil {
		return fmt.Errorf("error updating filter \"%s\": %w", targetFilter.Name, err)
	}
	return nil
}

func (s *FilterService) Delete(project, name string) error {
	return errors.New("unimplemented")
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

func (f *Filter) GetKind() ObjectKind {
	return f.Kind
}

func (left *Filter) Equals(right Object) bool {
	opts := cmp.Options{
		cmpopts.IgnoreUnexported(Filter{}),

		// sort FilterConditions
		cmp.Transformer("SortConditions", func(in []FilterCondition) []FilterCondition {
			out := make([]FilterCondition, len(in))
			copy(out, in) // copy input to avoid mutating it
			sort.Slice(out, func(i, j int) bool {
				return fmt.Sprintf("%+v", out[i]) < fmt.Sprintf("%+v", out[j])
			})
			return out
		}),

		// sort FilterOrders
		cmp.Transformer("SortOrders", func(in []FilterOrder) []FilterOrder {
			out := make([]FilterOrder, len(in))
			copy(out, in) // copy input to avoid mutating it
			sort.Slice(out, func(i, j int) bool {
				return fmt.Sprintf("%+v", out[i]) < fmt.Sprintf("%+v", out[j])
			})
			return out
		}),
	}
	return cmp.Equal(left, right, opts)
}
