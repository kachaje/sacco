INSERT INTO
  member (
    id,
    firstName,
    lastName,
    otherName,
    gender,
    title,
    maritalStatus,
    dateOfBirth,
    nationalId,
    utilityBillType,
    utilityBillNumber,
    fileNumber,
    oldFileNumber
  )
VALUES
  (
    1,
    "Mary",
    "Banda",
    "",
    "Female",
    "Miss",
    "Single",
    "1999-09-01",
    "DHFYR8475",
    "ESCOM",
    "29383746",
    "",
    ""
  );

INSERT INTO
  memberContact (
    memberId,
    postalAddress,
    residentialAddress,
    phoneNumber,
    homeVillage,
    homeTA,
    homeDistrict
  )
VALUES
  (
    1,
    "P.O. Box 1000, Lilongwe",
    "Area 2, Lilongwe",
    "09999999999",
    "Songwe",
    "Kyungu",
    "Karonga"
  );

INSERT INTO
  memberNominee (
    memberId,
    nextOfKinName,
    nextOfKinPhone,
    nextOfKinAddress
  )
VALUES
  (1, "John Banda", "0888888888", "Same as member");

INSERT INTO
  memberOccupation (
    memberId,
    employerName,
    netPay,
    jobTitle,
    employerAddress,
    highestQualification
  )
VALUES
  (
    1,
    "Sunseed Oil",
    72000,
    "Driver",
    "Kanengo, Lilongwe",
    "Secondary"
  );

INSERT INTO
  memberBeneficiary (memberId, name, percentage, contact)
VALUES
  (1, "Benefator 1", 10, "P.O. Box 1"),
  (1, "Benefator 2", 8, "P.O. Box 2"),
  (1, "Benefator 3", 5, "P.O. Box 3"),
  (1, "Benefator 4", 2, "P.O. Box 4");