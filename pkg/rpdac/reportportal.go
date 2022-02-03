package rpdac

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/b1zzu/reportportal-dashboards-as-code/pkg/reportportal"
	"gopkg.in/yaml.v2"
)

type ReportPortal struct {
	client *reportportal.Client

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	Dashboard IDashboardService
	Filter    IFilterService
}

type Object interface {
	GetName() string
	GetKind() ObjectKind
}

type GenericObject struct {
	Kind ObjectKind `json:"kind"`
	Name string     `json:"name"`
}

type service struct {
	client *reportportal.Client
}

func NewReportPortal(c *reportportal.Client) *ReportPortal {
	r := &ReportPortal{client: c}
	r.common.client = c
	r.Dashboard = (*DashboardService)(&r.common)
	r.Filter = (*FilterService)(&r.common)
	return r
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

	if o.Kind == -1 {
		log.Printf("warning: assuming kind '%s'", DashboardKind)
		o.Kind = DashboardKind
	}

	return o, nil
}
