package reportportal

import (
	"fmt"
	"net/url"
)

type FilterService service

type FilterList struct {
	Content []*Filter `json:"content"`
}

type Filter struct {
	Share       bool               `json:"share"`
	ID          int                `json:"id"`
	Name        string             `json:"name"`
	Type        string             `json:"type"`
	Description string             `json:"description"`
	Owner       string             `json:"onwer"`
	Conditions  []*FilterCondition `json:"conditions"`
	Orders      []*FilterOrder     `json:"orders"`
}

type FilterCondition struct {
	Condition      string `json:"condition"`
	FilteringField string `json:"filteringField"`
	Value          string `json:"value"`
}

type FilterOrder struct {
	IsAsc         bool   `json:"isAsc"`
	SortingColumn string `json:"sortingColumn"`
}

type FilterNotFoundError struct {
	Message string
}

type NewFilter struct {
	Name        string             `json:"name"`
	Type        string             `json:"type"`
	Description string             `json:"description"`
	Share       bool               `json:"share"`
	Conditions  []*FilterCondition `json:"conditions"`
	Orders      []*FilterOrder     `json:"orders"`
}

func NewFilterNotFoundError(projectName, filterName string) *DashboardNotFoundError {
	return &DashboardNotFoundError{Message: fmt.Sprintf("error filter with name \"%s\" in project \"%s\" not found", filterName, projectName)}
}

func (e *FilterNotFoundError) Error() string {
	return e.Message
}

func (s *FilterService) GetByID(projectName string, id int) (*Filter, *Response, error) {
	u := fmt.Sprintf("v1/%s/filter/%d", projectName, id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	f := new(Filter)
	resp, err := s.client.Do(req, f)
	if err != nil {
		return nil, resp, err
	}

	return f, resp, nil
}

func (s *FilterService) GetByName(projectName, name string) (*Filter, *Response, error) {
	u := fmt.Sprintf("v1/%s/filter?%s", projectName, url.Values{"filter.eq.name": []string{name}}.Encode())

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	fl := new(FilterList)
	resp, err := s.client.Do(req, fl)
	if err != nil {
		return nil, resp, err
	}

	if len(fl.Content) == 0 {
		return nil, resp, NewFilterNotFoundError(projectName, name)
	}

	return fl.Content[0], resp, nil
}

func (s *FilterService) Create(projectName string, f *NewFilter) (int, *Response, error) {
	u := fmt.Sprintf("v1/%v/filter", projectName)

	req, err := s.client.NewRequest("POST", u, f)
	if err != nil {
		return 0, nil, err
	}

	e := new(EntryCreated)
	resp, err := s.client.Do(req, e)
	if err != nil {
		return 0, resp, err
	}

	return e.ID, resp, nil
}
