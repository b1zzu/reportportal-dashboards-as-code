package reportportal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWidgetGet(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/test_project/widget/11", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		fmt.Fprint(w, `{
			"description": "",
			"owner": "dbizzarr",
			"share": true,
			"id": 3,
			"name": "Failed/Skipped/Passed [Last 7 days]",
			"widgetType": "statisticTrend",
			"contentParameters": {
				"contentFields": [
					"statistics$executions$passed",
					"statistics$executions$failed",
					"statistics$executions$skipped"
				],
				"itemsCount": 168,
				"widgetOptions": {
					"zoom": false,
					"timeline": "launch",
					"viewMode": "bar"
				}
			},
			"appliedFilters": [
				{
					"owner": "dbizzarr",
					"share": true,
					"id": 2,
					"name": "mk-e2e-test-suite",
					"conditions": [
						{
							"filteringField": "name",
							"condition": "eq",
							"value": "mk-e2e-test-suite"
						}
					],
					"orders": [
						{
							"sortingColumn": "startTime",
							"isAsc": false
						},
						{
							"sortingColumn": "number",
							"isAsc": false
						}
					],
					"type": "Launch"
				}
			],
			"content": {
				"result": [
					{
						"id": 38947,
						"number": 5412,
						"name": "mk-e2e-test-suite",
						"startTime": 1639145155978,
						"values": {
							"statistics$executions$passed": "147",
							"statistics$executions$skipped": "6"
						}
					},
					{
						"id": 38935,
						"number": 5411,
						"name": "mk-e2e-test-suite",
						"startTime": 1639141562776,
						"values": {
							"statistics$executions$passed": "147",
							"statistics$executions$skipped": "6"
						}
					}
				]
			}
		}`)
	})

	widget, _, err := client.Widget.Get("test_project", 11)
	if err != nil {
		t.Errorf("Widget.Get returned error: %v", err)
	}

	want := &Widget{
		Description: "",
		Owner:       "dbizzarr",
		Share:       true,
		ID:          3,
		Name:        "Failed/Skipped/Passed [Last 7 days]",
		WidgetType:  "statisticTrend",
		ContentParameters: &WidgetContentParameters{
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
		AppliedFilters: []*Filter{
			{
				Owner: "dbizzarr",
				Share: true,
				ID:    2,
				Name:  "mk-e2e-test-suite",
				Conditions: []*FilterCondition{
					{
						FilteringField: "name",
						Condition:      "eq",
						Value:          "mk-e2e-test-suite",
					},
				},
				Orders: []*FilterOrder{
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

	widget.Content = nil

	if !cmp.Equal(widget, want) {
		t.Errorf("Widget.Get returned %+v, want %+v", widget, want)
	}
}

func TestWidgetPost(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &NewWidget{
		Name:        "Failed/Skipped/Passed [Last 7 days]",
		Description: "",
		Share:       true,
		WidgetType:  "statisticTrend",
		ContentParameters: &WidgetContentParameters{
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
		Filters: []int{2},
	}

	mux.HandleFunc("/api/v1/test_project/widget", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{})

		v := new(NewWidget)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id": 1}`)
	})

	id, _, err := client.Widget.Post("test_project", input)
	if err != nil {
		t.Errorf("Widget.Post returned error: %v", err)
	}

	want := 1
	if id != want {
		t.Errorf("Widget.Post returned %+v, want %+v", id, want)
	}
}
