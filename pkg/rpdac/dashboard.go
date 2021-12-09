package rpdac

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"

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

func LoadDashboardFromFile(file string) (*Dashboard, error) {

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var d Dashboard
	err = yaml.Unmarshal(b, &d)
	if err != nil {
		return nil, err
	}

	return &d, nil
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
	return string(h.Sum(nil)[:4])
}
