package reportportal

type MockDashboardService struct {
	GetByNameM    func(projectName, name string) (*Dashboard, *Response, error)
	GetByIDM      func(projectName string, id int) (*Dashboard, *Response, error)
	CreateM       func(projectName string, d *NewDashboard) (int, *Response, error)
	UpdateM       func(projectName string, dashboardID int, d *UpdateDashboard) (string, *Response, error)
	DeleteM       func(projectName string, id int) (string, *Response, error)
	AddWidgetM    func(projectName string, dashboardID int, w *DashboardWidget) (string, *Response, error)
	RemoveWidgetM func(projectName string, dashboardID int, widgetID int) (string, *Response, error)
}

func (s *MockDashboardService) GetByName(projectName, name string) (*Dashboard, *Response, error) {
	return s.GetByNameM(projectName, name)
}

func (s *MockDashboardService) GetByID(projectName string, id int) (*Dashboard, *Response, error) {
	return s.GetByIDM(projectName, id)
}

func (s *MockDashboardService) Create(projectName string, d *NewDashboard) (int, *Response, error) {
	return s.CreateM(projectName, d)
}

func (s *MockDashboardService) Update(projectName string, dashboardID int, d *UpdateDashboard) (string, *Response, error) {
	return s.UpdateM(projectName, dashboardID, d)
}

func (s *MockDashboardService) Delete(projectName string, id int) (string, *Response, error) {
	return s.DeleteM(projectName, id)
}

func (s *MockDashboardService) AddWidget(projectName string, dashboardID int, w *DashboardWidget) (string, *Response, error) {
	return s.AddWidgetM(projectName, dashboardID, w)
}

func (s *MockDashboardService) RemoveWidget(projectName string, dashboardID int, widgetID int) (string, *Response, error) {
	return s.RemoveWidgetM(projectName, dashboardID, widgetID)
}

type MockWidgetService struct {
	GetM  func(projectName string, id int) (*Widget, *Response, error)
	PostM func(projectName string, w *NewWidget) (int, *Response, error)
}

func (s *MockWidgetService) Get(projectName string, id int) (*Widget, *Response, error) {
	return s.GetM(projectName, id)
}
func (s *MockWidgetService) Post(projectName string, w *NewWidget) (int, *Response, error) {
	return s.PostM(projectName, w)
}

type MockFilterService struct {
	GetByIDM   func(projectName string, id int) (*Filter, *Response, error)
	GetByNameM func(projectName, name string) (*Filter, *Response, error)
	CreateM    func(projectName string, f *NewFilter) (int, *Response, error)
	UpdateM    func(projectName string, id int, f *UpdateFilter) (string, *Response, error)
}

func (s *MockFilterService) GetByID(projectName string, id int) (*Filter, *Response, error) {
	return s.GetByIDM(projectName, id)
}
func (s *MockFilterService) GetByName(projectName, name string) (*Filter, *Response, error) {
	return s.GetByNameM(projectName, name)
}
func (s *MockFilterService) Create(projectName string, f *NewFilter) (int, *Response, error) {
	return s.CreateM(projectName, f)
}
func (s *MockFilterService) Update(projectName string, id int, f *UpdateFilter) (string, *Response, error) {
	return s.UpdateM(projectName, id, f)
}

type MockProjectSettingsService struct {
	GetM func(projectName string) (*ProjectSettings, *Response, error)
}

func (s *MockProjectSettingsService) Get(projectName string) (*ProjectSettings, *Response, error) {
	return s.GetM(projectName)
}
