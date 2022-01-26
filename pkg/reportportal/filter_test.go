package reportportal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFilterGetByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/test_project/filter/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		fmt.Fprint(w, `{
			"owner": "dbizzarr",
			"share": true,
			"id": 2,
			"name": "mk-e2e-test-suite",
			"description": "test",
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
		}`)
	})

	filter, _, err := client.Filter.GetByID("test_project", 2)
	if err != nil {
		t.Errorf("Filter.GetByID returned error: %v", err)
	}

	want := &Filter{
		Owner:       "dbizzarr",
		Share:       true,
		ID:          2,
		Name:        "mk-e2e-test-suite",
		Description: "test",
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
	}

	if !cmp.Equal(filter, want) {
		t.Errorf("Filter.GetByID returned %+v, want %+v", filter, want)
	}
}

func TestFilterGetByName(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/test_project/filter", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"filter.eq.name": "mk-e2e-test-suite",
		})
		fmt.Fprint(w, `{"content": [{
			"id": 2,
			"name": "mk-e2e-test-suite"
		}]}`)
	})

	filter, _, err := client.Filter.GetByName("test_project", "mk-e2e-test-suite")
	if err != nil {
		t.Errorf("Filter.GetByName returned error: %v", err)
	}

	want := &Filter{
		ID:   2,
		Name: "mk-e2e-test-suite",
	}

	if !cmp.Equal(filter, want) {
		t.Errorf("Filter.GetByName returned %+v, want %+v", filter, want)
	}
}

func TestFilterGetByName_NotFound(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/test_project/filter", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"filter.eq.name": "mk-e2e-test-suite",
		})
		fmt.Fprint(w, `{"content": []}`)
	})

	_, _, err := client.Filter.GetByName("test_project", "mk-e2e-test-suite")
	if _, ok := err.(*FilterNotFoundError); !ok {
		t.Errorf("Filter.GetByName returned error: %v, want FilterNotFoundError", err)
	}
}

func TestFilterCreate(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &NewFilter{
		Share:       true,
		Name:        "mk-e2e-test-suite",
		Description: "test",
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
	}

	mux.HandleFunc("/api/v1/test_project/filter", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, values{})

		v := new(NewFilter)
		json.NewDecoder(r.Body).Decode(v)

		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id": 1}`)
	})

	id, _, err := client.Filter.Create("test_project", input)
	if err != nil {
		t.Errorf("Filter.Create returned error: %v", err)
	}

	want := 1
	if id != want {
		t.Errorf("Filter.Create returned %+v, want %+v", id, want)
	}
}

func TestFilterUpdate(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &UpdateFilter{
		Share:       true,
		Name:        "mk-e2e-test-suite",
		Description: "test",
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
	}

	mux.HandleFunc("/api/v1/test_project/filter/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testFormValues(t, r, values{})

		v := new(UpdateFilter)
		json.NewDecoder(r.Body).Decode(v)

		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"message": "done"}`)
	})

	message, _, err := client.Filter.Update("test_project", 2, input)
	if err != nil {
		t.Errorf("Filter.Create returned error: %v", err)
	}

	want := "done"
	if message != want {
		t.Errorf("Filter.Create returned %+v, want %+v", message, want)
	}
}
