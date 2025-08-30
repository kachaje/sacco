package database

var (
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
	AccountArrayChildren = []string{
		"accountJournal",
		"accountTransaction",
	}
	SingleChildren = map[string][]string{
		"MemberLoanSingleChildren": MemberLoanSingleChildren,
		"MemberSingleChildren":     MemberSingleChildren,
	}
	ArrayChildren = map[string][]string{
		"AccountTransactionArrayChildren": AccountTransactionArrayChildren,
		"MemberArrayChildren":             MemberArrayChildren,
		"MemberBusinessArrayChildren":     MemberBusinessArrayChildren,
		"AccountArrayChildren":            AccountArrayChildren,
	}
	FloatFields = []string{
		"periodEmployedInMonths",
		"grossPay",
		"netPay",
		"totalCostOfGoods",
		"employeesWages",
		"netProfitLoss",
		"totalIncome",
		"ownSalary",
		"transport",
		"loanInterest",
		"utilities",
		"rentals",
		"otherCosts",
		"totalCosts",
		"financialYear",
		"loanAmount",
		"repaymentPeriodInMonths",
		"value",
		"amountRecommended",
		"amountApproved",
		"debit",
		"credit",
		"ownSalary",
		"transport",
		"utilities",
		"otherCosts",
		"totalCosts",
		"netProfitLoss",
		"financialYear",
		"employeesWages",
		"loanInterest",
		"rentals",
		"totalIncome",
		"totalCostOfGoods",
		"numberOfShares",
		"pricePerShare",
		"value",
		"password",
	}
	ParentModels = map[string][]string{
		"accountTransaction": {
			"account",
		},
		"memberOccupation": {
			"member",
			"memberLoanId",
		},
		"memberBeneficiary": {
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
		"memberLoanSecurity": {
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
		"memberContact": {
			"member",
		},
		"memberOccupationVerification": {
			"member",
			"memberLoan",
		},
		"memberBusiness": {
			"member",
			"memberLoan",
		},
		"memberLoanWitness": {
			"member",
			"memberLoan",
		},
		"memberNominee": {
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
	}
)
