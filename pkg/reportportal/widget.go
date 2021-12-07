package reportportal

import (
	"fmt"
)

type WidgetService service

type Widget struct {
	Description       string                   `json:"description"`
	Owner             string                   `json:"owner"`
	Share             bool                     `json:"share"`
	ID                int                      `json:"id"`
	Name              string                   `json:"name"`
	WidgetType        string                   `json:"widgetType"`
	ContentParameters *WidgetContentParameters `json:"contentParameters"`
	AppliedFilters    []*WidgetAppliedFilter   `json:"appliedFilters"`
	Content           interface{}              `json:"content"` // incomplete
}

type WidgetContentParameters struct {
	ContentFields []string               `json:"contentFields"`
	ItemsCount    int                    `json:"itemsCount"`
	WidgetOptions map[string]interface{} `json:"widgetOptions"`
}

type WidgetAppliedFilter struct {
	Share bool   `json:"share"`
	ID    int    `json:"id"`
	Name  string `json:"name"`
	// incomplete
}

func (s *WidgetService) Get(projectName string, id int) (*Widget, *Response, error) {
	u := fmt.Sprintf("v1/%v/widget/%v", projectName, id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	w := new(Widget)
	resp, err := s.client.Do(req, w)
	if err != nil {
		return nil, resp, err
	}

	return w, resp, nil
}

// {
//     "description": "",
//     "owner": "dbizzarr",
//     "share": true,
//     "id": 71,
//     "name": "Overall statistics [Last 7 days]",
//     "widgetType": "overallStatistics",
//     "contentParameters": {
//         "contentFields": [
//             "statistics$executions$total",
//             "statistics$executions$passed",
//             "statistics$executions$failed",
//             "statistics$executions$skipped",
//             "statistics$defects$product_bug$pb001",
//             "statistics$defects$product_bug$pb_qdy9r7uu9q9g",
//             "statistics$defects$automation_bug$ab001",
//             "statistics$defects$system_issue$si001",
//             "statistics$defects$system_issue$si_1iuqflmhg6hk6",
//             "statistics$defects$no_defect$nd001",
//             "statistics$defects$to_investigate$ti001",
//             "statistics$defects$system_issue$si_1h7o519q5xeg5",
//             "statistics$defects$automation_bug$ab_uv8mlzz5fqzn",
//             "statistics$defects$automation_bug$ab_1ien71b1ve81k",
//             "statistics$defects$automation_bug$ab_t4f3ctreg3sl"
//         ],
//         "itemsCount": 168,
//         "widgetOptions": {
//             "latest": false,
//             "viewMode": "panel"
//         }
//     },
//     "appliedFilters": [
//         {
//             "owner": "dbizzarr",
//             "share": true,
//             "id": 2,
//             "name": "mk-e2e-test-suite",
//             "conditions": [
//                 {
//                     "filteringField": "name",
//                     "condition": "eq",
//                     "value": "mk-e2e-test-suite"
//                 }
//             ],
//             "orders": [
//                 {
//                     "sortingColumn": "startTime",
//                     "isAsc": false
//                 },
//                 {
//                     "sortingColumn": "number",
//                     "isAsc": false
//                 }
//             ],
//             "type": "Launch"
//         }
//     ],
//     "content": {
//         "result": [
//             {
//                 "values": {
//                     "statistics$defects$system_issue$si001": 3,
//                     "statistics$defects$system_issue$si_1iuqflmhg6hk6": 0,
//                     "statistics$executions$passed": 22950,
//                     "statistics$executions$skipped": 2483,
//                     "statistics$defects$to_investigate$ti001": 53,
//                     "statistics$defects$automation_bug$ab001": 138,
//                     "statistics$defects$product_bug$pb001": 3,
//                     "statistics$executions$failed": 133,
//                     "statistics$executions$total": 25566
//                 }
//             }
//         ]
//     }
// }
