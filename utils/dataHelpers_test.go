package utils_test

import (
	"os"
	"path/filepath"
	"reflect"
	"sacco/utils"
	"sort"
	"testing"
)

func TestFlattenMap(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("..", "server", "database", "models", "fixtures", "sample.json"))
	if err != nil {
		t.Fatal(err)
	}

	data, err := utils.LoadYaml(string(content))
	if err != nil {
		t.Fatal(err)
	}

	result := utils.FlattenMap(data)

	target := map[string]any{
		"memberBeneficiaryId":                "member.memberBeneficiary.0.id",
		"memberBusinessId":                   "member.memberLoan.0.memberBusiness.id",
		"memberContactId":                    "member.memberContact.id",
		"memberId":                           "member.id",
		"memberLastYearBusinessHistoryId":    "member.memberLoan.0.memberBusiness.memberLastYearBusinessHistory.0.id",
		"memberLoanApprovalId":               "member.memberLoan.0.memberLoanApproval.id",
		"memberLoanId":                       "member.memberLoan.0.id",
		"memberNextYearBusinessProjectionId": "member.memberLoan.0.memberBusiness.memberNextYearBusinessProjection.0.id",
		"memberNomineeId":                    "member.memberNominee.id",
		"memberOccupationId":                 "member.memberLoan.0.memberOccupation.id",
		"memberOccupationVerificationId":     "member.memberLoan.0.memberOccupation.memberOccupationVerification.id",
	}

	if !reflect.DeepEqual(target, result) {
		t.Fatal("Test failed")
	}
}

func TestSetNestedValue(t *testing.T) {
	rawData := map[string]any{
		"memberBeneficiaryId":                "member.memberBeneficiary.0.id",
		"memberBusinessId":                   "member.memberLoan.0.memberBusiness.id",
		"memberContactId":                    "member.memberContact.id",
		"memberId":                           "member.id",
		"memberLastYearBusinessHistoryId":    "member.memberLoan.0.memberBusiness.memberLastYearBusinessHistory.0.id",
		"memberLoanApprovalId":               "member.memberLoan.0.memberLoanApproval.id",
		"memberLoanId":                       "member.memberLoan.0.id",
		"memberNextYearBusinessProjectionId": "member.memberLoan.0.memberBusiness.memberNextYearBusinessProjection.0.id",
		"memberNomineeId":                    "member.memberNominee.id",
		"memberOccupationId":                 "member.memberLoan.0.memberOccupation.id",
		"memberOccupationVerificationId":     "member.memberLoan.0.memberOccupation.memberOccupationVerification.id",
	}

	keys := []string{}

	for key := range rawData {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	data := map[string]any{}

	for i, key := range keys {
		value := rawData[key]

		utils.SetNestedValue(data, value.(string), i+1)
	}

	target := map[string]any{
		"member": map[string]any{
			"id": 4,
			"memberBeneficiary": map[string]any{
				"0": map[string]any{
					"id": 1,
				},
			},
			"memberContact": map[string]any{
				"id": 3,
			},
			"memberLoan": map[string]any{
				"0": map[string]any{
					"id": 7,
					"memberBusiness": map[string]any{
						"id": 2,
						"memberLastYearBusinessHistory": map[string]any{
							"0": map[string]any{
								"id": 5,
							},
						},
						"memberNextYearBusinessProjection": map[string]any{
							"0": map[string]any{
								"id": 8,
							},
						},
					},
					"memberLoanApproval": map[string]any{
						"id": 6,
					},
					"memberOccupation": map[string]any{
						"id": 10,
						"memberOccupationVerification": map[string]any{
							"id": 11,
						},
					},
				},
			},
			"memberNominee": map[string]any{
				"id": 9,
			},
		},
	}

	if !reflect.DeepEqual(target, data) {
		t.Fatal("Test failed")
	}
}

func TestDecodeKey(t *testing.T) {
	data := map[string]any{
		"child": map[string]any{
			"child1Id":     "2",
			"child1_1Id":   "3",
			"child1_1_1Id": "4",
			"child1_1_2": []map[string]any{
				{
					"child1_1_2Id": "5",
				},
				{
					"child1_1_2Id": "6",
				},
			},
			"child1_1_3Id":   "7",
			"child1_1_3_1Id": "8",
			"child2": []map[string]any{
				{
					"child2Id":   "9",
					"child2_1Id": "10",
					"child2_1_1": []map[string]any{
						{
							"child2_1_1Id": "11",
						},
						{
							"child2_1_1Id": "12",
						},
					},
				},
				{
					"child2Id": "13",
				},
			},
			"id": "1",
		},
	}
	target := map[string]any{
		"childId":                                  "1",
		"child.child1Id":                           "2",
		"child.child1_1Id":                         "3",
		"child.child1_1_1Id":                       "4",
		"child.child1_1_2.0.child1_1_2Id":          "5",
		"child.child1_1_2.1.child1_1_2Id":          "6",
		"child.child1_1_3Id":                       "7",
		"child.child1_1_3_1Id":                     "8",
		"child.child2.0.child2Id":                  "9",
		"child.child2.0.child2_1Id":                "10",
		"child.child2.0.child2_1_1.0.child2_1_1Id": "11",
		"child.child2.0.child2_1_1.1.child2_1_1Id": "12",
		"child.child2.1.child2Id":                  "13",
	}

	for key, value := range target {
		result, ok := utils.DecodeKey(key, data)
		if ok {
			if result != value {
				t.Fatalf("Test failed. Expecting: %v; Actual: %v", value, result)
			}
		} else {
			t.Fatalf("Test failed. Failed to fetch %s", key)
		}
	}
}

func TestFlattenKeys(t *testing.T) {
	data := map[string]any{
		"child": map[string]any{
			"child1Id":     "2",
			"child1_1Id":   "3",
			"child1_1_1Id": "4",
			"child1_1_2": []map[string]any{
				{
					"child1_1_2Id": "5",
				},
				{
					"child1_1_2Id": "6",
				},
			},
			"child1_1_3Id":   "7",
			"child1_1_3_1Id": "8",
			"child2": []map[string]any{
				{
					"child2Id":   "9",
					"child2_1Id": "10",
					"child2_1_1": []map[string]any{
						{
							"child2_1_1Id": "11",
						},
						{
							"child2_1_1Id": "12",
						},
					},
				},
				{
					"child2Id": "13",
				},
			},
			"id": "1",
		},
	}

	target := map[string]any{
		"childId":                                  "1",
		"child.child1Id":                           "2",
		"child.child1_1Id":                         "3",
		"child.child1_1_1Id":                       "4",
		"child.child1_1_2.0.child1_1_2Id":          "5",
		"child.child1_1_2.1.child1_1_2Id":          "6",
		"child.child1_1_3Id":                       "7",
		"child.child1_1_3_1Id":                     "8",
		"child.child2.0.child2Id":                  "9",
		"child.child2.0.child2_1Id":                "10",
		"child.child2.0.child2_1_1.0.child2_1_1Id": "11",
		"child.child2.0.child2_1_1.1.child2_1_1Id": "12",
		"child.child2.1.child2Id":                  "13",
	}

	result := utils.FlattenKeys(data, map[string]any{}, nil)

	if !reflect.DeepEqual(target, result) {
		t.Fatal("Test failed")
	}
}

func TestLoadKeys(t *testing.T) {
	data := map[string]any{
		"id":    "1",
		"field": "value",
		"child1": map[string]any{
			"id":    "2",
			"field": "value",
			"child1_1": map[string]any{
				"id":    "3",
				"field": "value",
				"child1_1_1": map[string]any{
					"id":    "4",
					"field": "value",
				},
				"child1_1_2": []map[string]any{
					{
						"id":    "5",
						"field": "value",
					},
					{
						"id":    "6",
						"field": "value",
					},
				},
				"child1_1_3": map[string]any{
					"id":    "7",
					"field": "value",
					"child1_1_3_1": map[string]any{
						"id":    "8",
						"field": "value",
					},
				},
			},
		},
		"child2": []map[string]any{
			{
				"id":    "9",
				"field": "value",
				"child2_1": map[string]any{
					"id":    "10",
					"field": "value",
					"child2_1_1": []map[string]any{
						{
							"id":    "11",
							"field": "value",
						},
						{
							"id":    "12",
							"field": "value",
						},
					},
				},
			},
			{
				"id":    "13",
				"field": "value",
			},
		},
	}
	target := map[string]any{
		"child1Id":     "2",
		"child1_1Id":   "3",
		"child1_1_1Id": "4",
		"child1_1_2": []map[string]any{
			{
				"child1_1_2Id": "5",
			},
			{
				"child1_1_2Id": "6",
			},
		},
		"child1_1_3Id":   "7",
		"child1_1_3_1Id": "8",
		"child2": []map[string]any{
			{
				"child2Id":   "9",
				"child2_1Id": "10",
				"child2_1_1": []map[string]any{
					{
						"child2_1_1Id": "11",
					},
					{
						"child2_1_1Id": "12",
					},
				},
			},
			{
				"child2Id": "13",
			},
		},
		"id": "1",
	}

	result := utils.LoadKeys(data, map[string]any{}, nil)

	if !reflect.DeepEqual(target, result) {
		t.Fatal("Test failed")
	}
}
