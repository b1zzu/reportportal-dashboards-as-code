package rpdac

import (
	"fmt"
	"log"
	"sort"

	"github.com/b1zzu/reportportal-dashboards-as-code/pkg/reportportal"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gopkg.in/yaml.v2"
)

type IFilterService interface {
	Get(project string, id int) (Object, error)
	Create(project string, o Object) error
	Apply(project string, o Object) error

	GetFilter(project string, id int) (*Filter, error)
	GetFilterByName(project, name string) (*Filter, error)
	CreateFilter(project string, f *Filter) error
	ApplyFilter(project string, f *Filter) error
}

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
	return s.GetFilter(project, id)
}

func (s *FilterService) Create(project string, o Object) error {
	return s.CreateFilter(project, o.(*Filter))
}

func (s *FilterService) Apply(project string, o Object) error {
	return s.ApplyFilter(project, o.(*Filter))
}

func (s *FilterService) GetFilter(project string, id int) (*Filter, error) {

	// retireve the filter defintion
	f, _, err := s.client.Filter.GetByID(project, id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving filter %d from project %s: %w", id, project, err)
	}

	return ToFilter(f), nil
}

func (s *FilterService) GetFilterByName(project, name string) (*Filter, error) {

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

func (s *FilterService) CreateFilter(project string, f *Filter) error {

	filterID, _, err := s.client.Filter.Create(project, FilterToNewFilter(f))
	if err != nil {
		return fmt.Errorf("error creating filter %s: %w", f.Name, err)
	}

	log.Printf("filter %s created with id: %d", f.Name, filterID)
	return nil
}

func (s *FilterService) ApplyFilter(project string, f *Filter) error {

	currentFilter, err := s.GetFilterByName(project, f.Name)
	if err != nil {
		return fmt.Errorf("error retrieving filter \"%s\" by name: %w", f.Name, err)
	}

	if currentFilter != nil {

		if currentFilter.Equals(f) {
			log.Printf("skip apply for filter \"%s\"", f.Name)
			return nil
		}

		return s.updateFilter(project, currentFilter, f)
	}

	return s.CreateFilter(project, f)

}

func (s *FilterService) updateFilter(project string, currentFilter, targetFilter *Filter) error {

	_, _, err := s.client.Filter.Update(project, currentFilter.origin.ID, FilterToUpdateFilter(targetFilter))
	if err != nil {
		return fmt.Errorf("error updating filter \"%s\": %w", targetFilter.Name, err)
	}

	log.Printf("update \"%s\" filter", targetFilter.Name)
	return nil
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

func (left *Filter) Equals(right *Filter) bool {
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
