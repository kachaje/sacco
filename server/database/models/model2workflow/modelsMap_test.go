package model2workflow_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sacco/server/parser"
	"testing"
)

func TestModelsMap(t *testing.T) {
	sample := map[string]any{}

	content, err := os.ReadFile(filepath.Join("..", "fixtures", "sample.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &sample)
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]string{}

	content, err = os.ReadFile(filepath.Join(".", "modelsMap.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	results := map[string]any{}

	session := parser.NewSession(nil, nil, nil)

	for key, value := range data {
		result, ok := session.DecodeKey(value, sample)
		if ok {
			results[key] = result
		} else {
			fmt.Println("missed:", key)
		}
	}

	target := map[string]any{
		"memberBeneficiaryId":                1.0,
		"memberBusinessId":                   1.0,
		"memberContactId":                    1.0,
		"memberId":                           1.0,
		"memberLastYearBusinessHistoryId":    1.0,
		"memberLoanApprovalId":               1.0,
		"memberLoanId":                       1.0,
		"memberLoanLiabilityId":              1.0,
		"memberLoanSecurityId":               1.0,
		"memberLoanWitnessId":                1.0,
		"memberNextYearBusinessProjectionId": 1.0,
		"memberNomineeId":                    1.0,
		"memberOccupationId":                 1.0,
		"memberOccupationVerificationId":     1.0,
		"memberSharesId":                     1.0,
	}

	if !reflect.DeepEqual(results, target) {
		t.Fatal("Test failed")
	}
}
