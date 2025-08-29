package database

var (
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
		"AccountArrayChildren":            AccountArrayChildren,
		"AccountTransactionArrayChildren": AccountTransactionArrayChildren,
		"MemberArrayChildren":             MemberArrayChildren,
		"MemberBusinessArrayChildren":     MemberBusinessArrayChildren,
	}
	FloatFields = []string{
		"amountRecommended",
		"amountApproved",
		"debit",
		"credit",
		"loanAmount",
		"repaymentPeriodInMonths",
		"value",
		"password",
		"numberOfShares",
		"pricePerShare",
		"value",
		"grossPay",
		"netPay",
		"periodEmployedInMonths",
		"employeesWages",
		"loanInterest",
		"utilities",
		"rentals",
		"netProfitLoss",
		"totalIncome",
		"totalCostOfGoods",
		"ownSalary",
		"transport",
		"otherCosts",
		"totalCosts",
		"financialYear",
		"totalIncome",
		"totalCostOfGoods",
		"employeesWages",
		"ownSalary",
		"utilities",
		"rentals",
		"otherCosts",
		"netProfitLoss",
		"financialYear",
		"transport",
		"loanInterest",
		"totalCosts",
	}
)
