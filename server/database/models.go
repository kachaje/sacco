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
	AccountTransactionArrayChildren = []string{
		"accountJournal",
	}
	MemberBusinessArrayChildren = []string{
		"memberLastYearBusinessHistory",
		"memberNextYearBusinessProjection",
	}
	AccountArrayChildren = []string{
		"accountJournal",
		"accountTransaction",
	}
	SingleChildren = map[string][]string{
		"MemberLoanSingleChildren": MemberLoanSingleChildren,
		"MemberSingleChildren":     MemberSingleChildren,
	}
	ArrayChildren = map[string][]string{
		"MemberArrayChildren":             MemberArrayChildren,
		"AccountTransactionArrayChildren": AccountTransactionArrayChildren,
		"MemberBusinessArrayChildren":     MemberBusinessArrayChildren,
		"AccountArrayChildren":            AccountArrayChildren,
	}
	FloatFields = []string{
		"netPay",
		"periodEmployedInMonths",
		"grossPay",
		"totalCostOfGoods",
		"rentals",
		"netProfitLoss",
		"totalIncome",
		"employeesWages",
		"ownSalary",
		"transport",
		"loanInterest",
		"utilities",
		"otherCosts",
		"totalCosts",
		"financialYear",
		"numberOfShares",
		"pricePerShare",
		"loanAmount",
		"repaymentPeriodInMonths",
		"value",
		"utilities",
		"otherCosts",
		"totalCosts",
		"totalIncome",
		"ownSalary",
		"transport",
		"loanInterest",
		"rentals",
		"netProfitLoss",
		"financialYear",
		"totalCostOfGoods",
		"employeesWages",
		"value",
		"amountRecommended",
		"amountApproved",
		"password",
		"credit",
		"debit",
	}
	ParentModels = map[string][]string{
		"memberBusiness": {
			"member",
			"memberLoan",
		},
		"memberContact": {
			"member",
		},
		"memberOccupation": {
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
		"memberShares": {
			"member",
		},
		"memberLoan": {
			"member",
		},
		"memberLoanSecurity": {
			"member",
			"memberLoan",
		},
		"memberOccupationVerification": {
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
		"memberLoanApproval": {
			"member",
			"memberLoan",
		},
		"accountJournal": {
			"account",
		},
		"memberNominee": {
			"member",
		},
		"memberLoanWitness": {
			"member",
			"memberLoan",
		},
		"accountTransaction": {
			"account",
		},
	}
)
