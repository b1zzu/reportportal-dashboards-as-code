package reportportal

import (
	"fmt"
	"net/url"
)

type DashboardService service

type DashboardList struct {
	Content []*Dashboard `json:"content"`
}

type Dashboard struct {
	Owner       string            `json:"owner"`
	Share       bool              `json:"share"`
	ID          int               `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Widgets     []DashboardWidget `json:"widgets"`
}

type DashboardWidget struct {
	WidgetID       int                     `json:"widgetId"`
	Share          bool                    `json:"share"`
	WidgetName     string                  `json:"widgetName"`
	WidgetType     string                  `json:"widgetType"`
	WidgetSize     DashboardWidgetSize     `json:"widgetSize"`
	WidgetPosition DashboardWidgetPosition `json:"widgetPosition"`
}

type DashboardWidgetSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type DashboardWidgetPosition struct {
	PositionX int `json:"positionX"`
	PositionY int `json:"positionY"`
}

type NewDashboard struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Share       bool   `json:"share"`
}

type UpdateDashboard struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Share       bool   `json:"share"`
}

type DashboardAddWidget struct {
	AddWidget *DashboardWidget `json:"addWidget"`
}

type DashboardNotFoundError struct {
	Message string
}

func NewDashboardNotFoundError(projectName, dashboardName string) *DashboardNotFoundError {
	return &DashboardNotFoundError{Message: fmt.Sprintf("error dashboard with name \"%s\" in project \"%s\" not found", dashboardName, projectName)}
}

func (e *DashboardNotFoundError) Error() string {
	return e.Message
}

func (s *DashboardService) GetByName(projectName, name string) (*Dashboard, *Response, error) {
	u := fmt.Sprintf("v1/%s/dashboard?%s", projectName, url.Values{"filter.eq.name": []string{name}}.Encode())

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	dl := new(DashboardList)
	resp, err := s.client.Do(req, dl)
	if err != nil {
		return nil, resp, err
	}

	if len(dl.Content) == 0 {
		return nil, resp, NewDashboardNotFoundError(projectName, name)
	}

	return dl.Content[0], resp, nil
}

func (s *DashboardService) GetByID(projectName string, id int) (*Dashboard, *Response, error) {
	u := fmt.Sprintf("v1/%s/dashboard/%d", projectName, id)

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

func (s *DashboardService) Create(projectName string, d *NewDashboard) (int, *Response, error) {
	u := fmt.Sprintf("v1/%v/dashboard", projectName)

	req, err := s.client.NewRequest("POST", u, d)
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

func (s *DashboardService) Update(projectName string, dashboardID int, d *UpdateDashboard) (string, *Response, error) {
	u := fmt.Sprintf("v1/%v/dashboard/%d", projectName, dashboardID)

	req, err := s.client.NewRequest("PUT", u, d)
	if err != nil {
		return "", nil, err
	}

	c := new(OperationCompletion)
	resp, err := s.client.Do(req, c)
	if err != nil {
		return "", resp, err
	}

	return c.Message, resp, nil
}

func (s *DashboardService) Delete(projectName string, id int) (string, *Response, error) {
	u := fmt.Sprintf("v1/%s/dashboard/%d", projectName, id)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return "", nil, err
	}

	c := new(OperationCompletion)
	resp, err := s.client.Do(req, c)
	if err != nil {
		return "", resp, err
	}

	return c.Message, resp, nil
}

func (s *DashboardService) AddWidget(projectName string, dashboardID int, w *DashboardWidget) (string, *Response, error) {
	u := fmt.Sprintf("v1/%v/dashboard/%d/add", projectName, dashboardID)

	req, err := s.client.NewRequest("PUT", u, &DashboardAddWidget{AddWidget: w})
	if err != nil {
		return "", nil, err
	}

	e := new(OperationCompletion)
	resp, err := s.client.Do(req, e)
	if err != nil {
		return "", resp, err
	}

	return e.Message, resp, nil
}

func (s *DashboardService) RemoveWidget(projectName string, dashboardID int, widgetID int) (string, *Response, error) {
	u := fmt.Sprintf("v1/%v/dashboard/%d/%d", projectName, dashboardID, widgetID)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return "", nil, err
	}

	e := new(OperationCompletion)
	resp, err := s.client.Do(req, e)
	if err != nil {
		return "", resp, err
	}

	return e.Message, resp, nil
}
