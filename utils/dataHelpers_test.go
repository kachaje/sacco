package utils_test

import (
	"os"
	"path/filepath"
	"reflect"
	"sacco/utils"
	"sort"
	"testing"
)

func TestFlattenMapIdMapOnly(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("..", "server", "database", "models", "fixtures", "sample.json"))
	if err != nil {
		t.Fatal(err)
	}

	data, err := utils.LoadYaml(string(content))
	if err != nil {
		t.Fatal(err)
	}

	result := utils.FlattenMap(data, true)

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

func TestFlattenMapAllData(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("..", "server", "database", "models", "fixtures", "sample.json"))
	if err != nil {
		t.Fatal(err)
	}

	data, err := utils.LoadYaml(string(content))
	if err != nil {
		t.Fatal(err)
	}

	result := utils.FlattenMap(data, false)

	target := map[string]any{
		"member.dateOfBirth":                                "1999-09-01",
		"member.fileNumber":                                 "",
		"member.firstName":                                  "Mary",
		"member.gender":                                     "Female",
		"member.id":                                         1,
		"member.lastName":                                   "Banda",
		"member.maritalStatus":                              "Single",
		"member.memberBeneficiary.0.contact":                "P.O. Box 1",
		"member.memberBeneficiary.0.id":                     1,
		"member.memberBeneficiary.0.memberId":               1,
		"member.memberBeneficiary.0.name":                   "Benefator 1",
		"member.memberBeneficiary.0.percentage":             10,
		"member.memberContact.email":                        any(nil),
		"member.memberContact.homeDistrict":                 "Lilongwe",
		"member.memberContact.homeTraditionalAuthority":     "Kabudula",
		"member.memberContact.homeVillage":                  "Thumba",
		"member.memberContact.id":                           1,
		"member.memberContact.memberId":                     1,
		"member.memberContact.postalAddress":                "P.O. Box 3200, Blantyre",
		"member.memberContact.residentialAddress":           "Chilomoni, Blantrye",
		"member.memberLoan.0.id":                            1,
		"member.memberLoan.0.loanAmount":                    200000,
		"member.memberLoan.0.loanPurpose":                   "School fees",
		"member.memberLoan.0.loanType":                      "PERSONAL",
		"member.memberLoan.0.memberBusiness.businessName":   "Vendors Galore",
		"member.memberLoan.0.memberBusiness.businessNature": "Vendor",
		"member.memberLoan.0.memberBusiness.id":             1,
		"member.memberLoan.0.memberBusiness.memberLastYearBusinessHistory.0.employeesWages":      50000,
		"member.memberLoan.0.memberBusiness.memberLastYearBusinessHistory.0.financialYear":       "2024",
		"member.memberLoan.0.memberBusiness.memberLastYearBusinessHistory.0.id":                  1,
		"member.memberLoan.0.memberBusiness.memberLastYearBusinessHistory.0.loanInterest":        0,
		"member.memberLoan.0.memberBusiness.memberLastYearBusinessHistory.0.memberBusinessId":    1,
		"member.memberLoan.0.memberBusiness.memberLastYearBusinessHistory.0.netProfitLoss":       715000,
		"member.memberLoan.0.memberBusiness.memberLastYearBusinessHistory.0.otherCosts":          0,
		"member.memberLoan.0.memberBusiness.memberLastYearBusinessHistory.0.ownSalary":           100000,
		"member.memberLoan.0.memberBusiness.memberLastYearBusinessHistory.0.rentals":             50000,
		"member.memberLoan.0.memberBusiness.memberLastYearBusinessHistory.0.totalCostOfGoods":    1000000,
		"member.memberLoan.0.memberBusiness.memberLastYearBusinessHistory.0.totalCosts":          1285000,
		"member.memberLoan.0.memberBusiness.memberLastYearBusinessHistory.0.totalIncome":         2000000,
		"member.memberLoan.0.memberBusiness.memberLastYearBusinessHistory.0.transport":           50000,
		"member.memberLoan.0.memberBusiness.memberLastYearBusinessHistory.0.utilities":           35000,
		"member.memberLoan.0.memberBusiness.memberLoanId":                                        1,
		"member.memberLoan.0.memberBusiness.memberNextYearBusinessProjection.0.employeesWages":   50000,
		"member.memberLoan.0.memberBusiness.memberNextYearBusinessProjection.0.financialYear":    "2025",
		"member.memberLoan.0.memberBusiness.memberNextYearBusinessProjection.0.id":               1,
		"member.memberLoan.0.memberBusiness.memberNextYearBusinessProjection.0.loanInterest":     0,
		"member.memberLoan.0.memberBusiness.memberNextYearBusinessProjection.0.memberBusinessId": 1,
		"member.memberLoan.0.memberBusiness.memberNextYearBusinessProjection.0.netProfitLoss":    715000,
		"member.memberLoan.0.memberBusiness.memberNextYearBusinessProjection.0.otherCosts":       0,
		"member.memberLoan.0.memberBusiness.memberNextYearBusinessProjection.0.ownSalary":        100000,
		"member.memberLoan.0.memberBusiness.memberNextYearBusinessProjection.0.rentals":          50000,
		"member.memberLoan.0.memberBusiness.memberNextYearBusinessProjection.0.totalCostOfGoods": 1500000,
		"member.memberLoan.0.memberBusiness.memberNextYearBusinessProjection.0.totalCosts":       1285000,
		"member.memberLoan.0.memberBusiness.memberNextYearBusinessProjection.0.totalIncome":      2500000,
		"member.memberLoan.0.memberBusiness.memberNextYearBusinessProjection.0.transport":        50000,
		"member.memberLoan.0.memberBusiness.memberNextYearBusinessProjection.0.utilities":        35000,
		"member.memberLoan.0.memberBusiness.tradingArea":                                         "Mtandire",
		"member.memberLoan.0.memberBusiness.yearsInBusiness":                                     3,
		"member.memberLoan.0.memberId":                                                         1,
		"member.memberLoan.0.memberLoanApproval.amountApproved":                                200000,
		"member.memberLoan.0.memberLoanApproval.amountRecommended":                             200000,
		"member.memberLoan.0.memberLoanApproval.approvalDate":                                  "2025-08-30",
		"member.memberLoan.0.memberLoanApproval.approvedBy":                                    "me",
		"member.memberLoan.0.memberLoanApproval.dateVerified":                                  "2025-08-30",
		"member.memberLoan.0.memberLoanApproval.denialOrPartialReason":                         any(nil),
		"member.memberLoan.0.memberLoanApproval.id":                                            1,
		"member.memberLoan.0.memberLoanApproval.loanStatus":                                    "APPROVED",
		"member.memberLoan.0.memberLoanApproval.memberLoanId":                                  1,
		"member.memberLoan.0.memberLoanApproval.verifiedBy":                                    "me",
		"member.memberLoan.0.memberOccupation.employerAddress":                                 "Kanengo",
		"member.memberLoan.0.memberOccupation.employerName":                                    "SOBO",
		"member.memberLoan.0.memberOccupation.employerPhone":                                   "0999888474",
		"member.memberLoan.0.memberOccupation.grossPay":                                        100000,
		"member.memberLoan.0.memberOccupation.highestQualification":                            "Secondary",
		"member.memberLoan.0.memberOccupation.id":                                              1,
		"member.memberLoan.0.memberOccupation.jobTitle":                                        "Driver",
		"member.memberLoan.0.memberOccupation.memberLoanId":                                    1,
		"member.memberLoan.0.memberOccupation.memberOccupationVerification.grossVerified":      "Yes",
		"member.memberLoan.0.memberOccupation.memberOccupationVerification.id":                 1,
		"member.memberLoan.0.memberOccupation.memberOccupationVerification.jobVerified":        "Yes",
		"member.memberLoan.0.memberOccupation.memberOccupationVerification.memberOccupationId": 1,
		"member.memberLoan.0.memberOccupation.memberOccupationVerification.netVerified":        "Yes",
		"member.memberLoan.0.memberOccupation.netPay":                                          90000,
		"member.memberLoan.0.memberOccupation.periodEmployedInMonths":                          36,
		"member.memberLoan.0.phoneNumber":                                                      any(nil),
		"member.memberLoan.0.repaymentPeriodInMonths":                                          12,
		"member.memberNominee.address":                                                         "Same as member",
		"member.memberNominee.id":                                                              1,
		"member.memberNominee.memberId":                                                        1,
		"member.memberNominee.name":                                                            "John Banda",
		"member.memberNominee.phoneNumber":                                                     "0888888888",
		"member.nationalId":                                                                    "DHFYR8475",
		"member.oldFileNumber":                                                                 "",
		"member.otherName":                                                                     "",
		"member.phoneNumber":                                                                   "09999999999",
		"member.title":                                                                         "Miss",
		"member.utilityBillNumber":                                                             "29383746",
		"member.utilityBillType":                                                               "ESCOM",
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
