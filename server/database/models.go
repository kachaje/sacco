package database

var (
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
	MemberOccupationSingleChildren = []string{
		"memberOccupationVerification",
	}
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
		"memberLoan",
	}
	MemberSingleChildren = []string{
		"memberContact",
		"memberNominee",
	}
	SingleChildren = map[string][]string{
		"MemberLoanSingleChildren":       MemberLoanSingleChildren,
		"MemberOccupationSingleChildren": MemberOccupationSingleChildren,
		"MemberSingleChildren":           MemberSingleChildren,
	}
	ArrayChildren = map[string][]string{
		"AccountArrayChildren":            AccountArrayChildren,
		"MemberBusinessArrayChildren":     MemberBusinessArrayChildren,
		"AccountTransactionArrayChildren": AccountTransactionArrayChildren,
		"MemberArrayChildren":             MemberArrayChildren,
	}
	FloatFields = []string{
		"value",
		"amountRecommended",
		"amountApproved",
		"grossPay",
		"netPay",
		"periodEmployedInMonths",
		"value",
		"loanAmount",
		"repaymentPeriodInMonths",
		"totalIncome",
		"totalCostOfGoods",
		"employeesWages",
		"utilities",
		"totalCosts",
		"netProfitLoss",
		"ownSalary",
		"transport",
		"loanInterest",
		"rentals",
		"otherCosts",
		"financialYear",
		"password",
		"debit",
		"credit",
		"totalIncome",
		"totalCostOfGoods",
		"employeesWages",
		"ownSalary",
		"loanInterest",
		"utilities",
		"totalCosts",
		"financialYear",
		"transport",
		"rentals",
		"otherCosts",
		"netProfitLoss",
		"numberOfShares",
		"pricePerShare",
	}
	ParentModels = map[string][]string{
		"memberLoanLiability": {
			"memberLoan",
		},
		"memberLoanApproval": {
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
		"memberLoanWitness": {
			"memberLoan",
		},
		"accountTransaction": {
			"account",
		},
		"memberContact": {
			"member",
		},
		"memberLoan": {
			"member",
		},
		"memberLastYearBusinessHistory": {
			"memberBusiness",
		},
		"memberOccupationVerification": {
			"memberOccupation",
		},
		"accountJournal": {
			"account",
		},
		"memberNominee": {
			"member",
		},
		"memberNextYearBusinessProjection": {
			"memberBusiness",
		},
		"memberShares": {
			"member",
		},
	}
	ModelsMap = map[string]string{
		"member":                           "member.id",
		"memberBeneficiary":                "member.memberBeneficiary.0.id",
		"memberBusiness":                   "member.memberLoan.0.memberBusiness.id",
		"memberContact":                    "member.memberContact.id",
		"memberNextYearBusinessProjection": "member.memberLoan.0.memberBusiness.memberNextYearBusinessProjection.0.id",
		"memberNominee":                    "member.memberNominee.id",
		"memberOccupationVerification":     "member.memberLoan.0.memberOccupation.memberOccupationVerification.id",
		"memberLastYearBusinessHistory":    "member.memberLoan.0.memberBusiness.memberLastYearBusinessHistory.0.id",
		"memberLoan":                       "member.memberLoan.0.id",
		"memberLoanApproval":               "member.memberLoan.0.memberLoanApproval.id",
		"memberOccupation":                 "member.memberLoan.0.memberOccupation.id",
	}
)
