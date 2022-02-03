package rpdac

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"

	"github.com/b1zzu/reportportal-dashboards-as-code/pkg/reportportal"
	"github.com/b1zzu/reportportal-dashboards-as-code/pkg/util"
	"gopkg.in/yaml.v2"
)

type IDashboardService interface {
	Get(project string, id int) (Object, error)

	GetDashboard(project string, id int) (*Dashboard, error)
	GetDashboardByName(project, name string) (*Dashboard, error)
	CreateDashboard(project string, d *Dashboard) error
	ApplyDashboard(project string, d *Dashboard) error
	DeleteDashboard(project, name string) error
}

type DashboardService service

type Dashboard struct {
	Kind        ObjectKind `json:"kind"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Widgets     []*Widget  `json:"widgets"`

	origin *reportportal.Dashboard
}

type Widget struct {
	Name              string                  `json:"name"`
	Description       string                  `json:"description"`
	WidgetType        string                  `json:"widgetType"`
	WidgetSize        WidgetSize              `json:"widgetSize"`
	WidgetPosition    WidgetPosition          `json:"widgetPosition"`
	Filters           []string                `json:"filters"`
	ContentParameters WidgetContentParameters `json:"contentParameters"`

	origin *reportportal.Widget
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

func (s *DashboardService) Get(project string, id int) (Object, error) {
	return s.GetDashboard(project, id)
}

func (s *DashboardService) GetDashboard(project string, id int) (*Dashboard, error) {

	// retireve the dashboard defintion
	d, _, err := s.client.Dashboard.GetByID(project, id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving dashboard %d from project %s: %w", id, project, err)
	}

	return s.loadDashboard(project, d)
}

func (s *DashboardService) GetDashboardByName(project, name string) (*Dashboard, error) {

	d, _, err := s.client.Dashboard.GetByName(project, name)
	if err != nil {
		if _, ok := err.(*reportportal.DashboardNotFoundError); ok {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return s.loadDashboard(project, d)
}

func (s *DashboardService) loadDashboard(project string, d *reportportal.Dashboard) (*Dashboard, error) {

	widgets := make([]*Widget, len(d.Widgets))

	decodeSubTypesMap, err := s.decodeSubTypseMap(project)
	if err != nil {
		return nil, err
	}

	dashboardHash := HashName(d.Name)

	// retrieve all widgets definitions
	for i, dw := range d.Widgets {
		w, _, err := s.client.Widget.Get(project, dw.WidgetID)
		if err != nil {
			return nil, fmt.Errorf("error retrieving widget %d from project %s: %w", dw.WidgetID, project, err)
		}

		widgets[i], err = ToWidget(w, &dw, dashboardHash, decodeSubTypesMap)
		if err != nil {
			return nil, err
		}
	}

	return ToDashboard(d, widgets), nil
}

func (s *DashboardService) CreateDashboard(project string, d *Dashboard) error {

	// resolve all filters
	filtersMap, err := s.filtersMap(project, d.Widgets)
	if err != nil {
		return err
	}

	encodeSubTypesMap, err := s.encodeSubTypseMap(project)
	if err != nil {
		return err
	}

	dashboardID, _, err := s.client.Dashboard.Create(project, &reportportal.NewDashboard{Name: d.Name, Description: d.Description, Share: true})
	if err != nil {
		return fmt.Errorf("error creating dashboard %s: %w", d.Name, err)
	}
	log.Printf("dashboard %s created with id: %d", d.Name, dashboardID)

	return s.createWidgets(project, dashboardID, d, filtersMap, encodeSubTypesMap)
}

func (s *DashboardService) createWidgets(
	project string,
	dashboardID int,
	dashboard *Dashboard,
	filtersMap map[string]int,
	encodeSubTypesMap map[string]string) error {

	dashboardHash := dashboard.HashName()

	for _, w := range dashboard.Widgets {

		nw, dw, err := FromWidget(dashboardHash, w, filtersMap, encodeSubTypesMap)
		if err != nil {
			return fmt.Errorf("error converting widget %s: %w", w.Name, err)
		}

		widgetID, _, err := s.client.Widget.Post(project, nw)
		if err != nil {
			return fmt.Errorf("error creating widget %s: %w", w.Name, err)
		}

		dw.WidgetID = widgetID

		_, _, err = s.client.Dashboard.AddWidget(project, dashboardID, dw)
		if err != nil {
			return fmt.Errorf("error adding widget %s to dashboard %s: %w", w.Name, dashboard.Name, err)
		}
		log.Printf("added \"%s\" widget to \"%s\" dashboard", w.Name, dashboard.Name)
	}
	return nil
}

// Create or Recreate the dashboard
func (s *DashboardService) ApplyDashboard(project string, d *Dashboard) error {

	currentDashboard, err := s.GetDashboardByName(project, d.Name)
	if err != nil {
		return fmt.Errorf("error retrieving dashboard \"%s\" by name: %w", d.Name, err)
	}

	if currentDashboard != nil {

		if currentDashboard.Equals(d) {
			log.Printf("skip apply for dashboard \"%s\"", d.Name)
			return nil
		}

		return s.updateDashboard(project, currentDashboard, d)
	}

	return s.CreateDashboard(project, d)
}

func (s *DashboardService) updateDashboard(project string, currentDashboard, targetDashboard *Dashboard) error {

	// resolve all filters
	filtersMap, err := s.filtersMap(project, targetDashboard.Widgets)
	if err != nil {
		return err
	}

	encodeSubTypesMap, err := s.encodeSubTypseMap(project)
	if err != nil {
		return err
	}

	dashboardID := currentDashboard.origin.ID

	// delete all widgets from the current dashboard so we can recreate them as expected by the target dashboard
	for _, w := range currentDashboard.Widgets {
		_, _, err := s.client.Dashboard.RemoveWidget(project, dashboardID, w.origin.ID)
		if err != nil {
			return fmt.Errorf("error removing widget \"%s\" from dashboard \"%s\": %w", w.Name, currentDashboard.Name, err)
		}
	}

	u := &reportportal.UpdateDashboard{Name: targetDashboard.Name, Description: targetDashboard.Description, Share: true}
	_, _, err = s.client.Dashboard.Update(project, dashboardID, u)
	if err != nil {
		return fmt.Errorf("error updating dashboard %s: %w", targetDashboard.Name, err)
	}
	log.Printf("updated \"%s\" dashboard", targetDashboard.Name)

	return s.createWidgets(project, dashboardID, targetDashboard, filtersMap, encodeSubTypesMap)
}

// Delete the Dashboard with the given name and Widgets created for it
func (s *DashboardService) DeleteDashboard(project, name string) error {

	d, _, err := s.client.Dashboard.GetByName(project, name)
	if err != nil {
		if _, ok := err.(*reportportal.DashboardNotFoundError); ok {
			// ignore
		} else {
			return err
		}
	}

	// because we have ignored the error in case of DashboardNotFoundError d can also be nil
	if d != nil {
		_, _, err = s.client.Dashboard.Delete(project, d.ID)
		if err != nil {
			return err
		}
	}

	// Widgets are deleted automatically if not used buy any dashboard
	return nil
}

func (s *DashboardService) decodeSubTypseMap(project string) (map[string]string, error) {
	ps, _, err := s.client.ProjectSettings.Get(project)
	if err != nil {
		return nil, err
	}

	decodeMap := make(map[string]string)
	for _, g := range ps.SubTypes {
		for _, t := range g {
			decodeMap[t.Locator] = t.ShortName
		}
	}
	return decodeMap, nil
}

func (s *DashboardService) encodeSubTypseMap(project string) (map[string]string, error) {

	m, err := s.decodeSubTypseMap(project)
	if err != nil {
		return nil, err
	}

	// the decode map is the inverse of the encode
	encodeMap := make(map[string]string)
	for k, v := range m {
		encodeMap[v] = k
	}
	return encodeMap, nil
}

func (s *DashboardService) filtersMap(project string, widgets []*Widget) (map[string]int, error) {
	filtersMap := make(map[string]int)
	for _, w := range widgets {
		for _, filterName := range w.Filters {

			if _, ok := filtersMap[filterName]; ok {
				// filter already resolved
				continue
			}

			f, _, err := s.client.Filter.GetByName(project, filterName)
			if err != nil {
				return nil, fmt.Errorf("error resolving filter \"%s\" in widget \"%s\": %w", filterName, w.Name, err)
			}

			filtersMap[filterName] = f.ID
		}
	}
	return filtersMap, nil
}

func ToDashboard(d *reportportal.Dashboard, widgets []*Widget) *Dashboard {

	return &Dashboard{
		Kind:        DashboardKind,
		Name:        d.Name,
		Description: d.Description,
		Widgets:     widgets,
		origin:      d,
	}
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

func ToWidget(w *reportportal.Widget, dw *reportportal.DashboardWidget, dashboardHash string, decodeSubTypesMap map[string]string) (*Widget, error) {

	name := strings.TrimSuffix(w.Name, fmt.Sprintf(" #%s", dashboardHash))

	filters := make([]string, len(w.AppliedFilters))
	for j, f := range w.AppliedFilters {
		filters[j] = f.Name
	}

	fields, err := DecodeFieldsSubTypes(w.ContentParameters.ContentFields, decodeSubTypesMap)
	if err != nil {
		return nil, fmt.Errorf("error decoding sub types in widget \"%s\": %w", w.Name, err)
	}

	return &Widget{
		Name:              name,
		Description:       w.Description,
		WidgetType:        w.WidgetType,
		WidgetSize:        WidgetSize{Width: dw.WidgetSize.Width, Height: dw.WidgetSize.Height},
		WidgetPosition:    WidgetPosition{PositionX: dw.WidgetPosition.PositionX, PositionY: dw.WidgetPosition.PositionY},
		Filters:           filters,
		ContentParameters: WidgetContentParameters{ContentFields: fields, ItemsCount: w.ContentParameters.ItemsCount, WidgetOptions: w.ContentParameters.WidgetOptions},
		origin:            w,
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
		ContentParameters: reportportal.WidgetContentParameters{
			ItemsCount:    w.ContentParameters.ItemsCount,
			ContentFields: fields,
			WidgetOptions: w.ContentParameters.WidgetOptions,
		},
	}

	dw := &reportportal.DashboardWidget{
		Share:          true,
		WidgetName:     w.Name,
		WidgetType:     w.WidgetType,
		WidgetSize:     reportportal.DashboardWidgetSize{Width: w.WidgetSize.Width, Height: w.WidgetSize.Height},
		WidgetPosition: reportportal.DashboardWidgetPosition{PositionX: w.WidgetPosition.PositionX, PositionY: w.WidgetPosition.PositionY},
	}

	return nw, dw, nil
}

func LoadDashboardFromFile(file []byte) (*Dashboard, error) {

	d := new(Dashboard)
	err := yaml.Unmarshal(file, d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Dashboard) GetName() string {
	return d.Name
}

func (d *Dashboard) GetKind() ObjectKind {
	return d.Kind
}

func (left *Dashboard) Equals(o util.Comparable) bool {
	if left == nil || o == nil {
		return left == o
	}

	right, ok := o.(*Dashboard)
	if !ok {
		return false
	}

	if left.Name != right.Name || left.Description != right.Description {
		return false
	}

	if left.Widgets == nil || right.Widgets == nil {
		return left.Widgets == nil && right.Widgets == nil
	}

	leftWidgets := make([]util.Comparable, len(left.Widgets))
	for i := range left.Widgets {
		leftWidgets[i] = left.Widgets[i]
	}

	rightWidgets := make([]util.Comparable, len(right.Widgets))
	for i := range right.Widgets {
		rightWidgets[i] = right.Widgets[i]
	}

	return util.CompareSlices(leftWidgets, rightWidgets)
}

func (d *Dashboard) Key() string {
	return d.Name
}

func (left *Widget) Equals(o util.Comparable) bool {
	if left == nil || o == nil {
		return left == o
	}

	right, ok := o.(*Widget)
	if !ok {
		return false
	}

	return left.Name == right.Name &&
		left.Description == right.Description &&
		left.WidgetType == right.WidgetType &&
		left.WidgetSize.Equals(&right.WidgetSize) &&
		util.CompareStringSlices(left.Filters, right.Filters) &&
		left.WidgetPosition.Equals(&right.WidgetPosition) &&
		left.ContentParameters.Equals(&right.ContentParameters)
}

func (d *Widget) Key() string {
	return d.Name
}

func (left *WidgetSize) Equals(right *WidgetSize) bool {

	if left == nil || right == nil {
		return left == right
	}

	return left.Height == right.Height &&
		left.Width == right.Width
}

func (left *WidgetPosition) Equals(right *WidgetPosition) bool {

	if left == nil || right == nil {
		return left == right
	}

	return left.PositionX == right.PositionX &&
		left.PositionY == right.PositionY
}

func (left *WidgetContentParameters) Equals(right *WidgetContentParameters) bool {

	if left == nil || right == nil {
		return left == right
	}

	return util.CompareStringSlices(left.ContentFields, right.ContentFields) &&
		left.ItemsCount == right.ItemsCount &&
		reflect.DeepEqual(right.WidgetOptions, left.WidgetOptions)
}

func (d *Dashboard) HashName() string {
	return HashName(d.Name)
}

func HashName(name string) string {
	h := sha1.New()
	io.WriteString(h, name)
	return hex.EncodeToString(h.Sum(nil))[:4]
}
