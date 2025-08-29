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
	SingleChildren = map[string][]string{
		"MemberLoanSingleChildren": MemberLoanSingleChildren,
		"MemberSingleChildren":     MemberSingleChildren,
	}
	ArrayChildren = map[string][]string{
		"AccountTransactionArrayChildren": AccountTransactionArrayChildren,
		"MemberBusinessArrayChildren":     MemberBusinessArrayChildren,
		"AccountArrayChildren":            AccountArrayChildren,
		"MemberArrayChildren":             MemberArrayChildren,
	}
	
	FloatFields = []string{
		"netPay", "grossPay", "periodEmployedInMonths", "yearsInBusiness",
		"totalIncome", "totalCostOfGoods", "employeesWages", "ownSalary",
		"transport", "loanInterest", "utilities", "rentals", "otherCosts",
		"totalCosts", "netProfitLoss", "numberOfShares", "pricePerShare",
		"loanAmount", "repaymentPeriodInMonths", "amountRecommended",
		"amountApproved", "value",
	}
)
