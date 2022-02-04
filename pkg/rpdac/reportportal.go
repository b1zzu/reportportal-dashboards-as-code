package rpdac

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/b1zzu/reportportal-dashboards-as-code/pkg/reportportal"
	"gopkg.in/yaml.v2"
)

type ReportPortal struct {
	client *reportportal.Client

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	Dashboard ServiceInterface
	Filter    ServiceInterface
}

type Object interface {
	GetName() string
	GetKind() ObjectKind
	Equals(o Object) bool
}

type GenericObject struct {
	Kind ObjectKind `json:"kind"`
	Name string     `json:"name"`
}

type ServiceInterface interface {
	Get(project string, id int) (Object, error)
	GetByName(project, name string) (Object, error)
	Create(project string, o Object) error
	Update(project string, current Object, target Object) error
	Delete(project, name string) error
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

func (r *ReportPortal) Service(kind ObjectKind) (ServiceInterface, error) {
	switch kind {
	case DashboardKind:
		return r.Dashboard, nil
	case FilterKind:
		return r.Filter, nil
	default:
		return nil, fmt.Errorf("error: object kind '%s' is not supported", kind.String())
	}
}

// Export the ObjectKind with the passed id from the passed project to the passed file.
//
// The file can be a relative or absoulte path to the file that will be written with the
// full content of the exported Object.
//
func (r *ReportPortal) Export(k ObjectKind, project string, id int, file string) error {

	s, err := r.Service(k)
	if err != nil {
		return err
	}

	// retrieve object from reportportal
	o, err := s.Get(project, id)
	if err != nil {
		return fmt.Errorf("error retrieving '%s' with id '%d' in project '%s': %w", k.String(), id, project, err)
	}

	// convert object to YAML
	b, err := yaml.Marshal(o)
	if err != nil {
		return fmt.Errorf("error marshal (encoding) '%s' with id '%d' in project '%s' to YAML: %w", k.String(), id, project, err)
	}

	// write object to file
	err = ioutil.WriteFile(file, b, 0644)
	if err != nil {
		return fmt.Errorf("error writing '%s' with id '%d' in project '%s' to file '%s': %w", k.String(), id, project, file, err)
	}

	log.Printf("%s with id '%d' in project '%s' exported to '%s'", k.String(), id, project, file)
	return nil
}

// Create a object/resource in ReportPortal from the passed file.
//
func (r *ReportPortal) Create(project, file string) error {

	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("error reading file '%s': %w", file, err)
	}

	o, err := UnmarshalObject(fileBytes)
	if err != nil {
		return fmt.Errorf("error unmarshal (decoding) file '%s': %w", file, err)
	}

	s, err := r.Service(o.GetKind())
	if err != nil {
		return err
	}

	err = s.Create(project, o)
	if err != nil {
		return fmt.Errorf("error creating %s from file '%s' in project '%s': %w", o.GetKind().String(), file, project, err)
	}

	log.Printf("%s with name '%s' from file '%s' created in project '%s'", o.GetKind().String(), o.GetName(), file, project)
	return nil
}

func (r *ReportPortal) Apply(project, file string, recursive bool) error {

	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		return fmt.Errorf("error '%s' is not a vailid file or directory: %w", file, err)
	} else if err != nil {
		return err
	}

	if info.IsDir() {

		if !recursive {
			return fmt.Errorf("error '%s' is a directory, use the `-r` option if you want to recursive apply all object in the directory", file)
		}

		failed := false
		err = filepath.WalkDir(file, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				log.Printf("Unknow error: %s", err)
				return nil
			}

			if d.IsDir() {
				// skip directories
				return nil
			}

			if !strings.HasSuffix(d.Name(), ".yml") && !strings.HasSuffix(d.Name(), ".yaml") {
				log.Printf("Ignore file '%s' because only .yml|.yaml are supported", path)
				return nil
			}

			if err := r.ApplyFile(project, path); err != nil {
				failed = true
				log.Printf("Failed to apply file '%s': %s", path, err)
			}
			return nil
		})
		if err != nil {
			return err
		}

		if failed {
			return errors.New("error applying one or more objects")
		}
		return nil

	} else {
		return r.ApplyFile(project, file)
	}
}

func (r *ReportPortal) ApplyFile(project, file string) error {

	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("error reading file '%s': %w", file, err)
	}

	o, err := UnmarshalObject(fileBytes)
	if err != nil {
		return fmt.Errorf("error unmarshal (decoding) file '%s': %w", file, err)
	}

	return r.ApplyObject(project, o)
}

func (r *ReportPortal) ApplyObject(project string, o Object) error {

	s, err := r.Service(o.GetKind())
	if err != nil {
		return err
	}

	current, err := s.GetByName(project, o.GetName())
	if err != nil {
		return fmt.Errorf("error retrieving %s with name '%s': %w", o.GetKind(), o.GetName(), err)
	}

	if current != nil {

		if current.Equals(o) {
			log.Printf("Skip apply %s with name '%s' in project '%s'", o.GetKind(), o.GetName(), project)
			return nil
		}

		if err = s.Update(project, current, o); err != nil {
			return err
		}
		log.Printf("%s with name '%s' updated in project '%s'", o.GetKind(), o.GetName(), project)
		return nil
	}

	if err = s.Create(project, o); err != nil {
		return err
	}
	log.Printf("%s with name '%s' created in project '%s'", o.GetKind(), o.GetName(), project)
	return nil
}

func UnmarshalObject(file []byte) (Object, error) {

	g := new(GenericObject)
	err := yaml.Unmarshal(file, g)
	if err != nil {
		return nil, err
	}

	var o Object
	switch g.Kind {
	case DashboardKind:
		o = new(Dashboard)
	case FilterKind:
		o = new(Filter)
	case UnknownKind:
		log.Printf("warning: assuming kind '%s'", DashboardKind.String())
		o = new(Dashboard)
	default:
		return nil, fmt.Errorf("error: object kind '%s' is not suppoerted from the export method", g.Kind.String())
	}

	err = yaml.Unmarshal(file, o)
	if err != nil {
		return nil, err
	}
	return o, nil
}
