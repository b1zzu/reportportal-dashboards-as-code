package rpdac

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Object interface {
	GetName() string
	GetKind() string
}

type GenericObject struct {
	Name string `json:"name"`
	Kind string `json:"string"`
}

func ToYaml(o Object) ([]byte, error) {
	b, err := yaml.Marshal(o)
	if err != nil {
		return []byte{}, fmt.Errorf("error marshal %s %s: %w", o.GetKind(), o.GetName(), err)
	}
	return b, nil
}

func WriteToFile(o Object, file string) error {

	y, err := ToYaml(o)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(file, y, 0660)
	if err != nil {
		return fmt.Errorf("error writing yaml %s %s to file %s: %w", o.GetKind(), o.GetName(), file, err)
	}
	return nil
}

func LoadFile(file string) ([]byte, error) {

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error loading '%s': %w", file, err)
	}
	return b, nil
}

func LoadObjectFromFile(file []byte) (*GenericObject, error) {

	o := new(GenericObject)
	err := yaml.Unmarshal(file, o)
	if err != nil {
		return nil, err
	}

	if o.Kind == "" {
		log.Printf("warning: assuming kind '%s'", DashboardKind)
		o.Kind = DashboardKind
	}

	return o, nil
}
