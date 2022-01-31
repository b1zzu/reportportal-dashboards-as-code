package reportportal

import "fmt"

type ProjectSettingsService service

type ProjectSettings struct {
	ProjectID int           `json:"project"`
	SubTypes  IssueSubTypes `json:"subTypes"`
}

type IssueSubTypes map[string][]IssueSubType

type IssueSubType struct {
	ID        int    `json:"id"`
	Locator   string `json:"locator"`
	TypeRef   string `json:"typeRef"`
	LongName  string `json:"longName"`
	ShortName string `json:"shortName"`
	Color     string `json:"color"`
}

func (s *ProjectSettingsService) Get(projectName string) (*ProjectSettings, *Response, error) {
	u := fmt.Sprintf("v1/%s/settings", projectName)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	ps := new(ProjectSettings)
	resp, err := s.client.Do(req, ps)
	if err != nil {
		return nil, resp, err
	}

	return ps, resp, nil
}
