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

	return r.loadDashboard(project, d)
}

func (r *ReportPortal) loadDashboard(project string, d *reportportal.Dashboard) (*Dashboard, error) {

	widgets := make([]*Widget, len(d.Widgets))

	decodeSubTypesMap, err := r.decodeSubTypseMap(project)
	if err != nil {
		return nil, err
	}

	dashboardHash := HashName(d.Name)

	// retrieve all widgets definitions
	for i, dw := range d.Widgets {
		w, _, err := r.client.Widget.Get(project, dw.WidgetID)
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

func (r *ReportPortal) GetFilter(project string, filterID int) (*Filter, error) {

	// retireve the filter defintion
	f, _, err := r.client.Filter.GetByID(project, filterID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving filter %d from project %s: %w", filterID, project, err)
	}

	return ToFilter(f), nil
}

func (r *ReportPortal) GetDashboardByName(project, dashboardName string) (*Dashboard, error) {

	d, _, err := r.client.Dashboard.GetByName(project, dashboardName)
	if err != nil {
		if _, ok := err.(*reportportal.DashboardNotFoundError); ok {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return r.loadDashboard(project, d)
}

func (r *ReportPortal) GetFilterByName(project, filterName string) (*Filter, error) {

	f, _, err := r.client.Filter.GetByName(project, filterName)
	if err != nil {
		if _, ok := err.(*reportportal.FilterNotFoundError); ok {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return ToFilter(f), nil
}

func (r *ReportPortal) CreateDashboard(project string, d *Dashboard) error {

	// resolve all filters
	filtersMap, err := r.filtersMap(project, d.Widgets)
	if err != nil {
		return err
	}

	encodeSubTypesMap, err := r.encodeSubTypseMap(project)
	if err != nil {
		return err
	}

	dashboardID, _, err := r.client.Dashboard.Create(project, &reportportal.NewDashboard{Name: d.Name, Description: d.Description, Share: true})
	if err != nil {
		return fmt.Errorf("error creating dashboard %s: %w", d.Name, err)
	}
	log.Printf("dashboard %s created with id: %d", d.Name, dashboardID)

	return r.createWidgets(project, dashboardID, d, filtersMap, encodeSubTypesMap)
}

func (r *ReportPortal) createWidgets(
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

		widgetID, _, err := r.client.Widget.Post(project, nw)
		if err != nil {
			return fmt.Errorf("error creating widget %s: %w", w.Name, err)
		}

		dw.WidgetID = widgetID

		_, _, err = r.client.Dashboard.AddWidget(project, dashboardID, dw)
		if err != nil {
			return fmt.Errorf("error adding widget %s to dashboard %s: %w", w.Name, dashboard.Name, err)
		}
		log.Printf("added \"%s\" widget to \"%s\" dashboard", w.Name, dashboard.Name)
	}
	return nil
}

func (r *ReportPortal) CreateFilter(project string, f *Filter) error {

	filterID, _, err := r.client.Filter.Create(project, FilterToNewFilter(f))
	if err != nil {
		return fmt.Errorf("error creating filter %s: %w", f.Name, err)
	}

	log.Printf("filter %s created with id: %d", f.Name, filterID)
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

// Create or Recreate the dashboard
func (r *ReportPortal) ApplyDashboard(project string, targetDashboard *Dashboard) error {

	currentDashboard, err := r.GetDashboardByName(project, targetDashboard.Name)
	if err != nil {
		return fmt.Errorf("error retrieving dashboard \"%s\" by name: %w", targetDashboard.Name, err)
	}

	if currentDashboard != nil {

		if currentDashboard.Equals(targetDashboard) {
			log.Printf("skip apply for dashboard \"%s\"", targetDashboard.Name)
			return nil
		}

		return r.updateDashboard(project, currentDashboard, targetDashboard)
	}

	return r.CreateDashboard(project, targetDashboard)
}

func (r *ReportPortal) ApplyFilter(project string, targetFilter *Filter) error {

	currentFilter, err := r.GetFilterByName(project, targetFilter.Name)
	if err != nil {
		return fmt.Errorf("error retrieving filter \"%s\" by name: %w", targetFilter.Name, err)
	}

	if currentFilter != nil {

		if currentFilter.Equals(targetFilter) {
			log.Printf("skip apply for filter \"%s\"", targetFilter.Name)
			return nil
		}

		return r.updateFilter(project, currentFilter, targetFilter)
	}

	return r.CreateFilter(project, targetFilter)

}

func (r *ReportPortal) updateDashboard(project string, currentDashboard, targetDashboard *Dashboard) error {

	// resolve all filters
	filtersMap, err := r.filtersMap(project, targetDashboard.Widgets)
	if err != nil {
		return err
	}

	encodeSubTypesMap, err := r.encodeSubTypseMap(project)
	if err != nil {
		return err
	}

	dashboardID := currentDashboard.origin.ID

	// delete all widgets from the current dashboard so we can recreate them as expected by the target dashboard
	for _, w := range currentDashboard.Widgets {
		_, _, err := r.client.Dashboard.RemoveWidget(project, dashboardID, w.origin.ID)
		if err != nil {
			return fmt.Errorf("error removing widget \"%s\" from dashboard \"%s\": %w", w.Name, currentDashboard.Name, err)
		}
	}

	u := &reportportal.UpdateDashboard{Name: targetDashboard.Name, Description: targetDashboard.Description, Share: true}
	_, _, err = r.client.Dashboard.Update(project, dashboardID, u)
	if err != nil {
		return fmt.Errorf("error updating dashboard %s: %w", targetDashboard.Name, err)
	}
	log.Printf("updated \"%s\" dashboard", targetDashboard.Name)

	return r.createWidgets(project, dashboardID, targetDashboard, filtersMap, encodeSubTypesMap)
}

func (r *ReportPortal) updateFilter(project string, currentFilter, targetFilter *Filter) error {

	_, _, err := r.client.Filter.Update(project, currentFilter.origin.ID, FilterToUpdateFilter(targetFilter))
	if err != nil {
		return fmt.Errorf("error updating filter \"%s\": %w", targetFilter.Name, err)
	}

	log.Printf("update \"%s\" filter", targetFilter.Name)
	return nil
}

func (r *ReportPortal) decodeSubTypseMap(project string) (map[string]string, error) {
	ps, _, err := r.client.ProjectSettings.Get(project)
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

func (r *ReportPortal) encodeSubTypseMap(project string) (map[string]string, error) {

	m, err := r.decodeSubTypseMap(project)
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

func (r *ReportPortal) filtersMap(project string, widgets []*Widget) (map[string]int, error) {
	filtersMap := make(map[string]int)
	for _, w := range widgets {
		for _, filterName := range w.Filters {

			if _, ok := filtersMap[filterName]; ok {
				// filter already resolved
				continue
			}

			f, _, err := r.client.Filter.GetByName(project, filterName)
			if err != nil {
				return nil, fmt.Errorf("error resolving filter \"%s\" in widget \"%s\": %w", filterName, w.Name, err)
			}

			filtersMap[filterName] = f.ID
		}
	}
	return filtersMap, nil
}
