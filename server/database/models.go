package database

var (
	AccountTransactionArrayChildren = []string{
		"accountJournal",
	}
	MemberArrayChildren = []string{
		"memberBeneficiary",
		"memberShares",
	}
	MemberSingleChildren = []string{
		"memberContact",
		"memberNominee",
		"memberOccupation",
		"memberLoan",
		"memberBusiness",
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
	SingleChildren = map[string][]string{
		"MemberSingleChildren":     MemberSingleChildren,
		"MemberLoanSingleChildren": MemberLoanSingleChildren,
	}
	ArrayChildren = map[string][]string{
		"AccountTransactionArrayChildren": AccountTransactionArrayChildren,
		"MemberArrayChildren":             MemberArrayChildren,
		"MemberBusinessArrayChildren":     MemberBusinessArrayChildren,
		"AccountArrayChildren":            AccountArrayChildren,
	}
	FloatFields = []string{
		"financialYear",
		"totalCostOfGoods",
		"employeesWages",
		"transport",
		"totalCosts",
		"netProfitLoss",
		"totalIncome",
		"ownSalary",
		"loanInterest",
		"utilities",
		"rentals",
		"otherCosts",
		"numberOfShares",
		"pricePerShare",
		"loanAmount",
		"repaymentPeriodInMonths",
		"value",
		"value",
		"totalIncome",
		"totalCostOfGoods",
		"ownSalary",
		"transport",
		"utilities",
		"rentals",
		"otherCosts",
		"totalCosts",
		"financialYear",
		"employeesWages",
		"loanInterest",
		"netProfitLoss",
		"amountRecommended",
		"amountApproved",
		"password",
		"grossPay",
		"netPay",
		"periodEmployedInMonths",
		"debit",
		"credit",
	}
	ParentModels = map[string][]string{
		"memberLoanWitness": {
			"member",
			"memberLoan",
		},
		"accountTransaction": {
			"account",
		},
		"memberContact": {
			"member",
		},
		"memberBusiness": {
			"member",
			"memberLoan",
		},
		"memberNextYearBusinessProjection": {
			"member",
			"memberLoan",
			"memberBusiness",
		},
		"memberShares": {
			"member",
		},
		"memberLoan": {
			"member",
		},
		"memberLoanLiability": {
			"member",
			"memberLoan",
		},
		"memberLoanSecurity": {
			"member",
			"memberLoan",
		},
		"memberBeneficiary": {
			"member",
		},
		"memberLastYearBusinessHistory": {
			"member",
			"memberLoan",
			"memberBusiness",
		},
		"memberLoanApproval": {
			"member",
			"memberLoan",
		},
		"memberNominee": {
			"member",
		},
		"memberOccupation": {
			"member",
		},
		"memberOccupationVerification": {
			"member",
			"memberLoan",
		},
		"accountJournal": {
			"account",
		},
	}
)
