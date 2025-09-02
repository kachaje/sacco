package database

var (
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
	MemberLoanSingleChildren = []string{
		"memberBusiness",
		"memberOccupation",
		"memberLoanLiability",
		"memberLoanSecurity",
		"memberLoanWitness",
		"memberLoanApproval",
	}
	SingleChildren = map[string][]string{
		"MemberOccupationSingleChildren": MemberOccupationSingleChildren,
		"MemberSingleChildren":           MemberSingleChildren,
		"MemberLoanSingleChildren":       MemberLoanSingleChildren,
	}
	ArrayChildren = map[string][]string{
		"AccountArrayChildren":            AccountArrayChildren,
		"MemberBusinessArrayChildren":     MemberBusinessArrayChildren,
		"AccountTransactionArrayChildren": AccountTransactionArrayChildren,
		"MemberArrayChildren":             MemberArrayChildren,
	}
	FloatFields = []string{
		"grossPay",
		"netPay",
		"periodEmployedInMonths",
		"value",
		"amountRecommended",
		"amountApproved",
		"utilities",
		"rentals",
		"otherCosts",
		"netProfitLoss",
		"financialYear",
		"totalIncome",
		"employeesWages",
		"ownSalary",
		"transport",
		"loanInterest",
		"totalCosts",
		"totalCostOfGoods",
		"numberOfShares",
		"pricePerShare",
		"debit",
		"credit",
		"otherCosts",
		"totalCostOfGoods",
		"employeesWages",
		"totalCosts",
		"netProfitLoss",
		"financialYear",
		"totalIncome",
		"ownSalary",
		"transport",
		"loanInterest",
		"utilities",
		"rentals",
		"loanAmount",
		"repaymentPeriodInMonths",
		"value",
		"password",
	}
	ParentModels = map[string][]string{
		"memberContact": {
			"member",
		},
		"memberNominee": {
			"member",
		},
		"memberOccupation": {
			"memberLoan",
		},
		"memberLoanLiability": {
			"memberLoan",
		},
		"memberLoanApproval": {
			"memberLoan",
		},
		"memberBusiness": {
			"memberLoan",
		},
		"memberNextYearBusinessProjection": {
			"memberBusiness",
		},
		"memberShares": {
			"member",
		},
		"memberOccupationVerification": {
			"memberOccupation",
		},
		"accountTransaction": {
			"account",
		},
		"accountJournal": {
			"account",
		},
		"memberBeneficiary": {
			"member",
		},
		"memberLastYearBusinessHistory": {
			"memberBusiness",
		},
		"memberLoan": {
			"member",
		},
		"memberLoanSecurity": {
			"memberLoan",
		},
		"memberLoanWitness": {
			"memberLoan",
		},
	}
	ModelsMap = map[string]string{
		"memberBusinessId":                   "member.memberLoan.0.memberBusiness.id",
		"memberContactId":                    "member.memberContact.id",
		"memberLastYearBusinessHistoryId":    "member.memberLoan.0.memberBusiness.memberLastYearBusinessHistory.0.id",
		"memberNextYearBusinessProjectionId": "member.memberLoan.0.memberBusiness.memberNextYearBusinessProjection.0.id",
		"memberOccupationId":                 "member.memberLoan.0.memberOccupation.id",
		"memberId":                           "member.id",
		"memberLoanApprovalId":               "member.memberLoan.0.memberLoanApproval.id",
		"memberLoanId":                       "member.memberLoan.0.id",
		"memberNomineeId":                    "member.memberNominee.id",
		"memberOccupationVerificationId":     "member.memberLoan.0.memberOccupation.memberOccupationVerification.id",
		"memberBeneficiaryId":                "member.memberBeneficiary.0.id",
	}
)
