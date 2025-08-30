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
		"grossPay",
		"netPay",
		"periodEmployedInMonths",
		"pricePerShare",
		"numberOfShares",
		"repaymentPeriodInMonths",
		"loanAmount",
		"financialYear",
		"totalIncome",
		"employeesWages",
		"ownSalary",
		"transport",
		"utilities",
		"totalCosts",
		"totalCostOfGoods",
		"loanInterest",
		"rentals",
		"otherCosts",
		"netProfitLoss",
		"value",
		"otherCosts",
		"totalCosts",
		"netProfitLoss",
		"financialYear",
		"employeesWages",
		"ownSalary",
		"transport",
		"utilities",
		"rentals",
		"totalIncome",
		"totalCostOfGoods",
		"loanInterest",
		"debit",
		"credit",
		"value",
		"amountApproved",
		"amountRecommended",
		"password",
	}
	ParentModels = map[string][]string{
		"memberOccupation": {
			"member",
			"memberLoanId",
		},
		"memberBusiness": {
			"member",
			"memberLoan",
		},
		"memberShares": {
			"member",
		},
		"memberLoan": {
			"member",
		},
		"memberOccupationVerification": {
			"member",
			"memberLoan",
		},
		"accountTransaction": {
			"account",
		},
		"memberNextYearBusinessProjection": {
			"member",
			"memberLoan",
			"memberBusiness",
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
		"memberLoanWitness": {
			"member",
			"memberLoan",
		},
		"accountJournal": {
			"account",
		},
		"memberContact": {
			"member",
		},
		"memberNominee": {
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
	}
)
