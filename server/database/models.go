package database

var (
	MemberBusinessArrayChildren = []string{
		"memberLastYearBusinessHistory",
		"memberNextYearBusinessProjection",
	}
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
	SingleChildren = map[string][]string{
		"MemberSingleChildren":     MemberSingleChildren,
		"MemberLoanSingleChildren": MemberLoanSingleChildren,
	}
	ArrayChildren = map[string][]string{
		"MemberBusinessArrayChildren":     MemberBusinessArrayChildren,
		"AccountTransactionArrayChildren": AccountTransactionArrayChildren,
		"MemberArrayChildren":             MemberArrayChildren,
		"AccountArrayChildren":            AccountArrayChildren,
	}
	FloatFields = []string{
		"rentals",
		"otherCosts",
		"financialYear",
		"employeesWages",
		"totalCosts",
		"netProfitLoss",
		"totalIncome",
		"totalCostOfGoods",
		"ownSalary",
		"transport",
		"loanInterest",
		"utilities",
		"value",
		"debit",
		"credit",
		"numberOfShares",
		"pricePerShare",
		"value",
		"amountRecommended",
		"amountApproved",
		"password",
		"grossPay",
		"netPay",
		"periodEmployedInMonths",
		"otherCosts",
		"totalIncome",
		"ownSalary",
		"totalCosts",
		"netProfitLoss",
		"financialYear",
		"totalCostOfGoods",
		"employeesWages",
		"transport",
		"loanInterest",
		"utilities",
		"rentals",
		"loanAmount",
		"repaymentPeriodInMonths",
	}
	ParentModels = map[string][]string{
		"memberContact": {
			"member",
		},
		"memberBeneficiary": {
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
		"memberLoanLiability": {
			"member",
			"memberLoan",
		},
		"accountTransaction": {
			"account",
		},
		"accountJournal": {
			"account",
		},
		"memberShares": {
			"member",
		},
		"memberLoanSecurity": {
			"member",
			"memberLoan",
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
			"memberLoan",
		},
		"memberLastYearBusinessHistory": {
			"member",
			"memberLoan",
			"memberBusiness",
		},
		"memberLoan": {
			"member",
		},
		"memberLoanWitness": {
			"member",
			"memberLoan",
		},
		"memberOccupationVerification": {
			"member",
			"memberLoan",
		},
	}
)
