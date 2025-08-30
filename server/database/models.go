package database

var (
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
	SingleChildren = map[string][]string{
		"MemberLoanSingleChildren": MemberLoanSingleChildren,
		"MemberSingleChildren":     MemberSingleChildren,
	}
	ArrayChildren = map[string][]string{
		"AccountArrayChildren":            AccountArrayChildren,
		"AccountTransactionArrayChildren": AccountTransactionArrayChildren,
		"MemberArrayChildren":             MemberArrayChildren,
		"MemberBusinessArrayChildren":     MemberBusinessArrayChildren,
	}
	FloatFields = []string{
		"totalCosts",
		"totalIncome",
		"employeesWages",
		"ownSalary",
		"loanInterest",
		"utilities",
		"netProfitLoss",
		"financialYear",
		"totalCostOfGoods",
		"transport",
		"rentals",
		"otherCosts",
		"numberOfShares",
		"pricePerShare",
		"loanAmount",
		"repaymentPeriodInMonths",
		"password",
		"grossPay",
		"netPay",
		"periodEmployedInMonths",
		"financialYear",
		"loanInterest",
		"rentals",
		"netProfitLoss",
		"totalIncome",
		"totalCostOfGoods",
		"employeesWages",
		"ownSalary",
		"transport",
		"utilities",
		"otherCosts",
		"totalCosts",
		"value",
		"amountRecommended",
		"amountApproved",
		"value",
		"debit",
		"credit",
	}
	ParentModels = map[string][]string{
		"memberNominee": {
			"member",
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
		"accountTransaction": {
			"account",
		},
		"memberOccupation": {
			"member",
			"memberLoan",
		},
		"memberLoanWitness": {
			"member",
			"memberLoan",
		},
		"memberContact": {
			"member",
		},
		"memberLastYearBusinessHistory": {
			"member",
			"memberLoan",
			"memberBusiness",
		},
		"memberLoanLiability": {
			"member",
			"memberLoan",
		},
		"memberLoanApproval": {
			"member",
			"memberLoan",
		},
		"memberBeneficiary": {
			"member",
		},
		"memberBusiness": {
			"member",
			"memberLoan",
		},
		"memberLoanSecurity": {
			"member",
			"memberLoan",
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
