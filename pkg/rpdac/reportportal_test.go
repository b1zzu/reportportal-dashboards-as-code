package rpdac

import (
	"errors"
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
	t.Helper()

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

func TestGetFilter(t *testing.T) {

	filter := &reportportal.Filter{
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
	}

	mockFilter := &reportportal.MockFilterService{
		GetByIDM: func(projectName string, id int) (*reportportal.Filter, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, id, 2)
			return filter, nil, nil
		},
	}

	r := NewReportPortal(&reportportal.Client{
		Filter: mockFilter,
	})

	got, err := r.GetFilter("test_project", 2)
	if err != nil {
		t.Errorf("ReportPortal.GetFilter returned error: %v", err)
	}

	want := &Filter{
		Kind:        "Filter",
		Name:        "mk-e2e-test-suite",
		Type:        "Launch",
		Description: "",
		Conditions: []FilterCondition{
			{FilteringField: "name", Condition: "eq", Value: "mk-e2e-test-suite"},
		},
		Orders: []FilterOrder{
			{SortingColumn: "startTime", IsAsc: false},
			{SortingColumn: "number", IsAsc: false},
		},
		origin: filter,
	}

	testDeepEqual(t, got, want, cmp.AllowUnexported(Filter{}))
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

	got, err := r.GetDashboardByName("test_project", "MK E2E Tests Overview")
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

	got, err := r.GetDashboardByName("test_project", "MK E2E Tests Overview")
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

	_, err := r.GetDashboardByName("test_project", "MK E2E Tests Overview")
	if err == nil {
		t.Errorf("ReportPortal.GetDashboardByName did not return the error")
	}
}

func TestGetFilterByName(t *testing.T) {

	filter := &reportportal.Filter{
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
	}

	mockFilter := &reportportal.MockFilterService{
		GetByNameM: func(projectName string, name string) (*reportportal.Filter, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, name, "mk-e2e-test-suite")
			return filter, nil, nil
		},
	}

	r := NewReportPortal(&reportportal.Client{
		Filter: mockFilter,
	})

	got, err := r.GetFilterByName("test_project", "mk-e2e-test-suite")
	if err != nil {
		t.Errorf("ReportPortal.GetFilterByName returned error: %v", err)
	}

	want := &Filter{
		Kind:        "Filter",
		Name:        "mk-e2e-test-suite",
		Type:        "Launch",
		Description: "",
		Conditions: []FilterCondition{
			{FilteringField: "name", Condition: "eq", Value: "mk-e2e-test-suite"},
		},
		Orders: []FilterOrder{
			{SortingColumn: "startTime", IsAsc: false},
			{SortingColumn: "number", IsAsc: false},
		},
		origin: filter,
	}

	testDeepEqual(t, got, want, cmp.AllowUnexported(Filter{}))
}

func TestGetFilterByName_NotFound(t *testing.T) {

	mockFilter := &reportportal.MockFilterService{
		GetByNameM: func(projectName string, name string) (*reportportal.Filter, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, name, "mk-e2e-test-suite")
			return nil, nil, reportportal.NewFilterNotFoundError(projectName, name)
		},
	}

	r := NewReportPortal(&reportportal.Client{
		Filter: mockFilter,
	})

	got, err := r.GetFilterByName("test_project", "mk-e2e-test-suite")
	if err != nil {
		t.Errorf("ReportPortal.GetFilterByName returned error: %v", err)
	}

	if got != nil {
		t.Errorf("ReportPortal.GetFilterByName want nil but got %+v", got)
	}
}

func TestGetFilterByName_Error(t *testing.T) {

	mockFilter := &reportportal.MockFilterService{
		GetByNameM: func(projectName string, name string) (*reportportal.Filter, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, name, "mk-e2e-test-suite")
			return nil, nil, errors.New("unexpected error")
		},
	}

	r := NewReportPortal(&reportportal.Client{
		Filter: mockFilter,
	})

	_, err := r.GetFilterByName("test_project", "mk-e2e-test-suite")
	if err == nil {
		t.Errorf("ReportPortal.GetFilterByName did not return the error")
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

	err := r.CreateDashboard("test_project", inputDashboard)
	if err != nil {
		t.Errorf("ReportPortal.CreateDashboard returned error: %v", err)
	}

	testDeepEqual(t, mockDashboard.Counter, reportportal.MockDashboardServiceCounter{Create: 1, AddWidget: 2})
	testDeepEqual(t, mockWidget.Counter, reportportal.MockWidgetServiceCounter{Post: 2})
	testDeepEqual(t, mockFilter.Counter, reportportal.MockFilterServiceCounter{GetByName: 1})
	testDeepEqual(t, mockProjectSettings.Counter, reportportal.MockProjectSettingsServiceCounter{Get: 1})
}

func TestCreateFilter(t *testing.T) {

	mockFilter := &reportportal.MockFilterService{
		CreateM: func(projectName string, f *reportportal.NewFilter) (int, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")

			want := &reportportal.NewFilter{
				Share: true,
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
			}

			testDeepEqual(t, f, want)
			return 2, nil, nil
		},
	}

	r := NewReportPortal(&reportportal.Client{
		Filter: mockFilter,
	})

	inputFilter := &Filter{
		Kind:        "Filter",
		Name:        "mk-e2e-test-suite",
		Description: "",
		Type:        "Launch",
		Conditions: []FilterCondition{
			{
				FilteringField: "name",
				Condition:      "eq",
				Value:          "mk-e2e-test-suite",
			},
		},
		Orders: []FilterOrder{
			{
				SortingColumn: "startTime",
				IsAsc:         false,
			},
			{
				SortingColumn: "number",
				IsAsc:         false,
			},
		},
	}

	err := r.CreateFilter("test_project", inputFilter)
	if err != nil {
		t.Errorf("ReportPortal.CreateFilter returned error: %v", err)
	}

	testDeepEqual(t, mockFilter.Counter, reportportal.MockFilterServiceCounter{Create: 1})
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

	err := r.DeleteDashboard("test_project", "MK E2E Tests Overview")
	if err != nil {
		t.Errorf("ReportPortal.DeleteDashboard returned error: %v", err)
	}

	testDeepEqual(t, mockDashboard.Counter, reportportal.MockDashboardServiceCounter{GetByName: 1, Delete: 1})
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
		Kind:        "Dashboard",
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

	err := r.ApplyDashboard("test_project", inputDashboard)
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
		Kind:        "Dashboard",
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

	err := r.ApplyDashboard("test_project", inputDashboard)
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
		Kind:        "Dashboard",
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

	err := r.ApplyDashboard("test_project", inputDashboard)
	if err != nil {
		t.Errorf("ReportPortal.ApplyDashboard returned error: %v", err)
	}

	testDeepEqual(t, mockDashboard.Counter, reportportal.MockDashboardServiceCounter{GetByName: 1})
	testDeepEqual(t, mockWidget.Counter, reportportal.MockWidgetServiceCounter{Get: 2})
	testDeepEqual(t, mockProjectSettings.Counter, reportportal.MockProjectSettingsServiceCounter{Get: 1})
}

func TestApplyFilter_Create(t *testing.T) {

	mockFilter := &reportportal.MockFilterService{
		GetByNameM: func(projectName, name string) (*reportportal.Filter, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, name, "mk-e2e-test-suite")
			return nil, nil, reportportal.NewFilterNotFoundError(projectName, name)
		},
		CreateM: func(projectName string, f *reportportal.NewFilter) (int, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")

			want := &reportportal.NewFilter{
				Share: true,
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
			}

			testDeepEqual(t, f, want)
			return 2, nil, nil
		},
	}

	r := NewReportPortal(&reportportal.Client{
		Filter: mockFilter,
	})

	inputFilter := &Filter{
		Kind:        "Filter",
		Name:        "mk-e2e-test-suite",
		Description: "",
		Type:        "Launch",
		Conditions: []FilterCondition{
			{
				FilteringField: "name",
				Condition:      "eq",
				Value:          "mk-e2e-test-suite",
			},
		},
		Orders: []FilterOrder{
			{
				SortingColumn: "startTime",
				IsAsc:         false,
			},
			{
				SortingColumn: "number",
				IsAsc:         false,
			},
		},
	}

	err := r.ApplyFilter("test_project", inputFilter)
	if err != nil {
		t.Errorf("ReportPortal.ApplyFilter returned error: %v", err)
	}

	testDeepEqual(t, mockFilter.Counter, reportportal.MockFilterServiceCounter{GetByName: 1, Create: 1})
}

func TestApplyFilter_Update(t *testing.T) {
	mockFilter := &reportportal.MockFilterService{
		GetByNameM: func(projectName, name string) (*reportportal.Filter, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, name, "mk-e2e-test-suite")
			return &reportportal.Filter{
				Owner: "dbizzarr",
				Share: true,
				ID:    2,
				Name:  "mk-e2e-test-suite",
				Conditions: []reportportal.FilterCondition{
					{
						FilteringField: "name",
						Condition:      "eq",
						Value:          "mk-e2e-test",
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
			}, nil, nil
		},
		UpdateM: func(projectName string, id int, f *reportportal.UpdateFilter) (string, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, id, 2)

			want := &reportportal.UpdateFilter{
				Share: true,
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
			}

			testDeepEqual(t, f, want)
			return "", nil, nil
		},
	}

	r := NewReportPortal(&reportportal.Client{
		Filter: mockFilter,
	})

	inputFilter := &Filter{
		Kind:        "Filter",
		Name:        "mk-e2e-test-suite",
		Description: "",
		Type:        "Launch",
		Conditions: []FilterCondition{
			{
				FilteringField: "name",
				Condition:      "eq",
				Value:          "mk-e2e-test-suite",
			},
		},
		Orders: []FilterOrder{
			{
				SortingColumn: "startTime",
				IsAsc:         false,
			},
			{
				SortingColumn: "number",
				IsAsc:         false,
			},
		},
	}

	err := r.ApplyFilter("test_project", inputFilter)
	if err != nil {
		t.Errorf("ReportPortal.ApplyFilter returned error: %v", err)
	}

	testDeepEqual(t, mockFilter.Counter, reportportal.MockFilterServiceCounter{GetByName: 1, Update: 1})

}

func TestApplyFilter_Skip(t *testing.T) {
	mockFilter := &reportportal.MockFilterService{
		GetByNameM: func(projectName, name string) (*reportportal.Filter, *reportportal.Response, error) {
			testEqual(t, projectName, "test_project")
			testEqual(t, name, "mk-e2e-test-suite")
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
				Type: "Launch",
			}, nil, nil
		},
	}

	r := NewReportPortal(&reportportal.Client{
		Filter: mockFilter,
	})

	inputFilter := &Filter{
		Kind:        "Filter",
		Name:        "mk-e2e-test-suite",
		Description: "",
		Type:        "Launch",
		Conditions: []FilterCondition{
			{
				FilteringField: "name",
				Condition:      "eq",
				Value:          "mk-e2e-test-suite",
			},
		},
		Orders: []FilterOrder{
			{
				SortingColumn: "startTime",
				IsAsc:         false,
			},
			{
				SortingColumn: "number",
				IsAsc:         false,
			},
		},
	}

	err := r.ApplyFilter("test_project", inputFilter)
	if err != nil {
		t.Errorf("ReportPortal.ApplyFilter returned error: %v", err)
	}

	testDeepEqual(t, mockFilter.Counter, reportportal.MockFilterServiceCounter{GetByName: 1})
}
