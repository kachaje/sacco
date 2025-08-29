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
	SingleChildren = map[string][]string{
		"MemberLoanSingleChildren": MemberLoanSingleChildren,
		"MemberSingleChildren":     MemberSingleChildren,
	}
	ArrayChildren = map[string][]string{
		"AccountArrayChildren":            AccountArrayChildren,
		"MemberArrayChildren":             MemberArrayChildren,
		"MemberBusinessArrayChildren":     MemberBusinessArrayChildren,
		"AccountTransactionArrayChildren": AccountTransactionArrayChildren,
	}
	FloatFields = []string{
		"transport",
		"rentals",
		"otherCosts",
		"financialYear",
		"totalCostOfGoods",
		"employeesWages",
		"ownSalary",
		"loanInterest",
		"utilities",
		"totalCosts",
		"netProfitLoss",
		"totalIncome",
		"loanAmount",
		"repaymentPeriodInMonths",
		"financialYear",
		"totalCostOfGoods",
		"ownSalary",
		"transport",
		"loanInterest",
		"rentals",
		"otherCosts",
		"totalCosts",
		"totalIncome",
		"employeesWages",
		"utilities",
		"netProfitLoss",
		"numberOfShares",
		"pricePerShare",
		"value",
		"amountRecommended",
		"amountApproved",
		"password",
		"debit",
		"credit",
		"grossPay",
		"netPay",
		"periodEmployedInMonths",
		"value",
	}
	ParentModels = map[string][]string{
		"memberContact": {
			"member",
		},
		"memberNextYearBusinessProjection": {
			"member",
			"memberLoan",
			"memberBusiness",
		},
		"memberLoan": {
			"member",
		},
		"memberOccupationVerification": {
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
		"memberLoanLiability": {
			"member",
			"memberLoan",
		},
		"memberLoanApproval": {
			"member",
			"memberLoan",
		},
		"account": {},
		"memberNominee": {
			"member",
		},
		"memberLoanWitness": {
			"member",
			"memberLoan",
		},
		"user": {},
		"accountJournal": {
			"account",
		},
		"member": {},
		"memberOccupation": {
			"member",
			"memberLoan",
		},
		"memberBusiness": {
			"member",
			"memberLoan",
		},
		"memberLoanSecurity": {
			"member",
			"memberLoan",
		},
		"accountTransaction": {
			"account",
		},
	}
)
