package rpdac

import (
	"reflect"
	"testing"
)

func TestDecodeFieldsSubTypes(t *testing.T) {

	inputFields := []string{
		"statistics$executions$total",
		"statistics$executions$passed",
		"statistics$executions$failed",
		"statistics$executions$skipped",
		"statistics$defects$product_bug$pb001",
		"statistics$defects$automation_bug$ab001",
		"statistics$defects$system_issue$si001",
		"statistics$defects$no_defect$nd001",
		"statistics$defects$to_investigate$ti001",
		"statistics$defects$system_issue$si_1iuqflmhg6hk6",
		"statistics$defects$product_bug$pb_qdy9r7uu9q9g",
		"statistics$defects$system_issue$si_1h7o519q5xeg5",
		"statistics$defects$automation_bug$ab_uv8mlzz5fqzn",
		"statistics$defects$automation_bug$ab_1ien71b1ve81k",
		"statistics$defects$automation_bug$ab_t4f3ctreg3sl",
	}
	inputSubTypesMap := map[string]string{
		"si_1iuqflmhg6hk6": "si1",
		"si_1h7o519q5xeg5": "si2",
		"pb_qdy9r7uu9q9g":  "pb1",
		"ab_uv8mlzz5fqzn":  "ab1",
		"ab_1ien71b1ve81k": "ab2",
		"ab_t4f3ctreg3sl":  "ab3",
	}

	expectedFields := []string{
		"statistics$executions$total",
		"statistics$executions$passed",
		"statistics$executions$failed",
		"statistics$executions$skipped",
		"statistics$defects$product_bug$pb001",
		"statistics$defects$automation_bug$ab001",
		"statistics$defects$system_issue$si001",
		"statistics$defects$no_defect$nd001",
		"statistics$defects$to_investigate$ti001",
		"statistics$defects$system_issue$si1",
		"statistics$defects$product_bug$pb1",
		"statistics$defects$system_issue$si2",
		"statistics$defects$automation_bug$ab1",
		"statistics$defects$automation_bug$ab2",
		"statistics$defects$automation_bug$ab3",
	}

	result, err := DecodeFieldsSubTypes(inputFields, inputSubTypesMap)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(result, expectedFields) {
		t.Errorf("Failed: got %+v but expected %v", result, expectedFields)
	}
}
