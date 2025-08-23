CREATE TABLE
  IF NOT EXISTS member (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    firstName TEXT NOT NULL,
    lastName TEXT NOT NULL,
    otherName TEXT,
    gender TEXT CHECK (gender IN ('Male', 'Female')),
    title TEXT,
    maritalStatus TEXT,
    dateOfBirth TEXT,
    nationalId TEXT,
    utilityBillType TEXT,
    utilityBillNumber TEXT,
    fileNumber TEXT,
    oldFileNumber TEXT,
    phoneNumber TEXT NOT NULL,
    memberIdNumber TEXT,
    shortMemberId TEXT,
    dateJoined TEXT,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
  );

CREATE TRIGGER IF NOT EXISTS memberUpdated AFTER
UPDATE ON member FOR EACH ROW BEGIN
UPDATE member
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS memberContact (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberId INTEGER NOT NULL,
    postalAddress TEXT,
    residentialAddress TEXT,
    email TEXT,
    homeVillage TEXT,
    homeTA TEXT,
    homeDistrict TEXT,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
  );

CREATE TRIGGER IF NOT EXISTS memberContactUpdated AFTER
UPDATE ON memberContact FOR EACH ROW BEGIN
UPDATE memberContact
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS memberNominee (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberId INTEGER NOT NULL,
    nomineeName TEXT,
    nomineePhone TEXT,
    nomineeAddress TEXT,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
  );

CREATE TRIGGER IF NOT EXISTS memberNomineeUpdated AFTER
UPDATE ON memberNominee FOR EACH ROW BEGIN
UPDATE memberNominee
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS memberOccupation (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberId INTEGER NOT NULL,
    employerName TEXT,
    employerAddress TEXT,
    employerPhone TEXT,
    jobTitle TEXT,
    periodEmployed REAL,
    grossPay REAL,
    netPay REAL,
    highestQualification TEXT,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
  );

CREATE TRIGGER IF NOT EXISTS memberOccupationUpdated AFTER
UPDATE ON memberOccupation FOR EACH ROW BEGIN
UPDATE memberOccupation
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS memberBeneficiary (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberId INTEGER NOT NULL,
    name TEXT,
    percentage REAL,
    contact TEXT,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
  );

CREATE TRIGGER IF NOT EXISTS memberBeneficiaryUpdated AFTER
UPDATE ON memberBeneficiary FOR EACH ROW BEGIN
UPDATE memberBeneficiary
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS memberBusiness (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberId INTEGER NOT NULL,
    yearsInBusiness REAL,
    businessNature TEXT,
    businessName TEXT,
    tradingArea TEXT,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
  );

CREATE TRIGGER IF NOT EXISTS memberBusinessUpdated AFTER
UPDATE ON memberBusiness FOR EACH ROW BEGIN
UPDATE memberBusiness
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS memberLastYearBusinessHistory (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberBusinessId INTEGER NOT NULL,
    memberLoanId INTEGER NOT NULL,
    financialYear TEXT,
    totalIncome REAL,
    totalCostOfGoods REAL,
    employeesWages REAL,
    ownSalary REAL,
    transport REAL,
    loanInterest REAL,
    utilities REAL,
    rentals REAL,
    otherCosts REAL,
    totalCosts REAL,
    netProfitLoss REAL,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
  );

CREATE TRIGGER IF NOT EXISTS memberLastYearBusinessHistoryUpdated AFTER
UPDATE ON memberLastYearBusinessHistory FOR EACH ROW BEGIN
UPDATE memberLastYearBusinessHistory
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS memberNextYearBusinessProjection (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberBusinessId INTEGER NOT NULL,
    memberLoanId INTEGER NOT NULL,
    financialYear TEXT,
    totalIncome REAL,
    totalCostOfGoods REAL,
    employeesWages REAL,
    ownSalary REAL,
    transport REAL,
    loanInterest REAL,
    utilities REAL,
    rentals REAL,
    otherCosts REAL,
    totalCosts REAL,
    netProfitLoss REAL,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
  );

CREATE TRIGGER IF NOT EXISTS memberNextYearBusinessProjectionUpdated AFTER
UPDATE ON memberNextYearBusinessProjection FOR EACH ROW BEGIN
UPDATE memberNextYearBusinessProjection
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS memberShares (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberId INTEGER NOT NULL,
    numberOfShares REAL,
    pricePerShare REAL,
    sharesType TEXT NOT NULL CHECK (sharesType IN ('FIXED', 'REDEEMABLE')),
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
  );

CREATE TRIGGER IF NOT EXISTS memberSharesUpdated AFTER
UPDATE ON memberShares FOR EACH ROW BEGIN
UPDATE memberShares
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS memberLoan (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberId INTEGER NOT NULL,
    employeeNumber TEXT,
    loanAmount REAL,
    repaymentPeriod REAL,
    loanPurpose TEXT,
    loanStatus TEXT NOT NULL CHECK (loanStatus IN ('PENDING', 'APPROVED', 'REJECTED')),
    loanType TEXT NOT NULL CHECK (
      loanType IN (
        'PERSONAL',
        'BUSINESS',
        'AGRICULTURAL',
        'EMERGENCY'
      )
    ),
    amountRecommended REAL,
    approvedBy TEXT,
    approvalDate TEXT,
    amountApproved REAL,
    verifiedBy TEXT,
    dateVerified TEXT,
    denialOrPartialReason TEXT,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
  );

CREATE TRIGGER IF NOT EXISTS memberLoanUpdated AFTER
UPDATE ON memberLoan FOR EACH ROW BEGIN
UPDATE memberLoan
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS memberLoanLiability (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberId INTEGER NOT NULL,
    memberLoanId INTEGER NOT NULL,
    description TEXT,
    value REAL,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
  );

CREATE TRIGGER IF NOT EXISTS memberLoanLiabilityUpdated AFTER
UPDATE ON memberLoanLiability FOR EACH ROW BEGIN
UPDATE memberLoanLiability
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS memberLoanSecurity (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberId INTEGER NOT NULL,
    memberLoanId INTEGER NOT NULL,
    description TEXT,
    value REAL,
    serialNumber TEXT,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
  );

CREATE TRIGGER IF NOT EXISTS memberLoanSecurityUpdated AFTER
UPDATE ON memberLoanSecurity FOR EACH ROW BEGIN
UPDATE memberLoanSecurity
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS memberLoanWitness (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberId INTEGER NOT NULL,
    memberLoanId INTEGER NOT NULL,
    name TEXT,
    telephone TEXT,
    address TEXT,
    date TEXT,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
  );

CREATE TRIGGER IF NOT EXISTS memberLoanWitnessUpdated AFTER
UPDATE ON memberLoanWitness FOR EACH ROW BEGIN
UPDATE memberLoanWitness
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS employmentVerification (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberId INTEGER NOT NULL,
    memberLoanId INTEGER NOT NULL,
    jobVerified TEXT NOT NULL CHECK (jobVerified IN ('Yes', 'No')) DEFAULT 'No',
    periodEmployed REAL,
    jobTitle TEXT,
    employerName TEXT,
    employerAddress TEXT,
    employerPhone TEXT,
    grossPay REAL,
    grossVerified TEXT NOT NULL CHECK (grossVerified IN ('Yes', 'No')) DEFAULT 'No',
    netPay REAL,
    netVerified TEXT NOT NULL CHECK (netVerified IN ('Yes', 'No')) DEFAULT 'No',
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
  );

CREATE TRIGGER IF NOT EXISTS employmentVerificationUpdated AFTER
UPDATE ON employmentVerification FOR EACH ROW BEGIN
UPDATE employmentVerification
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;