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
	MemberOccupationSingleChildren = []string{
		"memberOccupationVerification",
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
	AccountTransactionArrayChildren = []string{
		"accountJournal",
	}
	MemberBusinessArrayChildren = []string{
		"memberLastYearBusinessHistory",
		"memberNextYearBusinessProjection",
	}
	SingleChildren = map[string][]string{
		"MemberSingleChildren":           MemberSingleChildren,
		"MemberOccupationSingleChildren": MemberOccupationSingleChildren,
		"MemberLoanSingleChildren":       MemberLoanSingleChildren,
	}
	ArrayChildren = map[string][]string{
		"MemberArrayChildren":             MemberArrayChildren,
		"AccountArrayChildren":            AccountArrayChildren,
		"AccountTransactionArrayChildren": AccountTransactionArrayChildren,
		"MemberBusinessArrayChildren":     MemberBusinessArrayChildren,
	}
	FloatFields = []string{
		"grossPay",
		"netPay",
		"periodEmployedInMonths",
		"totalCosts",
		"totalCostOfGoods",
		"employeesWages",
		"transport",
		"loanInterest",
		"rentals",
		"otherCosts",
		"netProfitLoss",
		"financialYear",
		"totalIncome",
		"ownSalary",
		"utilities",
		"password",
		"debit",
		"credit",
		"amountRecommended",
		"amountApproved",
		"numberOfShares",
		"pricePerShare",
		"loanAmount",
		"repaymentPeriodInMonths",
		"value",
		"value",
		"transport",
		"loanInterest",
		"utilities",
		"rentals",
		"otherCosts",
		"financialYear",
		"totalIncome",
		"totalCosts",
		"netProfitLoss",
		"totalCostOfGoods",
		"employeesWages",
		"ownSalary",
	}
	ParentModels = map[string][]string{
		"memberContact": {
			"member",
		},
		"memberOccupation": {
			"memberLoan",
		},
		"memberNextYearBusinessProjection": {
			"memberBusiness",
		},
		"accountJournal": {
			"account",
		},
		"memberNominee": {
			"member",
		},
		"memberLoanWitness": {
			"memberLoan",
		},
		"memberLoanApproval": {
			"memberLoan",
		},
		"memberBeneficiary": {
			"member",
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
		"memberLoanSecurity": {
			"memberLoan",
		},
		"memberOccupationVerification": {
			"memberOccupation",
		},
		"accountTransaction": {
			"account",
		},
		"memberBusiness": {
			"memberLoan",
		},
		"memberLastYearBusinessHistory": {
			"memberBusiness",
		},
	}
)
