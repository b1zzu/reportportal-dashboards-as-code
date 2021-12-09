package rpdac

import (
	"fmt"
	"log"

	"github.com/b1zzu/reportportal-dashboards-as-code/pkg/reportportal"
)

type ReportPortal struct {
	client *reportportal.Client
}

func NewReportPortal(c *reportportal.Client) *ReportPortal {
	return &ReportPortal{client: c}
}

func (r *ReportPortal) GetDashboard(project string, dashboardID int) (*Dashboard, error) {

	// retireve the dashboard defintion
	d, _, err := r.client.Dashboard.GetByID(project, dashboardID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving dashboard %d from project %s: %w", dashboardID, project, err)
	}

	widgets := make([]*Widget, len(d.Widgets))

	// retrieve all widgets definitions
	for i, dw := range d.Widgets {
		w, _, err := r.client.Widget.Get(project, dw.WidgetID)
		if err != nil {
			return nil, fmt.Errorf("error retrieving widget %d from project %s: %w", dw.WidgetID, project, err)
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

func (r *ReportPortal) CreateDashboard(project string, d *Dashboard) error {

	hash := d.HashName()

	dashboardID, _, err := r.client.Dashboard.Create(project, &reportportal.NewDashboard{Name: d.Name, Share: true})
	if err != nil {
		return fmt.Errorf("error creating dashboard %s: %w", d.Name, err)
	}
	log.Printf("dashboard %s created with id: %d", d.Name, dashboardID)

	for _, w := range d.Widgets {

		nw := &reportportal.NewWidget{
			// For the rpdac tool the widget name is not unique across all dashboards, while fore ReportPortal it is,
			// by adding the dashboard name sha to the widget name we make it generic
			Name:        fmt.Sprintf("%s #%s", w.Name, hash),
			Description: w.Description,
			Share:       true,
			WidgetType:  w.WidgetType,
			Filters:     w.Filters,
			ContentParameters: &reportportal.WidgetContentParameters{
				ItemsCount:    w.ContentParameters.ItemsCount,
				ContentFields: w.ContentParameters.ContentFields,
				WidgetOptions: w.ContentParameters.WidgetOptions,
			},
		}
		widgetID, _, err := r.client.Widget.Post(project, nw)
		if err != nil {
			return fmt.Errorf("error creating widget %s: %w", w.Name, err)
		}
		log.Printf("widget %s created with id %d", w.Name, widgetID)

		dw := &reportportal.DashboardWidget{
			WidgetID:       widgetID,
			Share:          true,
			WidgetName:     w.Name,
			WidgetType:     w.WidgetType,
			WidgetSize:     &reportportal.DashboardWidgetSize{Width: w.WidgetSize.Width, Height: w.WidgetSize.Height},
			WidgetPosition: &reportportal.DashboardWidgetPosition{PositionX: w.WidgetPosition.PositionX, PositionY: w.WidgetPosition.PositionY},
		}

		_, _, err = r.client.Dashboard.AddWidget(project, dashboardID, dw)
		if err != nil {
			return fmt.Errorf("error adding widget %s to dashboard %s: %w", w.Name, d.Name, err)
		}
		log.Printf("widget %s added to dashboard %s", w.Name, d.Name)
	}
	return nil
}

// Delete the Dashboard with the given name and Widgets created for it
func (r *ReportPortal) DeleteDashboard(project, dashboard string) error {

	d, _, err := r.client.Dashboard.GetByName(project, dashboard)
	if err != nil {
		if _, ok := err.(*reportportal.DashboardNotFoundError); ok {
			// ignore
		} else {
			return err
		}
	}

	// because we have ignored the error in case of DashboardNotFoundError d can also be nil
	if d != nil {
		_, _, err = r.client.Dashboard.Delete(project, d.ID)
		if err != nil {
			return err
		}
	}

	// Widgets are deleted automatically if not used buy any dashboard
	return nil
}
