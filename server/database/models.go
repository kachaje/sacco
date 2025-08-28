package database

var (
	AccountArrayChildren = []string{
		"accountJournal",
		"accountTransaction",
	}
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

	SingleChildren = map[string][]string{
		"MemberLoanSingleChildren": MemberLoanSingleChildren,
		"MemberSingleChildren":     MemberSingleChildren,
	}

	ArrayChildren = map[string][]string{
		"MemberBusinessArrayChildren":     MemberBusinessArrayChildren,
		"MemberArrayChildren":             MemberArrayChildren,
		"AccountTransactionArrayChildren": AccountTransactionArrayChildren,
		"AccountArrayChildren":            AccountArrayChildren,
	}
)
