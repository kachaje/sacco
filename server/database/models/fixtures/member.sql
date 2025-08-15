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
    10,
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
    10,
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
  (10, "John Banda", "0888888888", "Same as member");