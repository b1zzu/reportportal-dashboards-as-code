package rpdac

import (
	"reflect"
	"testing"

	"github.com/b1zzu/reportportal-dashboards-as-code/pkg/reportportal"
)

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
				Kind:        "Dashboard",
				Name:        "Test",
				Description: "My test description",
				origin:      &reportportal.Dashboard{ID: 1},
			},
			right: &Dashboard{
				Kind:        "Dashboard",
				Name:        "Test",
				Description: "My test description",
			},
			expexct: true,
		},
		{
			description: "Compare dashboards with differt names should return false",
			left: &Dashboard{
				Kind:        "Dashboard",
				Name:        "Test One",
				Description: "My test description",
			},
			right: &Dashboard{
				Kind:        "Dashboard",
				Name:        "Test",
				Description: "My test description",
			},
			expexct: false,
		},
		{
			description: "Compare dashboards with differt description should return false",
			left: &Dashboard{
				Kind:        "Dashboard",
				Name:        "Test",
				Description: "My test description",
			},
			right: &Dashboard{
				Kind:        "Dashboard",
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
				WidgetSize: &WidgetSize{Width: 1, Height: 1},
			},
			right: &Widget{
				Name:       "Test",
				WidgetSize: &WidgetSize{Width: 2, Height: 1},
			},
			expexct: false,
		},
		{
			description: "Compare widgets with differt position should return false",
			left: &Widget{
				Name:           "Test",
				WidgetPosition: &WidgetPosition{PositionX: 4, PositionY: 9},
			},
			right: &Widget{
				Name:           "Test",
				WidgetPosition: &WidgetPosition{PositionX: 7, PositionY: 9},
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
				ContentParameters: &WidgetContentParameters{
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
				ContentParameters: &WidgetContentParameters{
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
				WidgetSize:     &WidgetSize{Width: 1, Height: 1},
				WidgetPosition: &WidgetPosition{PositionX: 4, PositionY: 9},
				Filters:        []string{"one", "two"},
				ContentParameters: &WidgetContentParameters{
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
				WidgetSize:     &WidgetSize{Width: 1, Height: 1},
				WidgetPosition: &WidgetPosition{PositionX: 4, PositionY: 9},
				Filters:        []string{"two", "one"},
				ContentParameters: &WidgetContentParameters{
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
