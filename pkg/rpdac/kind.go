package rpdac

import (
	"encoding/json"
)

type ObjectKind int

const (
	UnknownKind ObjectKind = iota
	DashboardKind
	FilterKind
)

var kinds = map[ObjectKind]string{
	DashboardKind: "Dashboard",
	FilterKind:    "Filter",
}

func (k ObjectKind) String() string {
	return kinds[k]
}

// MarshalJSON marshals the enum as a quoted json string
func (k ObjectKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.String())
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (k *ObjectKind) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	for kind, value := range kinds {
		if value == s {
			*k = kind
		}
	}
	return nil
}

func (k ObjectKind) MarshalYAML() (interface{}, error) {
	return k.String(), nil
}

func (k *ObjectKind) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	for kind, value := range kinds {
		if value == s {
			*k = kind
		}
	}
	return nil
}
