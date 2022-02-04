package rpdac

type MockServiceCounter struct {
	Get       int
	GetByName int
	Create    int
	Update    int
	Delete    int
}

type MockService struct {
	GetM       func(project string, id int) (Object, error)
	GetByNameM func(project, name string) (Object, error)
	CreateM    func(project string, o Object) error
	UpdateM    func(project string, current Object, target Object) error
	DeleteM    func(project, name string) error

	Counter MockServiceCounter
}

func (s *MockService) Get(project string, id int) (Object, error) {
	s.Counter.Get++
	return s.GetM(project, id)
}
func (s *MockService) GetByName(project, name string) (Object, error) {
	s.Counter.GetByName++
	return s.GetByNameM(project, name)
}
func (s *MockService) Create(project string, o Object) error {
	s.Counter.Create++
	return s.CreateM(project, o)
}
func (s *MockService) Update(project string, current Object, target Object) error {
	s.Counter.Update++
	return s.UpdateM(project, current, target)
}
func (s *MockService) Delete(project, name string) error {
	s.Counter.Delete++
	return s.DeleteM(project, name)
}

type MockObjectCounter struct {
	Equals int
}

type MockObject struct {
	Kind ObjectKind
	Name string

	Counter MockObjectCounter
}

func (m *MockObject) GetKind() ObjectKind {
	return m.Kind
}
func (m *MockObject) GetName() string {
	return m.Name
}
func (m *MockObject) Equals(o Object) bool {
	m.Counter.Equals++
	return m.GetKind() == o.GetKind() && m.GetName() == o.GetName()
}
