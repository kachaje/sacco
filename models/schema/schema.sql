CREATE TABLE
  IF NOT EXISTS account (
    name TEXT UNIQUE NOT NULL,
    number INTEGER UNIQUE NOT NULL,
    normal INTEGER NOT NULL,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
  );

CREATE TRIGGER IF NOT EXISTS accountUpdated AFTER
UPDATE ON account FOR EACH ROW BEGIN
UPDATE account
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS accountTransaction (
    id INTEGER,
    date TEXT NOT NULL,
    description TEXT,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
  );

CREATE TRIGGER IF NOT EXISTS accountTransactionUpdated AFTER
UPDATE ON accountTransaction FOR EACH ROW BEGIN
UPDATE accountTransaction
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS accountJournal (
    account INTEGER,
    accountTransactionId INTEGER,
    date TEXT NOT NULL,
    amount REAL,
    direction INTEGER,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (account) REFERENCES account (number) ON DELETE CASCADE,
    FOREIGN KEY (accountTransactionId) REFERENCES accountTransaction (id) ON DELETE CASCADE
  );

CREATE TRIGGER IF NOT EXISTS accountJournalUpdated AFTER
UPDATE ON accountJournal FOR EACH ROW BEGIN
UPDATE accountJournal
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    name TEXT NOT NULL,
    role TEXT NOT NULL,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
  );

CREATE TRIGGER IF NOT EXISTS userUpdated AFTER
UPDATE ON user FOR EACH ROW BEGIN
UPDATE user
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS userRole (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
  );

CREATE TRIGGER IF NOT EXISTS userRoleUpdated AFTER
UPDATE ON userRole FOR EACH ROW BEGIN
UPDATE userRole
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS memberIdsCache (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    idNumber TEXT NOT NULL UNIQUE,
    claimed INTEGER DEFAULT 0,
    memberId INTEGER,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
  );

CREATE TRIGGER IF NOT EXISTS memberIdsCacheUpdated AFTER
UPDATE ON memberIdsCache FOR EACH ROW BEGIN
UPDATE memberIdsCache
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS memberSavingIdsCache (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    idNumber TEXT NOT NULL UNIQUE,
    claimed INTEGER DEFAULT 0,
    memberSavingId INTEGER,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
  );

CREATE TRIGGER IF NOT EXISTS memberSavingIdNumberUpdated AFTER
UPDATE ON memberSavingIdsCache FOR EACH ROW BEGIN
UPDATE memberSavingIdsCache
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS member (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    firstName TEXT NOT NULL,
    lastName TEXT NOT NULL,
    otherName TEXT,
    gender TEXT CHECK (gender IN ('Male', 'Female')),
    title TEXT,
    maritalStatus TEXT CHECK (
      maritalStatus IN (
        'Married',
        'Single',
        'Widowed',
        'Divorced',
        'Other'
      )
    ),
    dateOfBirth TEXT,
    nationalId TEXT,
    utilityBillType TEXT,
    utilityBillNumber TEXT,
    phoneNumber TEXT NOT NULL,
    memberIdNumber TEXT,
    dateJoined TEXT DEFAULT CURRENT_TIMESTAMP,
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
    homeTraditionalAuthority TEXT,
    homeDistrict TEXT,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberId) REFERENCES member (id) ON DELETE CASCADE
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
  IF NOT EXISTS memberDependant (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberId INTEGER NOT NULL,
    name TEXT,
    phoneNumber TEXT,
    address TEXT,
    percentage REAL,
    isNominee INTEGER DEFAULT 0,
    relationship TEXT,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberId) REFERENCES member (id) ON DELETE CASCADE
  );

CREATE TRIGGER IF NOT EXISTS memberDependantUpdated AFTER
UPDATE ON memberDependant FOR EACH ROW BEGIN
UPDATE memberDependant
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS memberSaving (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberId INTEGER NOT NULL,
    memberSavingIdNumber TEXT,
    monthlySaving REAL,
    totalSaving REAL,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberId) REFERENCES member (id) ON DELETE CASCADE
  );

CREATE TRIGGER IF NOT EXISTS memberSavingUpdated AFTER
UPDATE ON memberSaving FOR EACH ROW BEGIN
UPDATE memberSaving
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS memberLoan (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberId INTEGER NOT NULL,
    memberSavingId INTEGER,
    phoneNumber TEXT,
    loanAmount REAL,
    repaymentPeriodInMonths REAL,
    loanPurpose TEXT,
    loanType TEXT NOT NULL CHECK (
      loanType IN (
        'PERSONAL',
        'BUSINESS',
        'AGRICULTURAL',
        'EMERGENCY'
      )
    ),
    loanStartDate TEXT,
    loanDueDate TEXT,
    monthlyInstallments REAL,
    interestRate REAL,
    amountPaid REAL,
    balanceAmount REAL,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberId) REFERENCES member (id) ON DELETE CASCADE
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
  IF NOT EXISTS memberBusiness (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberLoanId INTEGER NOT NULL,
    yearsInBusiness REAL,
    businessNature TEXT,
    businessName TEXT,
    tradingArea TEXT,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberLoanId) REFERENCES memberLoan (id) ON DELETE CASCADE
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
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberBusinessId) REFERENCES memberBusiness (id) ON DELETE CASCADE
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
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberBusinessId) REFERENCES memberBusiness (id) ON DELETE CASCADE
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
  IF NOT EXISTS memberOccupation (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberLoanId INTEGER NOT NULL,
    employerName TEXT,
    employerAddress TEXT,
    employerPhone TEXT,
    jobTitle TEXT,
    periodEmployedInMonths REAL,
    grossPay REAL,
    netPay REAL,
    highestQualification TEXT,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberLoanId) REFERENCES memberLoan (id) ON DELETE CASCADE
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
  IF NOT EXISTS memberLoanApproval (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberLoanId INTEGER NOT NULL,
    loanStatus TEXT NOT NULL CHECK (loanStatus IN ('PENDING', 'APPROVED', 'REJECTED')),
    amountRecommended REAL,
    approvedBy TEXT,
    approvalDate TEXT,
    amountApproved REAL,
    verifiedBy TEXT,
    dateVerified TEXT,
    denialOrPartialReason TEXT,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberLoanId) REFERENCES memberLoan (id) ON DELETE CASCADE
  );

CREATE TRIGGER IF NOT EXISTS memberLoanApprovalUpdated AFTER
UPDATE ON memberLoanApproval FOR EACH ROW BEGIN
UPDATE memberLoanApproval
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS memberLoanLiability (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberLoanId INTEGER NOT NULL,
    description TEXT,
    value REAL,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberLoanId) REFERENCES memberLoan (id) ON DELETE CASCADE
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
    memberLoanId INTEGER NOT NULL,
    description TEXT,
    value REAL,
    serialNumber TEXT,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberLoanId) REFERENCES memberLoan (id) ON DELETE CASCADE
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
    memberLoanId INTEGER NOT NULL,
    name TEXT,
    telephone TEXT,
    address TEXT,
    date TEXT,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberLoanId) REFERENCES memberLoan (id) ON DELETE CASCADE
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
  IF NOT EXISTS memberOccupationVerification (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberOccupationId INTEGER NOT NULL,
    jobVerified TEXT NOT NULL CHECK (jobVerified IN ('Yes', 'No')) DEFAULT 'No',
    grossVerified TEXT NOT NULL CHECK (grossVerified IN ('Yes', 'No')) DEFAULT 'No',
    netVerified TEXT NOT NULL CHECK (netVerified IN ('Yes', 'No')) DEFAULT 'No',
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberOccupationId) REFERENCES memberOccupation (id) ON DELETE CASCADE
  );

CREATE TRIGGER IF NOT EXISTS employmentVerificationUpdated AFTER
UPDATE ON memberOccupationVerification FOR EACH ROW BEGIN
UPDATE memberOccupationVerification
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS memberBusinessVerification (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberBusinessId INTEGER NOT NULL,
    businessVerified TEXT NOT NULL CHECK (businessVerified IN ('Yes', 'No')) DEFAULT 'No',
    grossIncomeVerified TEXT NOT NULL CHECK (grossIncomeVerified IN ('Yes', 'No')) DEFAULT 'No',
    netIncomeVerified TEXT NOT NULL CHECK (netIncomeVerified IN ('Yes', 'No')) DEFAULT 'No',
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberBusinessId) REFERENCES memberBusiness (id) ON DELETE CASCADE
  );

CREATE TRIGGER IF NOT EXISTS memberBusinessVerificationUpdated AFTER
UPDATE ON memberBusinessVerification FOR EACH ROW BEGIN
UPDATE memberBusinessVerification
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS insuranceProvider (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    providerName TEXT,
    providerAddress TEXT,
    providerPhoneNumber TEXT,
    providerContact TEXT,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
  );

CREATE TRIGGER IF NOT EXISTS insuranceProviderUpdated AFTER
UPDATE ON insuranceProvider FOR EACH ROW BEGIN
UPDATE insuranceProvider
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS memberLoanInsurance (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberLoanId INTEGER NOT NULL,
    insuranceProviderId INTEGER NOT NULL,
    date TEXT,
    amount REAL,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberLoanId) REFERENCES memberLoan (id) ON DELETE CASCADE
  );

CREATE TRIGGER IF NOT EXISTS memberLoanInsuranceUpdated AFTER
UPDATE ON memberLoanInsurance FOR EACH ROW BEGIN
UPDATE memberLoanInsurance
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS memberLoanPaymentSchedule (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberLoanId INTEGER NOT NULL,
    dueDate TEXT,
    amountDue REAL,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberLoanId) REFERENCES memberLoan (id) ON DELETE CASCADE
  );

CREATE TRIGGER IF NOT EXISTS memberLoanPaymentScheduleUpdated AFTER
UPDATE ON memberLoanPaymentSchedule FOR EACH ROW BEGIN
UPDATE memberLoanPaymentSchedule
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

CREATE TABLE
  IF NOT EXISTS notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberId INTEGER NOT NULL,
    notificationDate TEXT,
    notificationText TEXT,
    notificationRead INTEGER DEFAULT 0,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberLoanId) REFERENCES memberLoan (id) ON DELETE CASCADE
  );

CREATE TRIGGER IF NOT EXISTS notificationsUpdated AFTER
UPDATE ON notifications FOR EACH ROW BEGIN
UPDATE notifications
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;

---- START addMemberIdNumber TRIGGER ----
CREATE TRIGGER IF NOT EXISTS addMemberIdNumber AFTER INSERT ON member FOR EACH ROW BEGIN
UPDATE memberIdsCache
SET
  claimed = 1,
  memberId = NEW.id
WHERE
  id = (
    SELECT
      id
    FROM
      memberIdsCache
    WHERE
      claimed = 0
    ORDER BY
      id
    LIMIT
      1
  );

UPDATE member
SET
  memberIdNumber = (
    SELECT
      idNumber
    FROM
      memberIdsCache
    WHERE
      memberId = NEW.id
  )
WHERE
  id = NEW.id;

---- END addMemberIdNumber TRIGGER ----
END;

---- START addMemberIdNumber TRIGGER ----
CREATE TRIGGER IF NOT EXISTS addMemberIdNumber AFTER INSERT ON member FOR EACH ROW BEGIN
UPDATE memberSavingIdsCache
SET
  claimed = 1,
  memberSavingId = NEW.id
WHERE
  id = (
    SELECT
      id
    FROM
      memberSavingIdsCache
    WHERE
      claimed = 0
    ORDER BY
      id
    LIMIT
      1
  );

UPDATE member
SET
  memberSavingIdNumber = (
    SELECT
      idNumber
    FROM
      memberSavingIdsCache
    WHERE
      memberSavingId = NEW.id
  )
WHERE
  id = NEW.id;

---- END addMemberIdNumber TRIGGER ----
END;

INSERT
OR IGNORE INTO userRole (name)
VALUES
  ("Default"),
  ("Member"),
  ("Admin"),
  ("Cashier"),
  ("Accountant"),
  ("Loans Officer"),
  ("Manager");

INSERT
OR IGNORE INTO user (username, password, name, role)
VALUES
  (
    "default",
    "$2a$10$Xo4x3KiCkB3xGKvaCI4Hn.Be95DEiaIT3lbvHx/kOmyx7IqGY6ILK",
    "Default User",
    "Default"
  );

WITH RECURSIVE
  cnt (x) AS (
    SELECT
      1
    UNION ALL
    SELECT
      x + 1
    FROM
      cnt
    LIMIT
      999999
  ) INSERT
  OR IGNORE INTO memberIdsCache (idNumber)
SELECT
  CONCAT ('KSM', SUBSTR ('000000' || x, -6)) AS id
FROM
  cnt
ORDER BY
  RANDOM ();

WITH RECURSIVE
  cnt (x) AS (
    SELECT
      1
    UNION ALL
    SELECT
      x + 1
    FROM
      cnt
    LIMIT
      999999
  ) INSERT
  OR IGNORE INTO memberSavingIdsCache (idNumber)
SELECT
  CONCAT ('KSS', SUBSTR ('000000' || x, -6)) AS id
FROM
  cnt
ORDER BY
  RANDOM ();