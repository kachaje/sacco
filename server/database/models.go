package database

var (
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
	AccountTransactionArrayChildren = []string{
		"accountJournal",
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
		"MemberArrayChildren":             MemberArrayChildren,
		"MemberBusinessArrayChildren":     MemberBusinessArrayChildren,
		"AccountTransactionArrayChildren": AccountTransactionArrayChildren,
		"AccountArrayChildren":            AccountArrayChildren,
	}
	FloatFields = []string{
		"financialYear",
		"totalCostOfGoods",
		"employeesWages",
		"ownSalary",
		"transport",
		"loanInterest",
		"rentals",
		"otherCosts",
		"totalIncome",
		"utilities",
		"totalCosts",
		"netProfitLoss",
		"pricePerShare",
		"numberOfShares",
		"value",
		"password",
		"grossPay",
		"netPay",
		"periodEmployedInMonths",
		"loanAmount",
		"repaymentPeriodInMonths",
		"amountRecommended",
		"amountApproved",
		"value",
		"debit",
		"credit",
		"financialYear",
		"totalIncome",
		"ownSalary",
		"utilities",
		"rentals",
		"totalCosts",
		"netProfitLoss",
		"totalCostOfGoods",
		"employeesWages",
		"transport",
		"loanInterest",
		"otherCosts",
	}
	ParentModels = map[string][]string{
		"memberOccupationVerification": {
			"member",
			"memberLoan",
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
		"memberLoanSecurity": {
			"member",
			"memberLoan",
		},
		"accountTransaction": {
			"account",
		},
		"memberNominee": {
			"member",
		},
		"memberOccupation": {
			"member",
			"memberLoan",
		},
		"memberLoan": {
			"member",
		},
		"memberLoanWitness": {
			"member",
			"memberLoan",
		},
		"memberLoanApproval": {
			"member",
			"memberLoan",
		},
		"memberContact": {
			"member",
		},
		"memberBeneficiary": {
			"member",
		},
		"memberLoanLiability": {
			"member",
			"memberLoan",
		},
		"accountJournal": {
			"account",
		},
		"memberLastYearBusinessHistory": {
			"member",
			"memberLoan",
			"memberBusiness",
		},
	}
)
