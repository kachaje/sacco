package parser_test

import (
	"reflect"
	"sacco/server/parser"
	"testing"
)

func TestGetTokens(t *testing.T) {
	target := map[string]any{
		"op": "SUM",
		"terms": []any{
			"totalCostOfGoods",
			"employeesWages",
			"ownSalary",
			"transport",
			"loanInterest",
			"utilities",
			"rentals",
			"otherCosts",
		},
	}

	result := parser.GetTokens("SUM({{totalCostOfGoods}}, {{employeesWages}}, {{ownSalary}}, {{transport}}, {{loanInterest}}, {{utilities}}, {{rentals}}, {{otherCosts}})")

	if reflect.DeepEqual(target, result) {
		t.Fatal("Test failed")
	}
}

func TestCalculateFormulae(t *testing.T) {
	wf := parser.NewWorkflow(data, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	wf.Data = map[string]any{
		"totalCostOfGoods": "1000000",
		"employeesWages":   "500000",
		"ownSalary":        "100000",
		"transport":        "50000",
		"loanInterest":     "0",
		"utilities":        "35000",
		"rentals":          "50000",
		"otherCosts":       "0",
	}

	wf.FormulaFields["totalCosts"] = "SUM({{totalCostOfGoods}}, {{employeesWages}}, {{ownSalary}}, {{transport}}, {{loanInterest}}, {{utilities}}, {{rentals}}, {{otherCosts}})"
	wf.FormulaFields["netProfitLoss"] = "DIFF({{totalIncome}},{{totalCosts}})"

	wf.CalculateFormulae()
}
