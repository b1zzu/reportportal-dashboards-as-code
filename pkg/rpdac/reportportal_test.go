package rpdac

import (
	"testing"

	"github.com/b1zzu/reportportal-dashboards-as-code/pkg/reportportal"
	"github.com/google/go-cmp/cmp"
)

type Tesso struct {
	M int
}

func testDeepEqual(t *testing.T, got, want interface{}, opts ...cmp.Option) {
	t.Helper()

	if !cmp.Equal(got, want, opts...) {
		t.Errorf("Want (+) but got (-): %s", cmp.Diff(got, want, opts...))
	}
}

func testEqual(t *testing.T, got, want interface{}) {
	if got != want {
		t.Errorf("Want \"%s\" but got \"%s\"", want, got)
	}
}

func TestGetDashboard(t *testing.T) {

	dashboard := &reportportal.Dashboard{
		Owner: "dbizzarr",
		Share: true,
		ID:    1,
		Name:  "MK E2E Tests Overview",
		Widgets: []reportportal.DashboardWidget{
			{
				WidgetName: "Failed/Skipped/Passed [Last 7 days] #9eaf",
				WidgetID:   3,
				WidgetType: "statisticTrend",
				WidgetSize: reportportal.DashboardWidgetSize{
					Width:  12,
					Height: 6,
				},
				WidgetPosition: reportportal.DashboardWidgetPosition{
					PositionX: 0,
					PositionY: 13,
				},
				Share: true,
			},
			{
				WidgetName: "Unique bugs [Last 7 days] #9eaf",
				WidgetID:   67,
				WidgetType: "uniqueBugTable",
				WidgetSize: reportportal.DashboardWidgetSize{
					Width:  12,
					Height: 7,
				},
				WidgetPosition: reportportal.DashboardWidgetPosition{
					PositionX: 0,
					PositionY: 44,
				},
				Share: true,
			},
		},
	}

	widgets := map[int]*reportportal.Widget{
		3: {
			Description: "",
			Owner:       "dbizzarr",
			Share:       true,
			ID:          3,
			Name:        "Failed/Skipped/Passed [Last 7 days] #9eaf",
			WidgetType:  "statisticTrend",
			ContentParameters: reportportal.WidgetContentParameters{
				ContentFields: []string{
					"statistics$executions$passed",
					"statistics$executions$failed",
					"statistics$executions$skipped",
					"statistics$defects$system_issue$si_1iuqflmhg6hk6",
				},
				ItemsCount: 168,
				WidgetOptions: map[string]interface{}{
					"zoom":     false,
					"timeline": "launch",
					"viewMode": "bar",
				},
			},
			AppliedFilters: []reportportal.Filter{
				{
					Owner: "dbizzarr",
					Share: true,
					ID:    2,
					Name:  "mk-e2e-test-suite",
					Conditions: []reportportal.FilterCondition{
						{
							FilteringField: "name",
							Condition:      "eq",
							Value:          "mk-e2e-test-suite",
						},
					},
					Orders: []reportportal.FilterOrder{
						{
							SortingColumn: "startTime",
							IsAsc:         false,
						},
						{
							SortingColumn: "number",
							IsAsc:         false,
						},
					},
					Type: "Launch",
				},
			},
			Content: nil,
		},

		67: {
			Description: "",
			Owner:       "dbizzarr",
			Share:       true,
			ID:          67,
			Name:        "Unique bugs [Last 7 days] #9eaf",
			WidgetType:  "uniqueBugTable",
			ContentParameters: reportportal.WidgetContentParameters{
				ContentFields: []string{},
				ItemsCount:    168,
				WidgetOptions: map[string]interface{}{
					"latest": false,
				},
			},
			AppliedFilters: []reportportal.Filter{
				{
					Owner: "dbizzarr",
					Share: true,
					ID:    2,
					Name:  "mk-e2e-test-suite",
					Conditions: []reportportal.FilterCondition{
						{
							FilteringField: "name",
							Condition:      "eq",
							Value:          "mk-e2e-test-suite",
						},
					},
					Orders: []reportportal.FilterOrder{
						{
							SortingColumn: "startTime",
							IsAsc:         false,
						},
						{
							SortingColumn: "number",
							IsAsc:         false,
						},
					},
					Type: "Launch",
				},
			},
		},
	}

	projectSettings := &reportportal.ProjectSettings{
		ProjectID: 4,
		SubTypes: reportportal.IssueSubTypes{
			"NO_DEFECT": []reportportal.IssueSubType{{
				ID:        4,
				Locator:   "nd001",
				TypeRef:   "NO_DEFECT",
				LongName:  "No Defect",
				ShortName: "ND",
				Color:     "#777777",
			}},
			"SYSTEM_ISSUE": []reportportal.IssueSubType{{
				ID:        5,
				Locator:   "si001",
				TypeRef:   "SYSTEM_ISSUE",
				LongName:  "System Issue",
				ShortName: "SI",
				Color:     "#0274d1",
			}, {
				ID:        12,
				Locator:   "si_1iuqflmhg6hk6",
				TypeRef:   "SYSTEM_ISSUE",
				LongName:  "Kafka Cluster at Capacity",
				ShortName: "KCC",
				Color:     "#00b0ff",
			}},
		},
	}

	mockDashboard := &reportportal.MockDashboardService{
		GetByIDM: func(projectName string, id int) (*reportportal.Dashboard, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, id, 1)
			return dashboard, nil, nil
		},
	}

	mockWidget := &reportportal.MockWidgetService{
		GetM: func(projectName string, id int) (*reportportal.Widget, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			return widgets[id], nil, nil
		},
	}

	mockProjectSettings := &reportportal.MockProjectSettingsService{
		GetM: func(projectName string) (*reportportal.ProjectSettings, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			return projectSettings, nil, nil
		},
	}

	r := NewReportPortal(&reportportal.Client{
		Dashboard:       mockDashboard,
		Widget:          mockWidget,
		ProjectSettings: mockProjectSettings})

	got, err := r.GetDashboard("test_project", 1)
	if err != nil {
		t.Errorf("ReportPortal.GetDashboard returned error: %v", err)
	}

	want := &Dashboard{
		Kind:        DashboardKind,
		Name:        "MK E2E Tests Overview",
		Description: "",
		Widgets: []*Widget{
			{
				Name:        "Failed/Skipped/Passed [Last 7 days]",
				Description: "",
				WidgetType:  "statisticTrend",
				WidgetSize: WidgetSize{
					Width:  12,
					Height: 6,
				},
				WidgetPosition: WidgetPosition{
					PositionX: 0,
					PositionY: 13,
				},
				Filters: []string{"mk-e2e-test-suite"},
				ContentParameters: WidgetContentParameters{
					ContentFields: []string{
						"statistics$executions$passed",
						"statistics$executions$failed",
						"statistics$executions$skipped",
						"statistics$defects$system_issue$KCC",
					},
					ItemsCount: 168,
					WidgetOptions: map[string]interface{}{
						"timeline": "launch",
						"viewMode": "bar",
						"zoom":     false,
					},
				},
				origin: widgets[3],
			},
			{
				Name:        "Unique bugs [Last 7 days]",
				Description: "",
				WidgetType:  "uniqueBugTable",
				WidgetSize: WidgetSize{
					Width:  12,
					Height: 7,
				},
				WidgetPosition: WidgetPosition{
					PositionX: 0,
					PositionY: 44,
				},
				Filters: []string{"mk-e2e-test-suite"},
				ContentParameters: WidgetContentParameters{
					ContentFields: []string{},
					ItemsCount:    168,
					WidgetOptions: map[string]interface{}{
						"latest": false,
					},
				},
				origin: widgets[67],
			},
		},
		origin: dashboard,
	}

	opts := cmp.Options{
		cmp.AllowUnexported(Dashboard{}, Widget{}),
	}
	testDeepEqual(t, got, want, opts)
}
