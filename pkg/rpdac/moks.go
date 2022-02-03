package rpdac

type MockDashboardService struct {
	GetDashboardM       func(project string, id int) (*Dashboard, error)
	GetDashboardByNameM func(project, name string) (*Dashboard, error)
	CreateDashboardM    func(project string, d *Dashboard) error
	ApplyDashboardM     func(project string, d *Dashboard) error
	DeleteDashboardM    func(project, name string) error
}

func (s *MockDashboardService) Get(project string, id int) (Object, error) {
	return s.GetDashboard(project, id)
}
func (s *MockDashboardService) GetDashboard(project string, id int) (*Dashboard, error) {
	return s.GetDashboardM(project, id)
}
func (s *MockDashboardService) GetDashboardByName(project, name string) (*Dashboard, error) {
	return s.GetDashboardByNameM(project, name)
}
func (s *MockDashboardService) CreateDashboard(project string, d *Dashboard) error {
	return s.CreateDashboardM(project, d)
}
func (s *MockDashboardService) ApplyDashboard(project string, d *Dashboard) error {
	return s.ApplyDashboardM(project, d)
}
func (s *MockDashboardService) DeleteDashboard(project, name string) error {
	return s.DeleteDashboardM(project, name)
}

type MockFilterService struct {
	GetFilterM       func(project string, id int) (*Filter, error)
	GetFilterByNameM func(project, name string) (*Filter, error)
	CreateFilterM    func(project string, f *Filter) error
	ApplyFilterM     func(project string, f *Filter) error
}

func (s *MockFilterService) Get(project string, id int) (Object, error) {
	return s.GetFilter(project, id)
}
func (s *MockFilterService) GetFilter(project string, id int) (*Filter, error) {
	return s.GetFilterM(project, id)
}
func (s *MockFilterService) GetFilterByName(project, name string) (*Filter, error) {
	return s.GetFilterByNameM(project, name)
}
func (s *MockFilterService) CreateFilter(project string, f *Filter) error {
	return s.CreateFilterM(project, f)
}
func (s *MockFilterService) ApplyFilter(project string, f *Filter) error {
	return s.ApplyFilterM(project, f)
}
