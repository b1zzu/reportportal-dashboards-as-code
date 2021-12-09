package reportportal

import "fmt"

type FilterService service

type FilterList struct {
	Content []*Filter `json:"content"`
}

type Filter struct {
	Share bool   `json:"share"`
	ID    int    `json:"id"`
	Name  string `json:"name"`
	// incomplete
}

type FilterNotFoundError struct {
	Message string
}

func NewFilterNotFoundError(projectName, filterName string) *DashboardNotFoundError {
	return &DashboardNotFoundError{Message: fmt.Sprintf("error filter with name \"%s\" in project \"%s\" not found", filterName, projectName)}
}

func (e *FilterNotFoundError) Error() string {
	return e.Message
}

func (s *FilterService) GetByName(projectName, name string) (*Filter, *Response, error) {
	u := fmt.Sprintf("v1/%s/filter?filter.eq.name=%s", projectName, name)

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
