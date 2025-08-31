package database

var (
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
		"netProfitLoss",
		"financialYear",
		"totalIncome",
		"totalCostOfGoods",
		"loanInterest",
		"rentals",
		"otherCosts",
		"employeesWages",
		"ownSalary",
		"transport",
		"utilities",
		"totalCosts",
		"repaymentPeriodInMonths",
		"loanAmount",
		"password",
		"debit",
		"credit",
		"numberOfShares",
		"pricePerShare",
		"value",
		"totalCostOfGoods",
		"transport",
		"loanInterest",
		"utilities",
		"otherCosts",
		"totalIncome",
		"employeesWages",
		"ownSalary",
		"rentals",
		"totalCosts",
		"netProfitLoss",
		"financialYear",
		"value",
		"grossPay",
		"netPay",
		"periodEmployedInMonths",
		"amountRecommended",
		"amountApproved",
	}
	ParentModels = map[string][]string{
		"memberBeneficiary": {
			"member",
		},
		"memberNextYearBusinessProjection": {
			"memberBusiness",
		},
		"memberLoan": {
			"member",
		},
		"memberOccupationVerification": {
			"memberOccupation",
		},
		"accountJournal": {
			"account",
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
		"accountTransaction": {
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
		"memberLoanSecurity": {
			"memberLoan",
		},
		"memberNominee": {
			"member",
		},
		"memberOccupation": {
			"memberLoan",
		},
		"memberLoanApproval": {
			"memberLoan",
		},
	}
)
