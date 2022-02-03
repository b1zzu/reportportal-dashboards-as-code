package rpdac

import (
	"errors"
	"reflect"
	"testing"

	"github.com/b1zzu/reportportal-dashboards-as-code/pkg/reportportal"
	"github.com/google/go-cmp/cmp"
)

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

	got, err := r.Dashboard.GetDashboard("test_project", 1)
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

func TestGetDashboardByName(t *testing.T) {

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
		GetByNameM: func(projectName string, name string) (*reportportal.Dashboard, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, name, "MK E2E Tests Overview")
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

	got, err := r.Dashboard.GetDashboardByName("test_project", "MK E2E Tests Overview")
	if err != nil {
		t.Errorf("ReportPortal.GetDashboardByName returned error: %v", err)
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

func TestGetDashboardByName_NotFound(t *testing.T) {

	mockDashboard := &reportportal.MockDashboardService{
		GetByNameM: func(projectName string, name string) (*reportportal.Dashboard, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, name, "MK E2E Tests Overview")
			return nil, nil, reportportal.NewDashboardNotFoundError(projectName, name)
		},
	}

	r := NewReportPortal(&reportportal.Client{Dashboard: mockDashboard})

	got, err := r.Dashboard.GetDashboardByName("test_project", "MK E2E Tests Overview")
	if err != nil {
		t.Errorf("ReportPortal.GetDashboardByName returned error: %v", err)
	}

	if got != nil {
		t.Errorf("ReportPortal.GetDashboardByName want nil but got %+v", got)
	}
}

func TestGetDashboardByName_Error(t *testing.T) {

	mockDashboard := &reportportal.MockDashboardService{
		GetByNameM: func(projectName string, name string) (*reportportal.Dashboard, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, name, "MK E2E Tests Overview")
			return nil, nil, errors.New("unexpected error")
		},
	}

	r := NewReportPortal(&reportportal.Client{Dashboard: mockDashboard})

	_, err := r.Dashboard.GetDashboardByName("test_project", "MK E2E Tests Overview")
	if err == nil {
		t.Errorf("ReportPortal.GetDashboardByName did not return the error")
	}
}

func TestCreateDashboard(t *testing.T) {

	mockDashboard := &reportportal.MockDashboardService{
		CreateM: func(projectName string, d *reportportal.NewDashboard) (int, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")

			want := &reportportal.NewDashboard{
				Name:        "MK E2E Tests Overview",
				Description: "",
				Share:       true,
			}

			testDeepEqual(t, d, want)
			return 77, nil, nil
		},
		AddWidgetM: func(projectName string, dashboardID int, w *reportportal.DashboardWidget) (string, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, dashboardID, 77)

			want := map[int]*reportportal.DashboardWidget{
				3: {
					WidgetName: "Failed/Skipped/Passed [Last 7 days]",
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
				67: {
					WidgetName: "Unique bugs [Last 7 days]",
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
			}
			testDeepEqual(t, w, want[w.WidgetID])

			return "", nil, nil
		},
	}

	mockWidget := &reportportal.MockWidgetService{
		PostM: func(projectName string, w *reportportal.NewWidget) (int, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")

			want := map[string]struct {
				Widget *reportportal.NewWidget
				ID     int
			}{
				"Failed/Skipped/Passed [Last 7 days] #9eaf": {
					Widget: &reportportal.NewWidget{
						Name:        "Failed/Skipped/Passed [Last 7 days] #9eaf",
						Description: "",
						Share:       true,
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
						Filters: []int{2},
					},
					ID: 3,
				},
				"Unique bugs [Last 7 days] #9eaf": {
					Widget: &reportportal.NewWidget{
						Description: "",
						Share:       true,
						Name:        "Unique bugs [Last 7 days] #9eaf",
						WidgetType:  "uniqueBugTable",
						ContentParameters: reportportal.WidgetContentParameters{
							ContentFields: []string{},
							ItemsCount:    168,
							WidgetOptions: map[string]interface{}{
								"latest": false,
							},
						},
						Filters: []int{2},
					},
					ID: 67,
				},
			}
			testDeepEqual(t, w, want[w.Name].Widget)

			return want[w.Name].ID, nil, nil
		},
	}

	mockFilter := &reportportal.MockFilterService{
		GetByNameM: func(projectName, name string) (*reportportal.Filter, *reportportal.Response, error) {
			return &reportportal.Filter{
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
				Type: "Launch"}, nil, nil
		},
	}

	mockProjectSettings := &reportportal.MockProjectSettingsService{
		GetM: func(projectName string) (*reportportal.ProjectSettings, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")

			return &reportportal.ProjectSettings{
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
			}, nil, nil
		},
	}

	r := NewReportPortal(&reportportal.Client{
		Dashboard:       mockDashboard,
		Widget:          mockWidget,
		Filter:          mockFilter,
		ProjectSettings: mockProjectSettings,
	})

	inputDashboard := &Dashboard{
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
			},
		},
	}

	err := r.Dashboard.CreateDashboard("test_project", inputDashboard)
	if err != nil {
		t.Errorf("ReportPortal.CreateDashboard returned error: %v", err)
	}

	testDeepEqual(t, mockDashboard.Counter, reportportal.MockDashboardServiceCounter{Create: 1, AddWidget: 2})
	testDeepEqual(t, mockWidget.Counter, reportportal.MockWidgetServiceCounter{Post: 2})
	testDeepEqual(t, mockFilter.Counter, reportportal.MockFilterServiceCounter{GetByName: 1})
	testDeepEqual(t, mockProjectSettings.Counter, reportportal.MockProjectSettingsServiceCounter{Get: 1})
}

func TestApplyDashboard_Create(t *testing.T) {

	mockDashboard := &reportportal.MockDashboardService{
		GetByNameM: func(projectName, name string) (*reportportal.Dashboard, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, name, "MK E2E Tests Overview")
			return nil, nil, reportportal.NewDashboardNotFoundError(projectName, name)
		},
		CreateM: func(projectName string, d *reportportal.NewDashboard) (int, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")

			want := &reportportal.NewDashboard{
				Name:        "MK E2E Tests Overview",
				Description: "",
				Share:       true,
			}

			testDeepEqual(t, d, want)
			return 77, nil, nil
		},
		AddWidgetM: func(projectName string, dashboardID int, w *reportportal.DashboardWidget) (string, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, dashboardID, 77)
			return "", nil, nil
		},
	}

	mockWidget := &reportportal.MockWidgetService{
		PostM: func(projectName string, w *reportportal.NewWidget) (int, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			want := map[string]int{
				"Failed/Skipped/Passed [Last 7 days] #9eaf": 3,
				"Unique bugs [Last 7 days] #9eaf":           67,
			}
			return want[w.Name], nil, nil
		},
	}

	mockProjectSettings := &reportportal.MockProjectSettingsService{
		GetM: func(projectName string) (*reportportal.ProjectSettings, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			return &reportportal.ProjectSettings{
				ProjectID: 4,
				SubTypes: reportportal.IssueSubTypes{
					"NO_DEFECT": []reportportal.IssueSubType{},
				},
			}, nil, nil
		},
	}

	r := NewReportPortal(&reportportal.Client{
		Dashboard:       mockDashboard,
		Widget:          mockWidget,
		ProjectSettings: mockProjectSettings,
	})

	inputDashboard := &Dashboard{
		Kind:        DashboardKind,
		Name:        "MK E2E Tests Overview",
		Description: "",
		Widgets: []*Widget{
			{
				Name:              "Failed/Skipped/Passed [Last 7 days]",
				Description:       "",
				WidgetType:        "statisticTrend",
				WidgetSize:        WidgetSize{},
				WidgetPosition:    WidgetPosition{},
				Filters:           []string{},
				ContentParameters: WidgetContentParameters{},
			},
			{
				Name:              "Unique bugs [Last 7 days]",
				Description:       "",
				WidgetType:        "uniqueBugTable",
				WidgetSize:        WidgetSize{},
				WidgetPosition:    WidgetPosition{},
				Filters:           []string{},
				ContentParameters: WidgetContentParameters{},
			},
		},
	}

	err := r.Dashboard.ApplyDashboard("test_project", inputDashboard)
	if err != nil {
		t.Errorf("ReportPortal.ApplyDashboard returned error: %v", err)
	}

	testDeepEqual(t, mockDashboard.Counter, reportportal.MockDashboardServiceCounter{GetByName: 1, Create: 1, AddWidget: 2})
	testDeepEqual(t, mockWidget.Counter, reportportal.MockWidgetServiceCounter{Post: 2})
	testDeepEqual(t, mockProjectSettings.Counter, reportportal.MockProjectSettingsServiceCounter{Get: 1})
}

func TestApplyDashboard_Update(t *testing.T) {

	mockDashboard := &reportportal.MockDashboardService{
		GetByNameM: func(projectName, name string) (*reportportal.Dashboard, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, name, "MK E2E Tests Overview")
			return &reportportal.Dashboard{
				Owner:       "dbizzarr",
				Share:       true,
				ID:          1,
				Name:        "MK E2E Tests Overview",
				Description: "Dashboard to update",
				Widgets: []reportportal.DashboardWidget{
					{
						WidgetName:     "Failed/Skipped [Last 7 days] #9eaf",
						WidgetID:       2,
						WidgetType:     "statisticTrend",
						WidgetSize:     reportportal.DashboardWidgetSize{},
						WidgetPosition: reportportal.DashboardWidgetPosition{},
						Share:          true,
					},
				},
			}, nil, nil
		},
		UpdateM: func(projectName string, dashboardID int, d *reportportal.UpdateDashboard) (string, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, dashboardID, 1)

			want := &reportportal.UpdateDashboard{
				Share:       true,
				Name:        "MK E2E Tests Overview",
				Description: "",
			}
			testDeepEqual(t, d, want)

			return "", nil, nil
		},
		AddWidgetM: func(projectName string, dashboardID int, w *reportportal.DashboardWidget) (string, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, dashboardID, 1)

			want := map[int]*reportportal.DashboardWidget{
				3: {
					WidgetName:     "Failed/Skipped/Passed [Last 7 days]",
					WidgetID:       3,
					WidgetType:     "statisticTrend",
					WidgetSize:     reportportal.DashboardWidgetSize{},
					WidgetPosition: reportportal.DashboardWidgetPosition{},
					Share:          true,
				},
				67: {
					WidgetName:     "Unique bugs [Last 7 days]",
					WidgetID:       67,
					WidgetType:     "uniqueBugTable",
					WidgetSize:     reportportal.DashboardWidgetSize{},
					WidgetPosition: reportportal.DashboardWidgetPosition{},
					Share:          true,
				},
			}
			testDeepEqual(t, w, want[w.WidgetID])

			return "", nil, nil
		},
		RemoveWidgetM: func(projectName string, dashboardID, widgetID int) (string, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, dashboardID, 1)
			testEqual(t, widgetID, 2)
			return "", nil, nil
		},
	}

	mockWidget := &reportportal.MockWidgetService{
		GetM: func(projectName string, id int) (*reportportal.Widget, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, id, 2)

			return &reportportal.Widget{
				Description:       "",
				Owner:             "dbizzarr",
				Share:             true,
				ID:                2,
				Name:              "Failed/Skipped [Last 7 days] #9eaf",
				WidgetType:        "statisticTrend",
				ContentParameters: reportportal.WidgetContentParameters{},
				AppliedFilters:    []reportportal.Filter{},
				Content:           nil,
			}, nil, nil
		},
		PostM: func(projectName string, w *reportportal.NewWidget) (int, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")

			want := map[string]struct {
				Widget *reportportal.NewWidget
				ID     int
			}{
				"Failed/Skipped/Passed [Last 7 days] #9eaf": {
					Widget: &reportportal.NewWidget{
						Name:        "Failed/Skipped/Passed [Last 7 days] #9eaf",
						Description: "",
						Share:       true,
						WidgetType:  "statisticTrend",
						ContentParameters: reportportal.WidgetContentParameters{
							ContentFields: []string{},
						},
						Filters: []int{},
					},
					ID: 3,
				},
				"Unique bugs [Last 7 days] #9eaf": {
					Widget: &reportportal.NewWidget{
						Description: "",
						Share:       true,
						Name:        "Unique bugs [Last 7 days] #9eaf",
						WidgetType:  "uniqueBugTable",
						ContentParameters: reportportal.WidgetContentParameters{
							ContentFields: []string{},
						},
						Filters: []int{},
					},
					ID: 67,
				},
			}
			testDeepEqual(t, w, want[w.Name].Widget)

			return want[w.Name].ID, nil, nil
		},
	}

	mockProjectSettings := &reportportal.MockProjectSettingsService{
		GetM: func(projectName string) (*reportportal.ProjectSettings, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			return &reportportal.ProjectSettings{
				ProjectID: 4,
				SubTypes: reportportal.IssueSubTypes{
					"NO_DEFECT": []reportportal.IssueSubType{},
				},
			}, nil, nil
		},
	}

	r := NewReportPortal(&reportportal.Client{
		Dashboard:       mockDashboard,
		Widget:          mockWidget,
		ProjectSettings: mockProjectSettings,
	})

	inputDashboard := &Dashboard{
		Kind:        DashboardKind,
		Name:        "MK E2E Tests Overview",
		Description: "",
		Widgets: []*Widget{
			{
				Name:              "Failed/Skipped/Passed [Last 7 days]",
				Description:       "",
				WidgetType:        "statisticTrend",
				WidgetSize:        WidgetSize{},
				WidgetPosition:    WidgetPosition{},
				Filters:           []string{},
				ContentParameters: WidgetContentParameters{},
			},
			{
				Name:              "Unique bugs [Last 7 days]",
				Description:       "",
				WidgetType:        "uniqueBugTable",
				WidgetSize:        WidgetSize{},
				WidgetPosition:    WidgetPosition{},
				Filters:           []string{},
				ContentParameters: WidgetContentParameters{},
			},
		},
	}

	err := r.Dashboard.ApplyDashboard("test_project", inputDashboard)
	if err != nil {
		t.Errorf("ReportPortal.ApplyDashboard returned error: %v", err)
	}

	testDeepEqual(t, mockDashboard.Counter, reportportal.MockDashboardServiceCounter{GetByName: 1, Update: 1, AddWidget: 2, RemoveWidget: 1})
	testDeepEqual(t, mockWidget.Counter, reportportal.MockWidgetServiceCounter{Get: 1, Post: 2})
	testDeepEqual(t, mockProjectSettings.Counter, reportportal.MockProjectSettingsServiceCounter{Get: 2})
}

func TestApplyDashboard_Skip(t *testing.T) {

	mockDashboard := &reportportal.MockDashboardService{
		GetByNameM: func(projectName, name string) (*reportportal.Dashboard, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, name, "MK E2E Tests Overview")
			return &reportportal.Dashboard{
				Owner:       "dbizzarr",
				Share:       true,
				ID:          1,
				Name:        "MK E2E Tests Overview",
				Description: "",
				Widgets: []reportportal.DashboardWidget{
					{
						WidgetID:   3,
						Share:      true,
						WidgetName: "Failed/Skipped/Passed [Last 7 days]",
						WidgetType: "statisticTrend",
						WidgetSize: reportportal.DashboardWidgetSize{
							Width:  12,
							Height: 6,
						},
						WidgetPosition: reportportal.DashboardWidgetPosition{
							PositionX: 0,
							PositionY: 13,
						},
					},
					{
						WidgetID:   67,
						Share:      true,
						WidgetName: "Unique bugs [Last 7 days]",
						WidgetType: "uniqueBugTable",
						WidgetSize: reportportal.DashboardWidgetSize{
							Width:  12,
							Height: 7,
						},
						WidgetPosition: reportportal.DashboardWidgetPosition{
							PositionX: 0,
							PositionY: 44,
						},
					},
				},
			}, nil, nil
		},
	}

	mockWidget := &reportportal.MockWidgetService{
		GetM: func(projectName string, id int) (*reportportal.Widget, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")

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

			return widgets[id], nil, nil
		},
	}

	mockProjectSettings := &reportportal.MockProjectSettingsService{
		GetM: func(projectName string) (*reportportal.ProjectSettings, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			return &reportportal.ProjectSettings{
				ProjectID: 4,
				SubTypes: reportportal.IssueSubTypes{
					"NO_DEFECT": []reportportal.IssueSubType{},
				},
			}, nil, nil
		},
	}

	r := NewReportPortal(&reportportal.Client{
		Dashboard:       mockDashboard,
		Widget:          mockWidget,
		ProjectSettings: mockProjectSettings,
	})

	inputDashboard := &Dashboard{
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
					},
					ItemsCount: 168,
					WidgetOptions: map[string]interface{}{
						"timeline": "launch",
						"viewMode": "bar",
						"zoom":     false,
					},
				},
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
			},
		},
	}

	err := r.Dashboard.ApplyDashboard("test_project", inputDashboard)
	if err != nil {
		t.Errorf("ReportPortal.ApplyDashboard returned error: %v", err)
	}

	testDeepEqual(t, mockDashboard.Counter, reportportal.MockDashboardServiceCounter{GetByName: 1})
	testDeepEqual(t, mockWidget.Counter, reportportal.MockWidgetServiceCounter{Get: 2})
	testDeepEqual(t, mockProjectSettings.Counter, reportportal.MockProjectSettingsServiceCounter{Get: 1})
}

func TestDeleteDashboard(t *testing.T) {

	mockDashboard := &reportportal.MockDashboardService{
		GetByNameM: func(projectName, name string) (*reportportal.Dashboard, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, name, "MK E2E Tests Overview")
			return &reportportal.Dashboard{
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
				},
			}, nil, nil
		},
		DeleteM: func(projectName string, id int) (string, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, id, 1)
			return "", nil, nil
		},
	}

	r := NewReportPortal(&reportportal.Client{
		Dashboard: mockDashboard,
	})

	err := r.Dashboard.DeleteDashboard("test_project", "MK E2E Tests Overview")
	if err != nil {
		t.Errorf("ReportPortal.DeleteDashboard returned error: %v", err)
	}

	testDeepEqual(t, mockDashboard.Counter, reportportal.MockDashboardServiceCounter{GetByName: 1, Delete: 1})
}

func TestToDashboard(t *testing.T) {

	inputDashboard := &reportportal.Dashboard{
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

	inputWidgets := []*Widget{
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
				},
				ItemsCount: 168,
				WidgetOptions: map[string]interface{}{
					"timeline": "launch",
					"viewMode": "bar",
					"zoom":     false,
				},
			},
			origin: &reportportal.Widget{},
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
			origin: &reportportal.Widget{},
		},
	}

	got := ToDashboard(inputDashboard, inputWidgets)

	want := &Dashboard{
		Kind:        DashboardKind,
		Name:        "MK E2E Tests Overview",
		Description: "",
		Widgets:     inputWidgets,
		origin:      inputDashboard,
	}

	opts := cmp.Options{
		cmp.AllowUnexported(Dashboard{}, Widget{}),
	}
	if !cmp.Equal(got, want, opts) {
		t.Errorf("ToDashboard got: %+v, want: %+v", got, want)
	}
}

func TestToWidget(t *testing.T) {

	inputWidget := &reportportal.Widget{
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
	}

	inputDashboardWidget := &reportportal.DashboardWidget{
		WidgetName: "Failed/Skipped/Passed",
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
	}

	inputDecodeSubTypesMap := map[string]string{
		"si_1iuqflmhg6hk6": "si",
	}

	got, err := ToWidget(inputWidget, inputDashboardWidget, "9eaf", inputDecodeSubTypesMap)
	if err != nil {
		t.Errorf("ToWidget returned error: %v", err)
	}

	want := &Widget{
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
				"statistics$defects$system_issue$si",
			},
			ItemsCount: 168,
			WidgetOptions: map[string]interface{}{
				"timeline": "launch",
				"viewMode": "bar",
				"zoom":     false,
			},
		},
		origin: inputWidget,
	}

	opts := cmp.Options{
		cmp.AllowUnexported(Widget{}),
	}
	testDeepEqual(t, got, want, opts)
}

func TestToWidget_WithoutDashboardHastag(t *testing.T) {

	inputWidget := &reportportal.Widget{
		Description: "",
		Owner:       "dbizzarr",
		Share:       true,
		ID:          3,
		Name:        "Failed/Skipped/Passed [Last 7 days]",
		WidgetType:  "statisticTrend",
		Content:     nil,
	}

	inputDashboardWidget := &reportportal.DashboardWidget{
		WidgetName: "Failed/Skipped/Passed",
		WidgetID:   3,
		WidgetType: "statisticTrend",
		Share:      true,
	}

	inputDecodeSubTypesMap := map[string]string{
		"si_1iuqflmhg6hk6": "si",
	}

	got, err := ToWidget(inputWidget, inputDashboardWidget, "9eaf", inputDecodeSubTypesMap)
	if err != nil {
		t.Errorf("ToWidget returned error: %v", err)
	}

	want := &Widget{
		Name:        "Failed/Skipped/Passed [Last 7 days]",
		Description: "",
		WidgetType:  "statisticTrend",
		Filters:     []string{},
		ContentParameters: WidgetContentParameters{
			ContentFields: []string{},
		},
		origin: inputWidget,
	}

	opts := cmp.Options{
		cmp.AllowUnexported(Widget{}),
	}
	testDeepEqual(t, got, want, opts)
}

func TestFromWidget(t *testing.T) {

	inputWidget := &Widget{
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
				"statistics$defects$system_issue$si",
			},
			ItemsCount: 168,
			WidgetOptions: map[string]interface{}{
				"timeline": "launch",
				"viewMode": "bar",
				"zoom":     false,
			},
		},
		origin: nil,
	}

	inputFilterMap := map[string]int{
		"mk-e2e-test-suite": 2,
	}
	inputEncodeSubTypesMap := map[string]string{
		"si": "si_1iuqflmhg6hk6",
	}

	gotWidget, gotDashboardWidget, err := FromWidget("9eaf", inputWidget, inputFilterMap, inputEncodeSubTypesMap)
	if err != nil {
		t.Errorf("ToWidget returned error: %v", err)
	}

	wantWidget := &reportportal.NewWidget{
		Description: "",
		Share:       true,
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
		Filters: []int{2},
	}

	wantDashboardWidget := &reportportal.DashboardWidget{
		WidgetName: "Failed/Skipped/Passed [Last 7 days]",
		WidgetType: "statisticTrend",
		WidgetSize: reportportal.DashboardWidgetSize{
			Width:  12,
			Height: 6,
		},
		WidgetPosition: reportportal.DashboardWidgetPosition{
			PositionX: 0,
			PositionY: 13,
		},
		Share: true}

	testDeepEqual(t, gotWidget, wantWidget)
	testDeepEqual(t, gotDashboardWidget, wantDashboardWidget)
}

func TestDecodeFieldsSubTypes(t *testing.T) {

	inputFields := []string{
		"statistics$executions$total",
		"statistics$executions$passed",
		"statistics$executions$failed",
		"statistics$executions$skipped",
		"statistics$defects$product_bug$pb001",
		"statistics$defects$automation_bug$ab001",
		"statistics$defects$system_issue$si001",
		"statistics$defects$no_defect$nd001",
		"statistics$defects$to_investigate$ti001",
		"statistics$defects$system_issue$si_1iuqflmhg6hk6",
		"statistics$defects$product_bug$pb_qdy9r7uu9q9g",
		"statistics$defects$system_issue$si_1h7o519q5xeg5",
		"statistics$defects$automation_bug$ab_uv8mlzz5fqzn",
		"statistics$defects$automation_bug$ab_1ien71b1ve81k",
		"statistics$defects$automation_bug$ab_t4f3ctreg3sl",
	}
	inputSubTypesMap := map[string]string{
		"pb001":            "pb",
		"ab001":            "ab",
		"si001":            "si",
		"nd001":            "nd",
		"ti001":            "ti",
		"si_1iuqflmhg6hk6": "si1",
		"si_1h7o519q5xeg5": "si2",
		"pb_qdy9r7uu9q9g":  "pb1",
		"ab_uv8mlzz5fqzn":  "ab1",
		"ab_1ien71b1ve81k": "ab2",
		"ab_t4f3ctreg3sl":  "ab3",
	}

	expectedFields := []string{
		"statistics$executions$total",
		"statistics$executions$passed",
		"statistics$executions$failed",
		"statistics$executions$skipped",
		"statistics$defects$product_bug$pb",
		"statistics$defects$automation_bug$ab",
		"statistics$defects$system_issue$si",
		"statistics$defects$no_defect$nd",
		"statistics$defects$to_investigate$ti",
		"statistics$defects$system_issue$si1",
		"statistics$defects$product_bug$pb1",
		"statistics$defects$system_issue$si2",
		"statistics$defects$automation_bug$ab1",
		"statistics$defects$automation_bug$ab2",
		"statistics$defects$automation_bug$ab3",
	}

	result, err := DecodeFieldsSubTypes(inputFields, inputSubTypesMap)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(result, expectedFields) {
		t.Errorf("Failed: got %+v but expected %v", result, expectedFields)
	}
}

func TestDashboardEquals(t *testing.T) {

	tests := []*struct {
		description string
		left        *Dashboard
		right       *Dashboard
		expexct     bool
	}{
		{
			description: "Compare empty dashboards should return true",
			left:        &Dashboard{},
			right:       &Dashboard{},
			expexct:     true,
		},
		{
			description: "Compare equal dashboards but only one with the origin filed should return true",
			left: &Dashboard{
				Kind:        DashboardKind,
				Name:        "Test",
				Description: "My test description",
				origin:      &reportportal.Dashboard{ID: 1},
			},
			right: &Dashboard{
				Kind:        DashboardKind,
				Name:        "Test",
				Description: "My test description",
			},
			expexct: true,
		},
		{
			description: "Compare dashboards with differt names should return false",
			left: &Dashboard{
				Kind:        DashboardKind,
				Name:        "Test One",
				Description: "My test description",
			},
			right: &Dashboard{
				Kind:        DashboardKind,
				Name:        "Test",
				Description: "My test description",
			},
			expexct: false,
		},
		{
			description: "Compare dashboards with differt description should return false",
			left: &Dashboard{
				Kind:        DashboardKind,
				Name:        "Test",
				Description: "My test description",
			},
			right: &Dashboard{
				Kind:        DashboardKind,
				Name:        "Test",
				Description: "My updated test description",
			},
			expexct: false,
		},
		{
			description: "Compare dashboards wehre one with nil widgets should return false",
			left: &Dashboard{
				Widgets: []*Widget{},
			},
			right: &Dashboard{
				Widgets: nil,
			},
			expexct: false,
		},
		{
			description: "Compare dashboards with different widgets should return false",
			left: &Dashboard{
				Widgets: []*Widget{
					{Name: "One", Description: "One description"},
					{Name: "Two", Description: "Two description"},
				},
			},
			right: &Dashboard{
				Widgets: []*Widget{
					{Name: "Two", Description: "One description"},
					{Name: "Three", Description: "Three description"},
				},
			},
			expexct: false,
		},
		{
			description: "Compare dashboards with same widgets and same order should return true",
			left: &Dashboard{
				Widgets: []*Widget{
					{Name: "One", Description: "One description"},
					{Name: "Two", Description: "Two description"},
				},
			},
			right: &Dashboard{
				Widgets: []*Widget{
					{Name: "One", Description: "One description"},
					{Name: "Two", Description: "Two description"},
				},
			},
			expexct: true,
		},
		{
			description: "Compare dashboards with same widgets should return true",
			left: &Dashboard{
				Widgets: []*Widget{
					{Name: "One", Description: "One description"},
					{Name: "Two", Description: "Two description"},
				},
			},
			right: &Dashboard{
				Widgets: []*Widget{
					{Name: "Two", Description: "Two description"},
					{Name: "One", Description: "One description"},
				},
			},
			expexct: true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			r := test.left.Equals(test.right)
			if r != test.expexct {
				t.Errorf("expected '%t' but got '%t'", test.expexct, r)
			}
		})
	}
}

func TestWidgetEquals(t *testing.T) {

	tests := []*struct {
		description string
		left        *Widget
		right       *Widget
		expexct     bool
	}{
		{
			description: "Compare empty widgets should return true",
			left:        &Widget{},
			right:       &Widget{},
			expexct:     true,
		},
		{
			description: "Compare equal widgets but only one with the origin filed should return true",
			left: &Widget{
				Name:        "Test",
				Description: "My test description",
				origin:      &reportportal.Widget{ID: 1},
			},
			right: &Widget{
				Name:        "Test",
				Description: "My test description",
			},
			expexct: true,
		},
		{
			description: "Compare widgets with differt names should return false",
			left: &Widget{
				Name:        "Test One",
				Description: "My test description",
			},
			right: &Widget{
				Name:        "Test",
				Description: "My test description",
			},
			expexct: false,
		},
		{
			description: "Compare widgets with differt description should return false",
			left: &Widget{
				Name:        "Test",
				Description: "My test description",
			},
			right: &Widget{
				Name:        "Test",
				Description: "My other description",
			},
			expexct: false,
		},
		{
			description: "Compare widgets with differt widget type should return false",
			left: &Widget{
				Name:       "Test",
				WidgetType: "Test type",
			},
			right: &Widget{
				Name:       "Test",
				WidgetType: "Other test type",
			},
			expexct: false,
		},
		{
			description: "Compare widgets with differt size should return false",
			left: &Widget{
				Name:       "Test",
				WidgetSize: WidgetSize{Width: 1, Height: 1},
			},
			right: &Widget{
				Name:       "Test",
				WidgetSize: WidgetSize{Width: 2, Height: 1},
			},
			expexct: false,
		},
		{
			description: "Compare widgets with differt position should return false",
			left: &Widget{
				Name:           "Test",
				WidgetPosition: WidgetPosition{PositionX: 4, PositionY: 9},
			},
			right: &Widget{
				Name:           "Test",
				WidgetPosition: WidgetPosition{PositionX: 7, PositionY: 9},
			},
			expexct: false,
		},
		{
			description: "Compare widgets with differt filters should return false",
			left: &Widget{
				Name:    "Test",
				Filters: []string{"one", "two"},
			},
			right: &Widget{
				Name:    "Test",
				Filters: []string{"one", "three"},
			},
			expexct: false,
		},
		{
			description: "Compare widgets with differt content parameters should return false",
			left: &Widget{
				Name: "Test",
				ContentParameters: WidgetContentParameters{
					ContentFields: []string{"one"},
					ItemsCount:    1,
					WidgetOptions: map[string]interface{}{
						"Hello":  "world",
						"Sounds": true,
						"Numero": 1,
					},
				},
			},
			right: &Widget{
				Name: "Test",
				ContentParameters: WidgetContentParameters{
					ContentFields: []string{"two"},
					ItemsCount:    2,
					WidgetOptions: map[string]interface{}{
						"Hello":  "world",
						"Sounds": true,
						"Numero": 1,
					},
				},
			},
			expexct: false,
		},
		{
			description: "Compare equal widgets should return true",
			left: &Widget{
				Name:           "Test",
				Description:    "My test description",
				WidgetType:     "Test type",
				WidgetSize:     WidgetSize{Width: 1, Height: 1},
				WidgetPosition: WidgetPosition{PositionX: 4, PositionY: 9},
				Filters:        []string{"one", "two"},
				ContentParameters: WidgetContentParameters{
					ContentFields: []string{"three", "one"},
					ItemsCount:    1,
					WidgetOptions: map[string]interface{}{
						"Sounds": true,
						"Hello":  "world",
						"Numero": 1,
					},
				},
			},
			right: &Widget{
				Name:           "Test",
				Description:    "My test description",
				WidgetType:     "Test type",
				WidgetSize:     WidgetSize{Width: 1, Height: 1},
				WidgetPosition: WidgetPosition{PositionX: 4, PositionY: 9},
				Filters:        []string{"two", "one"},
				ContentParameters: WidgetContentParameters{
					ContentFields: []string{"one", "three"},
					ItemsCount:    1,
					WidgetOptions: map[string]interface{}{
						"Hello":  "world",
						"Sounds": true,
						"Numero": 1,
					},
				},
			},
			expexct: true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			r := test.left.Equals(test.right)
			if r != test.expexct {
				t.Errorf("expected '%t' but got '%t'", test.expexct, r)
			}
		})
	}
}

func TestHashName(t *testing.T) {

	input := "MK E2E Tests Overview"
	got := HashName(input)

	want := "9eaf"
	if got != want {
		t.Errorf("want %v but got %v", want, got)
	}
}
