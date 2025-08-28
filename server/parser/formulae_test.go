package parser_test

import (
	"sacco/server/parser"
	"testing"
)

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
