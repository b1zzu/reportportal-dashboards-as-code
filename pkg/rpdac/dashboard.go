package rpdac

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/b1zzu/reportportal-dashboards-as-code/pkg/reportportal"
	"gopkg.in/yaml.v2"
)

type Dashboard struct {
	Name    string    `json:"name"`
	Kind    string    `json:"kind"`
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
	return &Dashboard{Name: d.Name, Kind: "Dashboard", Widgets: widgets}
}

// convert 'statistics$defects$system_issue$xx_xxxxxxxxxxx' fields to 'statistics$defects$system_issue$shortname`
func DecodeFieldsSubTypes(fields []string, decodeMap map[string]string) ([]string, error) {

	result := make([]string, len(fields))
	for j, f := range fields {
		p := strings.Split(f, "$")
		if p[0] == "statistics" && p[1] == "defects" {
			s, ok := decodeMap[p[3]]
			if !ok {
				return nil, fmt.Errorf("error finding a map for the field \"%s\"", f)
			}
			p[3] = s // replace the locator with the short name
			result[j] = strings.Join(p, "$")
		} else {
			// keep it like it is
			result[j] = f
		}
	}
	return result, nil
}

// convert 'statistics$defects$system_issue$shortname' fields to 'statistics$defects$system_issue$xx_xxxxxxxxxxx'
func EncodeFieldsSubTypes(fields []string, encodeMap map[string]string) ([]string, error) {

	// because the encodeMap is the inverse of the decodeMap we can use the same
	// function but with the inverted map to encode the fields
	return DecodeFieldsSubTypes(fields, encodeMap)
}

func ToWidget(w *reportportal.Widget, dw *reportportal.DashboardWidget, decodeSubTypesMap map[string]string) (*Widget, error) {

	filters := make([]string, len(w.AppliedFilters))
	for j, f := range w.AppliedFilters {
		filters[j] = f.Name
	}

	fields, err := DecodeFieldsSubTypes(w.ContentParameters.ContentFields, decodeSubTypesMap)
	if err != nil {
		return nil, fmt.Errorf("error decoding sub types in widget \"%s\": %w", w.Name, err)
	}

	return &Widget{
		Name:              w.Name,
		Description:       w.Description,
		WidgetType:        w.WidgetType,
		WidgetSize:        &WidgetSize{Width: dw.WidgetSize.Width, Height: dw.WidgetSize.Height},
		WidgetPosition:    &WidgetPosition{PositionX: dw.WidgetPosition.PositionX, PositionY: dw.WidgetPosition.PositionY},
		Filters:           filters,
		ContentParameters: &WidgetContentParameters{ContentFields: fields, ItemsCount: w.ContentParameters.ItemsCount, WidgetOptions: w.ContentParameters.WidgetOptions},
	}, nil
}

func FromWidget(dashboardHash string, w *Widget, filtersMap map[string]int, encodeSubTypesMap map[string]string) (*reportportal.NewWidget, *reportportal.DashboardWidget, error) {

	filters := make([]int, len(w.Filters))
	for j, f := range w.Filters {
		filters[j] = filtersMap[f]
	}

	fields, err := EncodeFieldsSubTypes(w.ContentParameters.ContentFields, encodeSubTypesMap)
	if err != nil {
		return nil, nil, fmt.Errorf("error encoding sub types in widget \"%s\": %w", w.Name, err)
	}

	nw := &reportportal.NewWidget{
		// For the rpdac tool the widget name is not unique across all dashboards, while fore ReportPortal it is,
		// by adding the dashboard name sha to the widget name we make it generic
		Name:        fmt.Sprintf("%s #%s", w.Name, dashboardHash),
		Description: w.Description,
		Share:       true,
		WidgetType:  w.WidgetType,
		Filters:     filters,
		ContentParameters: &reportportal.WidgetContentParameters{
			ItemsCount:    w.ContentParameters.ItemsCount,
			ContentFields: fields,
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

	return nw, dw, nil
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

	if d.Kind != "" && d.Kind != "Dashboard" {
		return nil, fmt.Errorf("error invalid kind '%s' in file '%s'", d.Kind, file)
	}

	if d.Kind == "" {
		log.Printf("warning: assuming kind 'Dashboard' for file '%s'", file)
		d.Kind = "Dashboard"
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
