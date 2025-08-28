package database

var (
	MemberBusinessArrayChildren = []string{
		"memberLastYearBusinessHistory",
		"memberNextYearBusinessProjection",
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
		"AccountArrayChildren":            AccountArrayChildren,
		"AccountTransactionArrayChildren": AccountTransactionArrayChildren,
		"MemberArrayChildren":             MemberArrayChildren,
	}
)
