package database

var (
	MemberLoanSingleChildren = []string{
		"memberBusiness",
		"memberOccupation",
		"memberLoanLiability",
		"memberLoanSecurity",
		"memberLoanWitness",
		"memberLoanApproval",
	}
	AccountArrayChildren = []string{
		"accountJournal",
		"accountTransaction",
	}
	AccountTransactionArrayChildren = []string{
		"accountJournal",
	}
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
	MemberOccupationSingleChildren = []string{
		"memberOccupationVerification",
	}
	SingleChildren = map[string][]string{
		"MemberLoanSingleChildren":       MemberLoanSingleChildren,
		"MemberSingleChildren":           MemberSingleChildren,
		"MemberOccupationSingleChildren": MemberOccupationSingleChildren,
	}
	ArrayChildren = map[string][]string{
		"AccountArrayChildren":            AccountArrayChildren,
		"AccountTransactionArrayChildren": AccountTransactionArrayChildren,
		"MemberArrayChildren":             MemberArrayChildren,
		"MemberBusinessArrayChildren":     MemberBusinessArrayChildren,
	}
	FloatFields = []string{
		"password",
		"debit",
		"credit",
		"grossPay",
		"netPay",
		"periodEmployedInMonths",
		"totalIncome",
		"employeesWages",
		"ownSalary",
		"transport",
		"utilities",
		"rentals",
		"totalCosts",
		"netProfitLoss",
		"financialYear",
		"totalCostOfGoods",
		"loanInterest",
		"otherCosts",
		"loanAmount",
		"repaymentPeriodInMonths",
		"numberOfShares",
		"pricePerShare",
		"value",
		"amountRecommended",
		"amountApproved",
		"transport",
		"loanInterest",
		"financialYear",
		"totalCostOfGoods",
		"employeesWages",
		"utilities",
		"rentals",
		"otherCosts",
		"totalCosts",
		"netProfitLoss",
		"totalIncome",
		"ownSalary",
		"value",
	}
	ParentModels = map[string][]string{
		"memberBeneficiary": {
			"member",
		},
		"memberBusiness": {
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
		"memberNextYearBusinessProjection": {
			"memberBusiness",
		},
		"memberLoan": {
			"member",
		},
		"memberLoanWitness": {
			"memberLoan",
		},
		"accountTransaction": {
			"account",
		},
		"memberContact": {
			"member",
		},
		"memberShares": {
			"member",
		},
		"memberLoanLiability": {
			"memberLoan",
		},
		"memberLoanApproval": {
			"memberLoan",
		},
		"memberLastYearBusinessHistory": {
			"memberBusiness",
		},
		"memberLoanSecurity": {
			"memberLoan",
		},
		"memberOccupationVerification": {
			"memberOccupation",
		},
	}
)
