package reportportal

import (
	"fmt"
)

type IWidgetService interface {
	Get(projectName string, id int) (*Widget, *Response, error)
	Post(projectName string, w *NewWidget) (int, *Response, error)
}

type WidgetService service

type Widget struct {
	Description       string                  `json:"description"`
	Owner             string                  `json:"owner"`
	Share             bool                    `json:"share"`
	ID                int                     `json:"id"`
	Name              string                  `json:"name"`
	WidgetType        string                  `json:"widgetType"`
	ContentParameters WidgetContentParameters `json:"contentParameters"`
	AppliedFilters    []Filter                `json:"appliedFilters"`
	Content           interface{}             `json:"content"` // incomplete
}

type WidgetContentParameters struct {
	ContentFields []string               `json:"contentFields"`
	ItemsCount    int                    `json:"itemsCount"`
	WidgetOptions map[string]interface{} `json:"widgetOptions"`
}

type NewWidget struct {
	Name              string                  `json:"name"`
	Description       string                  `json:"description"`
	Share             bool                    `json:"share"`
	WidgetType        string                  `json:"widgetType"`
	ContentParameters WidgetContentParameters `json:"contentParameters"`
	Filters           []int                   `json:"filterIds"`
}

func (s *WidgetService) Get(projectName string, id int) (*Widget, *Response, error) {
	u := fmt.Sprintf("v1/%v/widget/%v", projectName, id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	w := new(Widget)
	resp, err := s.client.Do(req, w)
	if err != nil {
		return nil, resp, err
	}

	return w, resp, nil
}

func (s *WidgetService) Post(projectName string, w *NewWidget) (int, *Response, error) {
	u := fmt.Sprintf("v1/%s/widget", projectName)

	req, err := s.client.NewRequest("POST", u, w)
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
