package rpdac

import (
	"github.com/b1zzu/reportportal-dashboards-as-code/pkg/reportportal"
)

type ReportPortal struct {
	client *reportportal.Client

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	Dashboard IDashboardService
	Filter    IFilterService
}

type service struct {
	client *reportportal.Client
}

func NewReportPortal(c *reportportal.Client) *ReportPortal {
	r := &ReportPortal{client: c}
	r.common.client = c
	r.Dashboard = (*DashboardService)(&r.common)
	r.Filter = (*FilterService)(&r.common)
	return r
}
