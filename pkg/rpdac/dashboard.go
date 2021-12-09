package rpdac

import (
	"fmt"
	"os"

	"github.com/b1zzu/reportportal-dashboards-as-code/pkg/reportportal"
	"gopkg.in/yaml.v2"
)

type Dashboard struct {
	Name    string    `json:"name"`
	Widgets []*Widget `json:"widgets"`
}

type Widget struct {
	Name              string                   `json:"name"`
	Description       string                   `json:"description"`
	WidgetType        string                   `json:"widgetType"`
	WidgetSize        *WidgetSize              `json:"widgetSize"`
	WidgetPosition    *WidgetPosition          `json:"widgetPosition"`
	Filters           []int                    `json:"filters"`
	ContentParameters *WidgetContentParameters `json:"contentParameters"`
}

type WidgetSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type WidgetPosition struct {
	PositionX int `json:"positionX"`
	PositionY int `json:"positionY"`
}

type WidgetContentParameters struct {
	ContentFields []string               `json:"contentFields"`
	ItemsCount    int                    `json:"itemsCount"`
	WidgetOptions map[string]interface{} `json:"widgetOptions"`
}

func NewDashboard(name string, widgets []*Widget) *Dashboard {
	return &Dashboard{Name: name, Widgets: widgets}
}

func NewWidget(name, description, widgetType string, width, height, positionX, positionY int, filters []int, contentFields []string, itemsCount int, widgetOptions map[string]interface{}) *Widget {
	return &Widget{
		Name:              name,
		Description:       description,
		WidgetType:        widgetType,
		WidgetSize:        &WidgetSize{Width: width, Height: height},
		WidgetPosition:    &WidgetPosition{PositionX: positionX, PositionY: positionY},
		Filters:           filters,
		ContentParameters: &WidgetContentParameters{ContentFields: contentFields, ItemsCount: itemsCount, WidgetOptions: widgetOptions},
	}
}

func LoadDashboardFromReportPortal(client *reportportal.Client, project string, dashboardID int) (*Dashboard, error) {

	// retireve the dashboard defintion
	d, _, err := client.Dashboard.Get(project, dashboardID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving dashboard %d from project %s: %w", dashboardID, project, err)
	}

	widgets := make([]*Widget, len(d.Widgets))

	// retrieve all widgets definitions
	for i, dw := range d.Widgets {
		w, _, err := client.Widget.Get(project, dw.WidgetId)
		if err != nil {
			return nil, fmt.Errorf("error retrieving widget %d from project %s: %w", dw.WidgetId, project, err)
		}

		filters := make([]int, len(w.AppliedFilters))
		for j, f := range w.AppliedFilters {
			filters[j] = f.ID
		}

		widgets[i] = NewWidget(
			w.Name,
			w.Description,
			w.WidgetType,
			dw.WidgetSize.Width,
			dw.WidgetSize.Height,
			dw.WidgetPosition.PositionX,
			dw.WidgetPosition.PositionY,
			filters,
			w.ContentParameters.ContentFields,
			w.ContentParameters.ItemsCount,
			w.ContentParameters.WidgetOptions)
	}

	return NewDashboard(d.Name, widgets), nil
}

func (d *Dashboard) ToYaml() ([]byte, error) {
	b, err := yaml.Marshal(d)
	if err != nil {
		return []byte{}, fmt.Errorf("error marshal dashboard %s: %w", d.Name, err)
	}
	return b, nil
}

func (d *Dashboard) WriteToFile(file string) error {

	y, err := d.ToYaml()
	if err != nil {
		return err
	}

	err = os.WriteFile(file, y, 0660)
	if err != nil {
		return fmt.Errorf("error writing yaml dashboard %s to file %s: %w", d.Name, file, err)
	}
	return nil
}
