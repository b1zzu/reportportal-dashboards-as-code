package reportportal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDashboardGetByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/test_project/dashboard/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		fmt.Fprint(w, `{
			"owner": "dbizzarr",
			"share": true,
			"id": 1,
			"name": "MK E2E Tests Overview",
			"widgets": [
				{
					"widgetName": "Failed/Skipped/Passed",
					"widgetId": 3,
					"widgetType": "statisticTrend",
					"widgetSize": {
						"width": 12,
						"height": 6
					},
					"widgetPosition": {
						"positionX": 0,
						"positionY": 13
					},
					"share": true
				},
				{
					"widgetName": "Unique bugs [Last 7 days]",
					"widgetId": 67,
					"widgetType": "uniqueBugTable",
					"widgetSize": {
						"width": 12,
						"height": 7
					},
					"widgetPosition": {
						"positionX": 0,
						"positionY": 44
					},
					"share": true
				}
			]
		}`)
	})

	dashboard, _, err := client.Dashboard.GetByID("test_project", 1)
	if err != nil {
		t.Errorf("Dashboard.GetByID returned error: %v", err)
	}

	want := &Dashboard{
		Owner: "dbizzarr",
		Share: true,
		ID:    1,
		Name:  "MK E2E Tests Overview",
		Widgets: []*DashboardWidget{
			{
				WidgetName: "Failed/Skipped/Passed",
				WidgetID:   3,
				WidgetType: "statisticTrend",
				WidgetSize: &DashboardWidgetSize{
					Width:  12,
					Height: 6,
				},
				WidgetPosition: &DashboardWidgetPosition{
					PositionX: 0,
					PositionY: 13,
				},
				Share: true,
			},
			{
				WidgetName: "Unique bugs [Last 7 days]",
				WidgetID:   67,
				WidgetType: "uniqueBugTable",
				WidgetSize: &DashboardWidgetSize{
					Width:  12,
					Height: 7,
				},
				WidgetPosition: &DashboardWidgetPosition{
					PositionX: 0,
					PositionY: 44,
				},
				Share: true,
			},
		},
	}

	if !cmp.Equal(dashboard, want) {
		t.Errorf("Dashboard.GetByID returned %+v, want %+v", dashboard, want)
	}
}

func TestDashboardGetByName(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/test_project/dashboard", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"filter.eq.name": "MK E2E Tests Overview",
		})
		fmt.Fprint(w, `{"content": [{
			"owner": "dbizzarr",
			"share": true,
			"id": 1,
			"name": "MK E2E Tests Overview"
		}]}`)
	})

	dashboard, _, err := client.Dashboard.GetByName("test_project", "MK E2E Tests Overview")
	if err != nil {
		t.Errorf("Dashboard.GetByName returned error: %v", err)
	}

	want := &Dashboard{
		Owner: "dbizzarr",
		Share: true,
		ID:    1,
		Name:  "MK E2E Tests Overview",
	}

	if !cmp.Equal(dashboard, want) {
		t.Errorf("Dashboard.GetByName returned %+v, want %+v", dashboard, want)
	}
}

func TestDashboardGetByName_NotFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/test_project/dashboard", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"filter.eq.name": "MK E2E Tests Overview",
		})
		fmt.Fprint(w, `{"content": []}`)
	})

	_, _, err := client.Dashboard.GetByName("test_project", "MK E2E Tests Overview")
	if _, ok := err.(*DashboardNotFoundError); !ok {
		t.Errorf("Dashboard.GetByName returned error: %v, want DashboardNotFoundError", err)
	}
}

func TestDashboardCreate(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &NewDashboard{
		Share:       true,
		Name:        "MK E2E Tests Overview",
		Description: "Test",
	}

	mux.HandleFunc("/api/v1/test_project/dashboard", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{})

		v := new(NewDashboard)
		json.NewDecoder(r.Body).Decode(v)

		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id": 1}`)
	})

	id, _, err := client.Dashboard.Create("test_project", input)
	if err != nil {
		t.Errorf("Dashboard.Create returned error: %v", err)
	}

	want := 1
	if id != want {
		t.Errorf("Dashboard.Create returned %+v, want %+v", id, want)
	}
}

func TestDashboardUpdate(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &UpdateDashboard{
		Share:       true,
		Name:        "MK E2E Tests Overview",
		Description: "Test",
	}

	mux.HandleFunc("/api/v1/test_project/dashboard/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testFormValues(t, r, values{})

		v := new(UpdateDashboard)
		json.NewDecoder(r.Body).Decode(v)

		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"message": "done"}`)
	})

	message, _, err := client.Dashboard.Update("test_project", 2, input)
	if err != nil {
		t.Errorf("Dashboard.Create returned error: %v", err)
	}

	want := "done"
	if message != want {
		t.Errorf("Dashboard.Create returned %+v, want %+v", message, want)
	}
}

func TestDashboardDelete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/test_project/dashboard/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testFormValues(t, r, values{})

		fmt.Fprint(w, `{"message": "done"}`)
	})

	message, _, err := client.Dashboard.Delete("test_project", 2)
	if err != nil {
		t.Errorf("Dashboard.Delete returned error: %v", err)
	}

	want := "done"
	if message != want {
		t.Errorf("Dashboard.Delete returned %+v, want %+v", message, want)
	}
}

func TestDashboardAddWidget(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &DashboardWidget{
		WidgetName: "Failed/Skipped/Passed",
		WidgetID:   3,
		WidgetType: "statisticTrend",
		WidgetSize: &DashboardWidgetSize{
			Width:  12,
			Height: 6,
		},
		WidgetPosition: &DashboardWidgetPosition{
			PositionX: 0,
			PositionY: 13,
		},
		Share: true,
	}

	mux.HandleFunc("/api/v1/test_project/dashboard/2/add", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testFormValues(t, r, values{})

		v := new(DashboardAddWidget)
		json.NewDecoder(r.Body).Decode(v)

		if !cmp.Equal(v.AddWidget, input) {
			t.Errorf("Request body = %+v, want %+v", v.AddWidget, input)
		}

		fmt.Fprint(w, `{"message": "done"}`)
	})

	message, _, err := client.Dashboard.AddWidget("test_project", 2, input)
	if err != nil {
		t.Errorf("Dashboard.AddWidget returned error: %v", err)
	}

	want := "done"
	if message != want {
		t.Errorf("Dashboard.AddWidget returned %+v, want %+v", message, want)
	}
}

func TestDashboardRemoveWidget(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/test_project/dashboard/2/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testFormValues(t, r, values{})

		fmt.Fprint(w, `{"message": "done"}`)
	})

	message, _, err := client.Dashboard.RemoveWidget("test_project", 2, 1)
	if err != nil {
		t.Errorf("Dashboard.RemoveWidget returned error: %v", err)
	}

	want := "done"
	if message != want {
		t.Errorf("Dashboard.RemoveWidget returned %+v, want %+v", message, want)
	}
}
