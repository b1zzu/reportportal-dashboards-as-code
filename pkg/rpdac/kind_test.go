package rpdac

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestObjectKindMarshalJSON(t *testing.T) {

	input := &GenericObject{
		Kind: DashboardKind,
		Name: "Name",
	}

	got, err := json.Marshal(input)
	if err != nil {
		t.Errorf("json.Marshal retunred error: %s", err)
	}

	want := "{\"kind\":\"Dashboard\",\"name\":\"Name\"}"
	testEqual(t, string(got), want)
}

func TestObjectKindMarshalYAML(t *testing.T) {

	input := &GenericObject{
		Kind: DashboardKind,
		Name: "Name",
	}

	got, err := yaml.Marshal(input)
	if err != nil {
		t.Errorf("yaml.Marshal retunred error: %s", err)
	}

	want := `kind: Dashboard
name: Name
`
	testEqual(t, string(got), want)
}

func TestObjectKindUnmarshalJSON(t *testing.T) {

	input := "{\"kind\":\"Dashboard\",\"name\":\"Name\"}"

	got := new(GenericObject)
	err := json.Unmarshal([]byte(input), got)
	if err != nil {
		t.Errorf("json.Unmarshal retunred error: %s", err)
	}

	want := &GenericObject{
		Kind: DashboardKind,
		Name: "Name",
	}
	testDeepEqual(t, got, want)
}

func TestObjectKindUnmarshalYAML(t *testing.T) {

	input := `kind: Dashboard
name: Name
`

	got := new(GenericObject)
	err := yaml.Unmarshal([]byte(input), got)
	if err != nil {
		t.Errorf("yaml.Unmarshal retunred error: %s", err)
	}

	want := &GenericObject{
		Kind: DashboardKind,
		Name: "Name",
	}

	testDeepEqual(t, got, want)
}
