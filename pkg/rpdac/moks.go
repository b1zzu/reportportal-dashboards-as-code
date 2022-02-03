package rpdac

type MockDashboardServiceCounter struct {
	GetDashboard       int
	GetDashboardByName int
	CreateDashboard    int
	ApplyDashboard     int
	DeleteDashboard    int
}

type MockDashboardService struct {
	GetDashboardM       func(project string, id int) (*Dashboard, error)
	GetDashboardByNameM func(project, name string) (*Dashboard, error)
	CreateDashboardM    func(project string, d *Dashboard) error
	ApplyDashboardM     func(project string, d *Dashboard) error
	DeleteDashboardM    func(project, name string) error

	Counter MockDashboardServiceCounter
}

func (s *MockDashboardService) Get(project string, id int) (Object, error) {
	return s.GetDashboard(project, id)
}
func (s *MockDashboardService) Create(project string, o Object) error {
	return s.CreateDashboard(project, o.(*Dashboard))
}
func (s *MockDashboardService) Apply(project string, o Object) error {
	return s.ApplyDashboard(project, o.(*Dashboard))
}
func (s *MockDashboardService) GetDashboard(project string, id int) (*Dashboard, error) {
	s.Counter.GetDashboard++
	return s.GetDashboardM(project, id)
}
func (s *MockDashboardService) GetDashboardByName(project, name string) (*Dashboard, error) {
	s.Counter.GetDashboardByName++
	return s.GetDashboardByNameM(project, name)
}
func (s *MockDashboardService) CreateDashboard(project string, d *Dashboard) error {
	s.Counter.CreateDashboard++
	return s.CreateDashboardM(project, d)
}
func (s *MockDashboardService) ApplyDashboard(project string, d *Dashboard) error {
	s.Counter.ApplyDashboard++
	return s.ApplyDashboardM(project, d)
}
func (s *MockDashboardService) DeleteDashboard(project, name string) error {
	s.Counter.DeleteDashboard++
	return s.DeleteDashboardM(project, name)
}

type MockFilterServiceCounter struct {
	GetFilter       int
	GetFilterByName int
	CreateFilter    int
	ApplyFilter     int
}

type MockFilterService struct {
	GetFilterM       func(project string, id int) (*Filter, error)
	GetFilterByNameM func(project, name string) (*Filter, error)
	CreateFilterM    func(project string, f *Filter) error
	ApplyFilterM     func(project string, f *Filter) error

	Counter MockFilterServiceCounter
}

func (s *MockFilterService) Get(project string, id int) (Object, error) {
	return s.GetFilter(project, id)
}
func (s *MockFilterService) Create(project string, o Object) error {
	return s.CreateFilter(project, o.(*Filter))
}
func (s *MockFilterService) Apply(project string, o Object) error {
	return s.ApplyFilter(project, o.(*Filter))
}
func (s *MockFilterService) GetFilter(project string, id int) (*Filter, error) {
	s.Counter.GetFilter++
	return s.GetFilterM(project, id)
}
func (s *MockFilterService) GetFilterByName(project, name string) (*Filter, error) {
	s.Counter.GetFilterByName++
	return s.GetFilterByNameM(project, name)
}
func (s *MockFilterService) CreateFilter(project string, f *Filter) error {
	s.Counter.CreateFilter++
	return s.CreateFilterM(project, f)
}
func (s *MockFilterService) ApplyFilter(project string, f *Filter) error {
	s.Counter.ApplyFilter++
	return s.ApplyFilterM(project, f)
}
