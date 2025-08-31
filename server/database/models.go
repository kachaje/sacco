package database

var (
	AccountTransactionArrayChildren = []string{
		"accountJournal",
	}
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
	AccountArrayChildren = []string{
		"accountJournal",
		"accountTransaction",
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
	MemberOccupationSingleChildren = []string{
		"memberOccupationVerification",
	}
	SingleChildren = map[string][]string{
		"MemberSingleChildren":           MemberSingleChildren,
		"MemberLoanSingleChildren":       MemberLoanSingleChildren,
		"MemberOccupationSingleChildren": MemberOccupationSingleChildren,
	}
	ArrayChildren = map[string][]string{
		"AccountTransactionArrayChildren": AccountTransactionArrayChildren,
		"MemberBusinessArrayChildren":     MemberBusinessArrayChildren,
		"MemberArrayChildren":             MemberArrayChildren,
		"AccountArrayChildren":            AccountArrayChildren,
	}
	FloatFields = []string{
		"value",
		"value",
		"amountRecommended",
		"amountApproved",
		"debit",
		"credit",
		"password",
		"financialYear",
		"totalIncome",
		"totalCostOfGoods",
		"employeesWages",
		"ownSalary",
		"loanInterest",
		"rentals",
		"otherCosts",
		"transport",
		"utilities",
		"totalCosts",
		"netProfitLoss",
		"otherCosts",
		"totalCostOfGoods",
		"employeesWages",
		"ownSalary",
		"rentals",
		"totalCosts",
		"netProfitLoss",
		"financialYear",
		"totalIncome",
		"transport",
		"loanInterest",
		"utilities",
		"numberOfShares",
		"pricePerShare",
		"loanAmount",
		"repaymentPeriodInMonths",
		"periodEmployedInMonths",
		"grossPay",
		"netPay",
	}
	ParentModels = map[string][]string{
		"memberLoanSecurity": {
			"memberLoan",
		},
		"accountTransaction": {
			"account",
		},
		"memberContact": {
			"member",
		},
		"memberBeneficiary": {
			"member",
		},
		"memberBusiness": {
			"memberLoan",
		},
		"memberLoanLiability": {
			"memberLoan",
		},
		"memberLoanApproval": {
			"memberLoan",
		},
		"accountJournal": {
			"account",
		},
		"memberOccupationVerification": {
			"memberOccupation",
		},
		"memberLoanWitness": {
			"memberLoan",
		},
		"memberLastYearBusinessHistory": {
			"memberBusiness",
		},
		"memberNextYearBusinessProjection": {
			"memberBusiness",
		},
		"memberShares": {
			"member",
		},
		"memberLoan": {
			"member",
		},
		"memberNominee": {
			"member",
		},
		"memberOccupation": {
			"memberLoan",
		},
	}
)
