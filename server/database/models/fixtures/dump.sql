PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE member (
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
INSERT INTO member VALUES(1,'Mary','Banda','','Female','Miss','Single','1999-09-01','DHFYR8475','ESCOM','29383746','','','09999999999',NULL,NULL,NULL,1,'2025-09-01 06:48:29','2025-09-01 06:48:29');
CREATE TABLE memberContact (
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
INSERT INTO memberContact VALUES(1,1,'P.O. Box 3200, Blantyre','Chilomoni, Blantrye',NULL,'Thumba','Kabudula','Lilongwe',1,'2025-09-01 06:48:29','2025-09-01 06:48:29');
INSERT INTO memberContact VALUES(2,1,'P.O. Box 1000, Lilongwe','Area 2, Lilongwe',NULL,'Songwe','Kyungu','Karonga',1,'2025-09-01 06:48:29','2025-09-01 06:48:29');
CREATE TABLE memberNominee (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberId INTEGER NOT NULL,
    name TEXT,
    phoneNumber TEXT,
    address TEXT,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberId) REFERENCES member (id) ON DELETE CASCADE
  );
INSERT INTO memberNominee VALUES(1,1,'John Banda','0888888888','Same as member',1,'2025-09-01 06:48:29','2025-09-01 06:48:29');
CREATE TABLE memberOccupation (
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
INSERT INTO memberOccupation VALUES(1,1,'SOBO','Kanengo','0999888474','Driver',36.0,100000.0,90000.0,'Secondary',1,'2025-09-01 06:48:29','2025-09-01 06:48:29');
CREATE TABLE memberBeneficiary (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberId INTEGER NOT NULL,
    name TEXT,
    percentage REAL,
    contact TEXT,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberId) REFERENCES member (id) ON DELETE CASCADE
  );
INSERT INTO memberBeneficiary VALUES(1,1,'Benefator 1',10.0,'P.O. Box 1',1,'2025-09-01 06:48:29','2025-09-01 06:48:29');
INSERT INTO memberBeneficiary VALUES(2,1,'Benefator 2',8.0,'P.O. Box 2',1,'2025-09-01 06:48:29','2025-09-01 06:48:29');
CREATE TABLE memberBusiness (
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
INSERT INTO memberBusiness VALUES(1,1,3.0,'Vendor','Vendors Galore','Mtandire',1,'2025-09-01 06:48:29','2025-09-01 06:48:29');
CREATE TABLE memberLastYearBusinessHistory (
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
INSERT INTO memberLastYearBusinessHistory VALUES(1,1,'2024',2000000.0,1000000.0,50000.0,100000.0,50000.0,0.0,35000.0,50000.0,0.0,1285000.0,715000.0,1,'2025-09-01 06:48:29','2025-09-01 06:48:29');
CREATE TABLE memberNextYearBusinessProjection (
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
INSERT INTO memberNextYearBusinessProjection VALUES(1,1,'2025',2500000.0,1500000.0,50000.0,100000.0,50000.0,0.0,35000.0,50000.0,0.0,1285000.0,715000.0,1,'2025-09-01 06:48:29','2025-09-01 06:48:29');
CREATE TABLE memberShares (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberId INTEGER NOT NULL,
    numberOfShares REAL,
    pricePerShare REAL,
    sharesType TEXT NOT NULL CHECK (sharesType IN ('FIXED', 'REDEEMABLE')),
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberId) REFERENCES member (id) ON DELETE CASCADE
  );
CREATE TABLE memberLoan (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberId INTEGER NOT NULL,
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
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberId) REFERENCES member (id) ON DELETE CASCADE
  );
INSERT INTO memberLoan VALUES(1,1,NULL,200000.0,12.0,'School fees','PERSONAL',1,'2025-09-01 06:48:29','2025-09-01 06:48:29');
CREATE TABLE memberLoanApproval (
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
INSERT INTO memberLoanApproval VALUES(1,1,'APPROVED',200000.0,'me','2025-08-30',200000.0,'me','2025-08-30',NULL,1,'2025-09-01 06:48:29','2025-09-01 06:48:29');
CREATE TABLE memberLoanLiability (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    memberLoanId INTEGER NOT NULL,
    description TEXT,
    value REAL,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (memberLoanId) REFERENCES memberLoan (id) ON DELETE CASCADE
  );
CREATE TABLE memberLoanSecurity (
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
CREATE TABLE memberLoanWitness (
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
CREATE TABLE memberOccupationVerification (
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
CREATE TABLE user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    name TEXT NOT NULL,
    role TEXT NOT NULL,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
  );
INSERT INTO user VALUES(1,'admin','$2a$10$Xo4x3KiCkB3xGKvaCI4Hn.Be95DEiaIT3lbvHx/kOmyx7IqGY6ILK','Admin User','admin',1,'2025-09-01 06:48:17','2025-09-01 06:48:17');
CREATE TABLE role (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
  );
INSERT INTO role VALUES(1,'Member',1,'2025-09-01 06:48:17','2025-09-01 06:48:17');
INSERT INTO role VALUES(2,'Admin',1,'2025-09-01 06:48:17','2025-09-01 06:48:17');
INSERT INTO role VALUES(3,'Cashier',1,'2025-09-01 06:48:17','2025-09-01 06:48:17');
INSERT INTO role VALUES(4,'Accountant',1,'2025-09-01 06:48:17','2025-09-01 06:48:17');
INSERT INTO role VALUES(5,'Loans Officer',1,'2025-09-01 06:48:17','2025-09-01 06:48:17');
INSERT INTO role VALUES(6,'Manager',1,'2025-09-01 06:48:17','2025-09-01 06:48:17');
CREATE TABLE account (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    accountName TEXT NOT NULL UNIQUE,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
  );
CREATE TABLE accountTransaction (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    accountTransactionDate TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    description text NOT NULL,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
  );
CREATE TABLE accountJournal (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    accountTransactionId INTEGER NOT NULL,
    accountId INTEGER NOT NULL,
    debit REAL NOT NULL,
    credit REAL NOT NULL,
    active INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (accountTransactionId) REFERENCES accountTransaction (id) ON DELETE CASCADE,
    FOREIGN KEY (accountId) REFERENCES account (id) ON DELETE CASCADE,
    CHECK (
      debit >= 0
      AND credit >= 0
    ),
    CHECK (
      debit > 0
      OR credit > 0
    )
  );
DELETE FROM sqlite_sequence;
INSERT INTO sqlite_sequence VALUES('role',6);
INSERT INTO sqlite_sequence VALUES('user',1);
INSERT INTO sqlite_sequence VALUES('member',1);
INSERT INTO sqlite_sequence VALUES('memberContact',2);
INSERT INTO sqlite_sequence VALUES('memberNominee',1);
INSERT INTO sqlite_sequence VALUES('memberBeneficiary',2);
INSERT INTO sqlite_sequence VALUES('memberLoan',1);
INSERT INTO sqlite_sequence VALUES('memberBusiness',1);
INSERT INTO sqlite_sequence VALUES('memberLastYearBusinessHistory',1);
INSERT INTO sqlite_sequence VALUES('memberNextYearBusinessProjection',1);
INSERT INTO sqlite_sequence VALUES('memberLoanApproval',1);
INSERT INTO sqlite_sequence VALUES('memberOccupation',1);
CREATE TRIGGER memberUpdated AFTER
UPDATE ON member FOR EACH ROW BEGIN
UPDATE member
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;
CREATE TRIGGER memberContactUpdated AFTER
UPDATE ON memberContact FOR EACH ROW BEGIN
UPDATE memberContact
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;
CREATE TRIGGER memberNomineeUpdated AFTER
UPDATE ON memberNominee FOR EACH ROW BEGIN
UPDATE memberNominee
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;
CREATE TRIGGER memberOccupationUpdated AFTER
UPDATE ON memberOccupation FOR EACH ROW BEGIN
UPDATE memberOccupation
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;
CREATE TRIGGER memberBeneficiaryUpdated AFTER
UPDATE ON memberBeneficiary FOR EACH ROW BEGIN
UPDATE memberBeneficiary
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;
CREATE TRIGGER memberBusinessUpdated AFTER
UPDATE ON memberBusiness FOR EACH ROW BEGIN
UPDATE memberBusiness
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;
CREATE TRIGGER memberLastYearBusinessHistoryUpdated AFTER
UPDATE ON memberLastYearBusinessHistory FOR EACH ROW BEGIN
UPDATE memberLastYearBusinessHistory
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;
CREATE TRIGGER memberNextYearBusinessProjectionUpdated AFTER
UPDATE ON memberNextYearBusinessProjection FOR EACH ROW BEGIN
UPDATE memberNextYearBusinessProjection
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;
CREATE TRIGGER memberSharesUpdated AFTER
UPDATE ON memberShares FOR EACH ROW BEGIN
UPDATE memberShares
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;
CREATE TRIGGER memberLoanUpdated AFTER
UPDATE ON memberLoan FOR EACH ROW BEGIN
UPDATE memberLoan
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;
CREATE TRIGGER memberLoanApprovalUpdated AFTER
UPDATE ON memberLoanApproval FOR EACH ROW BEGIN
UPDATE memberLoanApproval
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;
CREATE TRIGGER memberLoanLiabilityUpdated AFTER
UPDATE ON memberLoanLiability FOR EACH ROW BEGIN
UPDATE memberLoanLiability
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;
CREATE TRIGGER memberLoanSecurityUpdated AFTER
UPDATE ON memberLoanSecurity FOR EACH ROW BEGIN
UPDATE memberLoanSecurity
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;
CREATE TRIGGER memberLoanWitnessUpdated AFTER
UPDATE ON memberLoanWitness FOR EACH ROW BEGIN
UPDATE memberLoanWitness
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;
CREATE TRIGGER employmentVerificationUpdated AFTER
UPDATE ON memberOccupationVerification FOR EACH ROW BEGIN
UPDATE memberOccupationVerification
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;
CREATE TRIGGER userUpdated AFTER
UPDATE ON user FOR EACH ROW BEGIN
UPDATE user
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;
CREATE TRIGGER accountUpdated AFTER
UPDATE ON account FOR EACH ROW BEGIN
UPDATE account
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;
CREATE TRIGGER accountTransactionUpdated AFTER
UPDATE ON accountTransaction FOR EACH ROW BEGIN
UPDATE accountTransaction
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;
CREATE TRIGGER accountJournalUpdated AFTER
UPDATE ON accountJournal FOR EACH ROW BEGIN
UPDATE accountJournal
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = OLD.id;

END;
COMMIT;
