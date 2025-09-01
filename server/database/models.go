package database

var (
	MemberBusinessArrayChildren = []string{
		"memberLastYearBusinessHistory",
		"memberNextYearBusinessProjection",
	}
	AccountTransactionArrayChildren = []string{
		"accountJournal",
	}
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
	SingleChildren = map[string][]string{
		"MemberLoanSingleChildren":       MemberLoanSingleChildren,
		"MemberSingleChildren":           MemberSingleChildren,
		"MemberOccupationSingleChildren": MemberOccupationSingleChildren,
	}
	ArrayChildren = map[string][]string{
		"MemberBusinessArrayChildren":     MemberBusinessArrayChildren,
		"AccountTransactionArrayChildren": AccountTransactionArrayChildren,
		"AccountArrayChildren":            AccountArrayChildren,
		"MemberArrayChildren":             MemberArrayChildren,
	}
	FloatFields = []string{
		"otherCosts",
		"netProfitLoss",
		"totalCostOfGoods",
		"employeesWages",
		"loanInterest",
		"rentals",
		"totalCosts",
		"financialYear",
		"totalIncome",
		"ownSalary",
		"transport",
		"utilities",
		"utilities",
		"otherCosts",
		"transport",
		"loanInterest",
		"rentals",
		"totalCosts",
		"netProfitLoss",
		"financialYear",
		"totalIncome",
		"totalCostOfGoods",
		"employeesWages",
		"ownSalary",
		"amountRecommended",
		"amountApproved",
		"debit",
		"credit",
		"periodEmployedInMonths",
		"grossPay",
		"netPay",
		"value",
		"password",
		"value",
		"numberOfShares",
		"pricePerShare",
		"loanAmount",
		"repaymentPeriodInMonths",
	}
	ParentModels = map[string][]string{
		"memberNominee": {
			"member",
		},
		"memberLastYearBusinessHistory": {
			"memberBusiness",
		},
		"memberNextYearBusinessProjection": {
			"memberBusiness",
		},
		"memberLoanApproval": {
			"memberLoan",
		},
		"accountJournal": {
			"account",
		},
		"memberLoanWitness": {
			"memberLoan",
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
		"memberLoanSecurity": {
			"memberLoan",
		},
		"accountTransaction": {
			"account",
		},
		"memberLoanLiability": {
			"memberLoan",
		},
		"memberOccupationVerification": {
			"memberOccupation",
		},
		"memberContact": {
			"member",
		},
		"memberShares": {
			"member",
		},
		"memberLoan": {
			"member",
		},
	}
	ModelsMap = map[string]string{}
)
