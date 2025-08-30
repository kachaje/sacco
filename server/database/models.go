package database

var (
	MemberBusinessArrayChildren = []string{
		"memberLastYearBusinessHistory",
		"memberNextYearBusinessProjection",
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
	MemberLoanSingleChildren = []string{
		"memberBusiness",
		"memberOccupation",
		"memberLoanLiability",
		"memberLoanSecurity",
		"memberLoanWitness",
		"memberOccupationVerification",
		"memberLoanApproval",
	}
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
	SingleChildren = map[string][]string{
		"MemberSingleChildren":           MemberSingleChildren,
		"MemberLoanSingleChildren":       MemberLoanSingleChildren,
		"MemberOccupationSingleChildren": MemberOccupationSingleChildren,
	}
	ArrayChildren = map[string][]string{
		"MemberBusinessArrayChildren":     MemberBusinessArrayChildren,
		"MemberArrayChildren":             MemberArrayChildren,
		"AccountArrayChildren":            AccountArrayChildren,
		"AccountTransactionArrayChildren": AccountTransactionArrayChildren,
	}
	FloatFields = []string{
		"grossPay",
		"netPay",
		"periodEmployedInMonths",
		"password",
		"financialYear",
		"totalIncome",
		"employeesWages",
		"transport",
		"loanInterest",
		"utilities",
		"otherCosts",
		"totalCosts",
		"totalCostOfGoods",
		"ownSalary",
		"rentals",
		"netProfitLoss",
		"numberOfShares",
		"pricePerShare",
		"loanAmount",
		"repaymentPeriodInMonths",
		"value",
		"amountRecommended",
		"amountApproved",
		"employeesWages",
		"transport",
		"utilities",
		"rentals",
		"netProfitLoss",
		"financialYear",
		"totalCostOfGoods",
		"ownSalary",
		"loanInterest",
		"otherCosts",
		"totalCosts",
		"totalIncome",
		"value",
		"debit",
		"credit",
	}
	ParentModels = map[string][]string{
		"memberContact": {
			"member",
		},
		"memberOccupation": {
			"memberLoan",
		},
		"memberBusiness": {
			"memberLoan",
		},
		"memberLoanWitness": {
			"memberLoan",
		},
		"memberOccupationVerification": {
			"memberOccupation",
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
		"memberLoan": {
			"member",
		},
		"memberLoanLiability": {
			"memberLoan",
		},
		"memberLoanApproval": {
			"memberLoan",
		},
		"memberNominee": {
			"member",
		},
		"memberNextYearBusinessProjection": {
			"memberBusiness",
		},
		"memberLoanSecurity": {
			"memberLoan",
		},
		"accountTransaction": {
			"account",
		},
		"accountJournal": {
			"account",
		},
	}
)
