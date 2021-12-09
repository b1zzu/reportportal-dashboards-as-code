package rpdac

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"

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
	Filters           []string                 `json:"filters"`
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

func ToDashboard(d *reportportal.Dashboard, widgets []*Widget) *Dashboard {
	return &Dashboard{Name: d.Name, Widgets: widgets}
}

func ToWidget(w *reportportal.Widget, dw *reportportal.DashboardWidget) *Widget {

	filters := make([]string, len(w.AppliedFilters))
	for j, f := range w.AppliedFilters {
		filters[j] = f.Name
	}

	return &Widget{
		Name:              w.Name,
		Description:       w.Description,
		WidgetType:        w.WidgetType,
		WidgetSize:        &WidgetSize{Width: dw.WidgetSize.Width, Height: dw.WidgetSize.Height},
		WidgetPosition:    &WidgetPosition{PositionX: dw.WidgetPosition.PositionX, PositionY: dw.WidgetPosition.PositionY},
		Filters:           filters,
		ContentParameters: &WidgetContentParameters{ContentFields: w.ContentParameters.ContentFields, ItemsCount: w.ContentParameters.ItemsCount, WidgetOptions: w.ContentParameters.WidgetOptions},
	}
}

func FromWidget(dashboardHash string, w *Widget, widgetFilters []int) (*reportportal.NewWidget, *reportportal.DashboardWidget) {

	nw := &reportportal.NewWidget{
		// For the rpdac tool the widget name is not unique across all dashboards, while fore ReportPortal it is,
		// by adding the dashboard name sha to the widget name we make it generic
		Name:        fmt.Sprintf("%s #%s", w.Name, dashboardHash),
		Description: w.Description,
		Share:       true,
		WidgetType:  w.WidgetType,
		Filters:     widgetFilters,
		ContentParameters: &reportportal.WidgetContentParameters{
			ItemsCount:    w.ContentParameters.ItemsCount,
			ContentFields: w.ContentParameters.ContentFields,
			WidgetOptions: w.ContentParameters.WidgetOptions,
		},
	}

	dw := &reportportal.DashboardWidget{
		Share:          true,
		WidgetName:     w.Name,
		WidgetType:     w.WidgetType,
		WidgetSize:     &reportportal.DashboardWidgetSize{Width: w.WidgetSize.Width, Height: w.WidgetSize.Height},
		WidgetPosition: &reportportal.DashboardWidgetPosition{PositionX: w.WidgetPosition.PositionX, PositionY: w.WidgetPosition.PositionY},
	}

	return nw, dw
}

func LoadDashboardFromFile(file string) (*Dashboard, error) {

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	d := new(Dashboard)
	err = yaml.Unmarshal(b, d)
	if err != nil {
		return nil, err
	}

	return d, nil
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

	err = ioutil.WriteFile(file, y, 0660)
	if err != nil {
		return fmt.Errorf("error writing yaml dashboard %s to file %s: %w", d.Name, file, err)
	}
	return nil
}

func (d *Dashboard) HashName() string {
	h := sha1.New()
	io.WriteString(h, d.Name)
	return hex.EncodeToString(h.Sum(nil))[:4]
}
