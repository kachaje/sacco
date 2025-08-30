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
	MemberLoanSingleChildren = []string{
		"memberBusiness",
		"memberOccupation",
		"memberLoanLiability",
		"memberLoanSecurity",
		"memberLoanWitness",
		"memberOccupationVerification",
		"memberLoanApproval",
	}
	MemberBusinessArrayChildren = []string{
		"memberLastYearBusinessHistory",
		"memberNextYearBusinessProjection",
	}
	AccountTransactionArrayChildren = []string{
		"accountJournal",
	}
	MemberOccupationSingleChildren = []string{
		"memberOccupationVerification",
	}
	AccountArrayChildren = []string{
		"accountJournal",
		"accountTransaction",
	}
	SingleChildren = map[string][]string{
		"MemberSingleChildren":           MemberSingleChildren,
		"MemberLoanSingleChildren":       MemberLoanSingleChildren,
		"MemberOccupationSingleChildren": MemberOccupationSingleChildren,
	}
	ArrayChildren = map[string][]string{
		"MemberArrayChildren":             MemberArrayChildren,
		"MemberBusinessArrayChildren":     MemberBusinessArrayChildren,
		"AccountTransactionArrayChildren": AccountTransactionArrayChildren,
		"AccountArrayChildren":            AccountArrayChildren,
	}
	FloatFields = []string{
		"repaymentPeriodInMonths",
		"loanAmount",
		"totalCostOfGoods",
		"employeesWages",
		"rentals",
		"otherCosts",
		"totalCosts",
		"ownSalary",
		"transport",
		"loanInterest",
		"utilities",
		"netProfitLoss",
		"financialYear",
		"totalIncome",
		"value",
		"value",
		"amountApproved",
		"amountRecommended",
		"password",
		"debit",
		"credit",
		"totalCosts",
		"netProfitLoss",
		"totalIncome",
		"totalCostOfGoods",
		"employeesWages",
		"transport",
		"otherCosts",
		"financialYear",
		"ownSalary",
		"loanInterest",
		"utilities",
		"rentals",
		"grossPay",
		"netPay",
		"periodEmployedInMonths",
		"numberOfShares",
		"pricePerShare",
	}
	ParentModels = map[string][]string{
		"memberNominee": {
			"member",
		},
		"memberLoan": {
			"member",
		},
		"memberBeneficiary": {
			"member",
		},
		"memberNextYearBusinessProjection": {
			"memberBusiness",
		},
		"memberLoanLiability": {
			"memberLoan",
		},
		"memberLoanSecurity": {
			"memberLoan",
		},
		"memberLoanWitness": {
			"memberLoan",
		},
		"memberLoanApproval": {
			"memberLoan",
		},
		"accountJournal": {
			"account",
		},
		"memberContact": {
			"member",
		},
		"memberBusiness": {
			"memberLoan",
		},
		"memberLastYearBusinessHistory": {
			"memberBusiness",
		},
		"memberOccupationVerification": {
			"memberOccupation",
		},
		"accountTransaction": {
			"account",
		},
		"memberOccupation": {
			"memberLoan",
		},
		"memberShares": {
			"member",
		},
	}
)
