package reportportal

import (
	"fmt"
)

type DashboardService service

type Dashboard struct {
	Owner   string             `json:"owner"`
	Share   bool               `json:"share"`
	ID      int                `json:"id"`
	Name    string             `json:"name"`
	Widgets []*DashboardWidget `json:"widgets"`
}

type DashboardWidget struct {
	WidgetName     string                   `json:"widgetName"`
	WidgetId       int                      `json:"widgetId"`
	WidgetType     string                   `json:"widgetType"`
	WidgetSize     *DashboardWidgetSize     `json:"widgetSize"`
	WidgetPosition *DashboardWidgetPosition `json:"widgetPosition"`
	Share          bool                     `json:"share"`
}

type DashboardWidgetSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type DashboardWidgetPosition struct {
	PositionX int `json:"positionX"`
	PositionY int `json:"positionY"`
}

func (s *DashboardService) Get(projectName string, id int) (*Dashboard, *Response, error) {
	u := fmt.Sprintf("v1/%v/dashboard/%v", projectName, id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	d := new(Dashboard)
	resp, err := s.client.Do(req, d)
	if err != nil {
		return nil, resp, err
	}

	return d, resp, nil
}
