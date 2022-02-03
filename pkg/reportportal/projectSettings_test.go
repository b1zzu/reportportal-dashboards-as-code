package reportportal

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPojectSettingsGet(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/test_project/settings", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		fmt.Fprint(w, `{
				"project": 4,
				"subTypes": {
					"NO_DEFECT": [
						{
							"id": 4,
							"locator": "nd001",
							"typeRef": "NO_DEFECT",
							"longName": "No Defect",
							"shortName": "ND",
							"color": "#777777"
						}
					],
					"AUTOMATION_BUG": [
						{
							"id": 2,
							"locator": "ab001",
							"typeRef": "AUTOMATION_BUG",
							"longName": "Automation Bug",
							"shortName": "AB",
							"color": "#f7d63e"
						},
						{
							"id": 18,
							"locator": "ab_uv8mlzz5fqzn",
							"typeRef": "AUTOMATION_BUG",
							"longName": "Product Breaking Change",
							"shortName": "PBC",
							"color": "#f50057"
						}
					]
				}
			}`)
	})

	projectSettings, _, err := client.ProjectSettings.Get("test_project")
	if err != nil {
		t.Errorf("ProjectSettings.Get returned error: %v", err)
	}

	want := &ProjectSettings{
		ProjectID: 4,
		SubTypes: IssueSubTypes{
			"NO_DEFECT": []IssueSubType{{
				ID:        4,
				Locator:   "nd001",
				TypeRef:   "NO_DEFECT",
				LongName:  "No Defect",
				ShortName: "ND",
				Color:     "#777777",
			}},
			"AUTOMATION_BUG": []IssueSubType{{
				ID:        2,
				Locator:   "ab001",
				TypeRef:   "AUTOMATION_BUG",
				LongName:  "Automation Bug",
				ShortName: "AB",
				Color:     "#f7d63e",
			}, {
				ID:        18,
				Locator:   "ab_uv8mlzz5fqzn",
				TypeRef:   "AUTOMATION_BUG",
				LongName:  "Product Breaking Change",
				ShortName: "PBC",
				Color:     "#f50057",
			}},
		},
	}

	if !cmp.Equal(projectSettings, want) {
		t.Errorf("ProjectSettings.Get returned %+v, want %+v", projectSettings, want)
	}
}
