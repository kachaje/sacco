package parser_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sacco/server/parser"
	"testing"
)

func TestUpdateSessionFlags(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("..", "database", "models", "fixtures", "sample.json"))
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]any{}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	session := parser.NewSession(nil, nil, nil)
	session.UpdateActiveData(data, 0)

	err = session.UpdateSessionFlags(nil)
	if err != nil {
		t.Fatal(err)
	}

	target := map[string]any{
		"memberBeneficiary": []map[string]any{
			{
				"memberBeneficiaryId": "1",
			},
			{
				"memberBeneficiaryId": "2",
			},
		},
		"memberContactId": "1",
		"memberId":        "1",
		"memberLoan": []map[string]any{
			{
				"memberBusinessId": "1",
				"memberLastYearBusinessHistory": []map[string]any{
					{
						"memberLastYearBusinessHistoryId": "1",
					},
				},
				"memberLoanApprovalId": "1",
				"memberLoanId":         "1",
				"memberNextYearBusinessProjection": []map[string]any{
					{
						"memberNextYearBusinessProjectionId": "1",
					},
				},
				"memberOccupationId":             "1",
				"memberOccupationVerificationId": "1",
			},
		},
		"memberNomineeId": "1",
	}

	if !reflect.DeepEqual(target, session.GlobalIds) {
		t.Fatal("Test failed")
	}
}
