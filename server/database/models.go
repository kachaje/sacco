package database

var (
	MemberArrayChildren = []string{
		"memberBeneficiary",
		"memberShares",
		"memberLoan",
	}
	MemberSingleChildren = []string{
		"memberContact",
		"memberNominee",
	}
	MemberBusinessArrayChildren = []string{
		"memberLastYearBusinessHistory",
		"memberNextYearBusinessProjection",
	}
	AccountArrayChildren = []string{
		"accountJournal",
		"accountTransaction",
	}
	MemberOccupationSingleChildren = []string{
		"memberOccupationVerification",
	}
	MemberLoanSingleChildren = []string{
		"memberBusiness",
		"memberOccupation",
		"memberLoanLiability",
		"memberLoanSecurity",
		"memberLoanWitness",
		"memberLoanApproval",
	}
	AccountTransactionArrayChildren = []string{
		"accountJournal",
	}
	SingleChildren = map[string][]string{
		"MemberSingleChildren":           MemberSingleChildren,
		"MemberOccupationSingleChildren": MemberOccupationSingleChildren,
		"MemberLoanSingleChildren":       MemberLoanSingleChildren,
	}
	ArrayChildren = map[string][]string{
		"MemberArrayChildren":             MemberArrayChildren,
		"MemberBusinessArrayChildren":     MemberBusinessArrayChildren,
		"AccountArrayChildren":            AccountArrayChildren,
		"AccountTransactionArrayChildren": AccountTransactionArrayChildren,
	}
	FloatFields = []string{
		"totalIncome",
		"totalCostOfGoods",
		"loanInterest",
		"utilities",
		"totalCosts",
		"financialYear",
		"employeesWages",
		"ownSalary",
		"transport",
		"rentals",
		"otherCosts",
		"netProfitLoss",
		"numberOfShares",
		"pricePerShare",
		"value",
		"amountRecommended",
		"amountApproved",
		"value",
		"debit",
		"credit",
		"grossPay",
		"netPay",
		"periodEmployedInMonths",
		"loanAmount",
		"repaymentPeriodInMonths",
		"password",
		"financialYear",
		"employeesWages",
		"transport",
		"loanInterest",
		"rentals",
		"otherCosts",
		"totalCosts",
		"netProfitLoss",
		"totalIncome",
		"totalCostOfGoods",
		"ownSalary",
		"utilities",
	}
	ParentModels = map[string][]string{
		"memberBeneficiary": {
			"member",
		},
		"memberBusiness": {
			"memberLoan",
		},
		"memberNextYearBusinessProjection": {
			"memberBusiness",
		},
		"memberShares": {
			"member",
		},
		"memberLoanLiability": {
			"memberLoan",
		},
		"memberLoanWitness": {
			"memberLoan",
		},
		"memberLoanApproval": {
			"memberLoan",
		},
		"memberLoanSecurity": {
			"memberLoan",
		},
		"accountJournal": {
			"account",
		},
		"memberNominee": {
			"member",
		},
		"memberOccupation": {
			"memberLoan",
		},
		"memberLoan": {
			"member",
		},
		"accountTransaction": {
			"account",
		},
		"memberContact": {
			"member",
		},
		"memberLastYearBusinessHistory": {
			"memberBusiness",
		},
		"memberOccupationVerification": {
			"memberOccupation",
		},
	}
)
