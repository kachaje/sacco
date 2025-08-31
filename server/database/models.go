package database

var (
	AccountArrayChildren = []string{
		"accountJournal",
		"accountTransaction",
	}
	AccountTransactionArrayChildren = []string{
		"accountJournal",
	}
	MemberOccupationSingleChildren = []string{
		"memberOccupationVerification",
	}
	MemberBusinessArrayChildren = []string{
		"memberLastYearBusinessHistory",
		"memberNextYearBusinessProjection",
	}
	MemberLoanSingleChildren = []string{
		"memberBusiness",
		"memberOccupation",
		"memberLoanLiability",
		"memberLoanSecurity",
		"memberLoanWitness",
		"memberOccupationVerification",
		"memberLoanApproval",
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
	SingleChildren = map[string][]string{
		"MemberOccupationSingleChildren": MemberOccupationSingleChildren,
		"MemberLoanSingleChildren":       MemberLoanSingleChildren,
		"MemberSingleChildren":           MemberSingleChildren,
	}
	ArrayChildren = map[string][]string{
		"AccountArrayChildren":            AccountArrayChildren,
		"AccountTransactionArrayChildren": AccountTransactionArrayChildren,
		"MemberBusinessArrayChildren":     MemberBusinessArrayChildren,
		"MemberArrayChildren":             MemberArrayChildren,
	}
	FloatFields = []string{
		"financialYear",
		"totalIncome",
		"totalCostOfGoods",
		"employeesWages",
		"ownSalary",
		"utilities",
		"rentals",
		"otherCosts",
		"transport",
		"loanInterest",
		"totalCosts",
		"netProfitLoss",
		"value",
		"password",
		"grossPay",
		"netPay",
		"periodEmployedInMonths",
		"totalIncome",
		"totalCostOfGoods",
		"transport",
		"loanInterest",
		"otherCosts",
		"netProfitLoss",
		"financialYear",
		"employeesWages",
		"ownSalary",
		"utilities",
		"rentals",
		"totalCosts",
		"numberOfShares",
		"pricePerShare",
		"amountRecommended",
		"amountApproved",
		"loanAmount",
		"repaymentPeriodInMonths",
		"value",
		"debit",
		"credit",
	}
	ParentModels = map[string][]string{
		"memberNextYearBusinessProjection": {
			"memberBusiness",
		},
		"memberLoanLiability": {
			"memberLoan",
		},
		"accountTransaction": {
			"account",
		},
		"memberNominee": {
			"member",
		},
		"memberOccupation": {
			"memberLoan",
		},
		"memberBeneficiary": {
			"member",
		},
		"memberLastYearBusinessHistory": {
			"memberBusiness",
		},
		"memberShares": {
			"member",
		},
		"memberLoanApproval": {
			"memberLoan",
		},
		"memberContact": {
			"member",
		},
		"memberBusiness": {
			"memberLoan",
		},
		"memberLoan": {
			"member",
		},
		"memberLoanSecurity": {
			"memberLoan",
		},
		"accountJournal": {
			"account",
		},
		"memberLoanWitness": {
			"memberLoan",
		},
		"memberOccupationVerification": {
			"memberOccupation",
		},
	}
)
