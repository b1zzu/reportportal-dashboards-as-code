package rpdac

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Object interface {
	GetName() string
	GetKind() string
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
