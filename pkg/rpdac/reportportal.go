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

		widgets[i] = ToWidget(w, dw)
	}

	return ToDashboard(d, widgets), nil
}

func (r *ReportPortal) CreateDashboard(project string, d *Dashboard) error {

	dashboardHash := d.HashName()

	// resolve all filters
	filtersMap := make(map[string]int)
	for _, w := range d.Widgets {
		for _, filterName := range w.Filters {

			if _, ok := filtersMap[filterName]; ok {
				// filter already resolved
				continue
			}

			f, _, err := r.client.Filter.GetByName(project, filterName)
			if err != nil {
				return fmt.Errorf("error resolving filter \"%s\" in widget \"%s\" in dashboard \"%s\": %w", filterName, w.Name, d.Name, err)
			}

			filtersMap[filterName] = f.ID
		}
	}

	dashboardID, _, err := r.client.Dashboard.Create(project, &reportportal.NewDashboard{Name: d.Name, Share: true})
	if err != nil {
		return fmt.Errorf("error creating dashboard %s: %w", d.Name, err)
	}
	log.Printf("dashboard %s created with id: %d", d.Name, dashboardID)

	for _, w := range d.Widgets {

		filters := make([]int, len(w.Filters))
		for j, f := range w.Filters {
			filters[j] = filtersMap[f]
		}

		nw, dw := FromWidget(dashboardHash, w, filters)

		widgetID, _, err := r.client.Widget.Post(project, nw)
		if err != nil {
			return fmt.Errorf("error creating widget %s: %w", w.Name, err)
		}
		log.Printf("widget %s created with id %d", w.Name, widgetID)

		dw.WidgetID = widgetID

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
