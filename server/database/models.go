package database

var (
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
	SingleChildren = map[string][]string{
		"MemberOccupationSingleChildren": MemberOccupationSingleChildren,
		"MemberLoanSingleChildren":       MemberLoanSingleChildren,
		"MemberSingleChildren":           MemberSingleChildren,
	}
	ArrayChildren = map[string][]string{
		"MemberBusinessArrayChildren":     MemberBusinessArrayChildren,
		"AccountArrayChildren":            AccountArrayChildren,
		"AccountTransactionArrayChildren": AccountTransactionArrayChildren,
		"MemberArrayChildren":             MemberArrayChildren,
	}
	FloatFields = []string{
		"totalCostOfGoods",
		"employeesWages",
		"loanInterest",
		"utilities",
		"otherCosts",
		"totalCosts",
		"financialYear",
		"totalIncome",
		"ownSalary",
		"transport",
		"rentals",
		"netProfitLoss",
		"amountRecommended",
		"amountApproved",
		"password",
		"debit",
		"credit",
		"grossPay",
		"netPay",
		"periodEmployedInMonths",
		"financialYear",
		"ownSalary",
		"loanInterest",
		"otherCosts",
		"totalIncome",
		"totalCostOfGoods",
		"employeesWages",
		"transport",
		"utilities",
		"rentals",
		"totalCosts",
		"netProfitLoss",
		"loanAmount",
		"repaymentPeriodInMonths",
		"value",
		"value",
		"numberOfShares",
		"pricePerShare",
	}
	ParentModels = map[string][]string{
		"memberContact": {
			"member",
		},
		"memberNextYearBusinessProjection": {
			"memberBusiness",
		},
		"memberLoanWitness": {
			"memberLoan",
		},
		"memberOccupationVerification": {
			"memberOccupation",
		},
		"memberNominee": {
			"member",
		},
		"memberLoanApproval": {
			"memberLoan",
		},
		"accountJournal": {
			"account",
		},
		"memberOccupation": {
			"memberLoan",
		},
		"memberBeneficiary": {
			"member",
		},
		"memberBusiness": {
			"memberLoan",
		},
		"memberLastYearBusinessHistory": {
			"memberBusiness",
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
		"memberShares": {
			"member",
		},
		"accountTransaction": {
			"account",
		},
	}
)
