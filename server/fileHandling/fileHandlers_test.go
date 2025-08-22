package filehandling_test

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sacco/server/database"
	filehandling "sacco/server/fileHandling"
	"sacco/server/parser"
	"sacco/utils"
	"testing"
)

func TestSaveDataOne(t *testing.T) {
	phoneNumber := "0999888777"
	sourceFolder := filepath.Join("..", "database", "models", "fixtures", "cache", phoneNumber)
	cacheFolder := filepath.Join(".", "tmp2", "cache")

	os.MkdirAll(filepath.Join(cacheFolder, phoneNumber), 0755)

	defer func() {
		os.RemoveAll(filepath.Join(".", "tmp2"))
	}()

	for _, file := range []string{
		"memberOccupation.json",
	} {
		src, err := os.Open(filepath.Join(sourceFolder, file))
		if err != nil {
			t.Fatal(err)
			continue
		}
		defer src.Close()

		dst, err := os.Create(filepath.Join(cacheFolder, phoneNumber, file))
		if err != nil {
			t.Fatal(err)
			continue
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			t.Fatal(err)
			continue
		}

		_, err = os.Stat(dst.Name())
		if os.IsNotExist(err) {
			t.Fatalf("Test failed. Failed to create %s", dst.Name())
		}
	}

	session := &parser.Session{
		AddedModels: map[string]bool{},
	}

	sessionId := "sample"

	sessions := make(map[string]*parser.Session)

	sessions[sessionId] = session

	saveFunc := func(
		a map[string]any,
		b string,
		c int,
	) (*int64, error) {
		var id int64 = 13

		return &id, nil
	}

	model := "member"

	err := filehandling.SaveData(map[string]any{}, &model, &phoneNumber, &sessionId, &cacheFolder, nil, saveFunc, sessions, nil)
	if err != nil {
		t.Fatal(err)
	}

	sessions[sessionId].AddedModels["memberContact"] = true

	data := map[string]any{
		"dateOfBirth":       "1999-09-01",
		"phoneNumber":       "09999999999",
		"fileNumber":        "",
		"firstName":         "Mary",
		"gender":            "Female",
		"id":                1,
		"lastName":          "Banda",
		"maritalStatus":     "Single",
		"nationalId":        "DHFYR8475",
		"oldFileNumber":     "",
		"otherName":         "",
		"title":             "Miss",
		"utilityBillNumber": "29383746",
		"utilityBillType":   "ESCOM",
	}

	err = filehandling.SaveData(data, &model, &phoneNumber, &sessionId, &cacheFolder, nil, saveFunc, sessions, nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range []string{
		"memberContact.json",
		"memberOccupation.json",
		"memberBeneficiary.json",
		"memberNominee.json",
	} {
		filename := filepath.Join(cacheFolder, phoneNumber, file)

		_, err = os.Stat(filename)
		if !os.IsNotExist(err) {
			t.Fatalf("Test failed. Expected file %s to be deleted by now", filename)
		}
	}
}

func TestSaveDataAll(t *testing.T) {
	t.Skip()

	phoneNumber := "0999888777"
	sourceFolder := filepath.Join("..", "database", "models", "fixtures", "cache", phoneNumber)
	cacheFolder := filepath.Join(".", "tmp1", "cache")

	os.MkdirAll(filepath.Join(cacheFolder, phoneNumber), 0755)

	defer func() {
		os.RemoveAll(filepath.Join(".", "tmp1"))
	}()

	for _, file := range []string{
		"memberContact.json",
		"memberOccupation.json",
		"memberBeneficiary.json",
		"memberNominee.json",
	} {
		src, err := os.Open(filepath.Join(sourceFolder, file))
		if err != nil {
			t.Fatal(err)
			continue
		}
		defer src.Close()

		dst, err := os.Create(filepath.Join(cacheFolder, phoneNumber, file))
		if err != nil {
			t.Fatal(err)
			continue
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			t.Fatal(err)
			continue
		}

		_, err = os.Stat(dst.Name())
		if os.IsNotExist(err) {
			t.Fatalf("Test failed. Failed to create %s", dst.Name())
		}
	}

	session := &parser.Session{}

	sessionId := "sample"

	sessions := make(map[string]*parser.Session)

	sessions[sessionId] = session

	saveFunc := func(
		a map[string]any,
		b string,
		c int,
	) (*int64, error) {
		var id int64 = 13

		return &id, nil
	}

	model := "member"

	err := filehandling.SaveData(map[string]any{}, &model, &phoneNumber, &sessionId, &cacheFolder, nil, saveFunc, sessions, nil)
	if err != nil {
		t.Fatal(err)
	}

	sessions[sessionId].ContactsAdded = true

	data := map[string]any{
		"dateOfBirth":       "1999-09-01",
		"phoneNumber":       "09999999999",
		"fileNumber":        "",
		"firstName":         "Mary",
		"gender":            "Female",
		"id":                1,
		"lastName":          "Banda",
		"maritalStatus":     "Single",
		"nationalId":        "DHFYR8475",
		"oldFileNumber":     "",
		"otherName":         "",
		"title":             "Miss",
		"utilityBillNumber": "29383746",
		"utilityBillType":   "ESCOM",
	}

	err = filehandling.SaveData(data, &model, &phoneNumber, &sessionId, &cacheFolder, nil, saveFunc, sessions, nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range []string{
		"memberContact.json",
		"memberOccupation.json",
		"memberBeneficiary.json",
		"memberNominee.json",
	} {
		filename := filepath.Join(cacheFolder, phoneNumber, file)

		_, err = os.Stat(filename)
		if !os.IsNotExist(err) {
			t.Fatalf("Test failed. Expected file %s to be deleted by now", filename)
		}
	}

}

func TestRerunFailedSaves(t *testing.T) {
	t.Skip()

	phoneNumber := "0999888777"
	cacheFolder := filepath.Join(".", "tmpReruns", "cache")
	sessionFolder := filepath.Join(cacheFolder, phoneNumber)

	err := os.MkdirAll(sessionFolder, 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll(filepath.Join(".", "tmpReruns"))
	}()

	data := []map[string]any{
		{
			"contact":    "029293836",
			"memberId":   1,
			"name":       "John Banda",
			"percentage": 10,
		},
		{
			"contact":    "038373662",
			"memberId":   1,
			"name":       "Jean Phiri",
			"percentage": 6,
		},
	}

	payload, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(sessionFolder, "memberBeneficiary.json"), payload, 0644)
	if err != nil {
		t.Fatal(err)
	}

	payload, err = json.MarshalIndent(map[string]any{
		"homeDistrict":       "Karonga",
		"homeTA":             "Kyungu",
		"homeVillage":        "Songwe",
		"id":                 1,
		"memberId":           1,
		"phoneNumber":        "09999999999",
		"postalAddress":      "P.O. Box 1000, Lilongwe",
		"residentialAddress": "Area 2, Lilongwe",
	}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(sessionFolder, "memberContact.json"), payload, 0644)
	if err != nil {
		t.Fatal(err)
	}

	count := 0

	saveFunc := func(
		a map[string]any,
		b string,
		c int,
	) (*int64, error) {
		count++

		return nil, nil
	}

	var id int64 = 1

	session := &parser.Session{
		ActiveMemberData: map[string]any{},
		MemberId:         &id,
	}

	sessionId := "sample"

	sessions := make(map[string]*parser.Session)

	sessions[sessionId] = session

	err = filehandling.RerunFailedSaves(&phoneNumber, &sessionId, &cacheFolder, saveFunc, sessions)
	if err != nil {
		t.Fatal(err)
	}

	if count != 2 {
		t.Fatalf("Test failed. Expected: 2; Actual: %v", count)
	}
}

func TestHandleBeneficiaries(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer func() {
		db.Close()

		if _, err := os.Stat("memberBeneficiary.json"); !os.IsNotExist(err) {
			os.Remove("memberBeneficiary.json")
		}
	}()

	data := map[string]any{
		"contact1":    "P.O. Box 1234",
		"id1":         1,
		"memberId1":   1,
		"name1":       "Benefator 1",
		"percentage1": 35,
		"contact2":    "P.O. Box 5678",
		"id2":         2,
		"memberId2":   1,
		"name2":       "Benefator 2",
		"percentage2": 25,
	}

	phoneNumber := "0999888777"
	sessionId := "sample"
	cacheFolder := "./tmp"

	var id int64 = 1

	sessions := map[string]*parser.Session{
		sessionId: {
			MemberId:         &id,
			ActiveMemberData: map[string]any{},
			AddedModels:      map[string]bool{},
		},
	}

	err := filehandling.HandleBeneficiaries(data, &phoneNumber, &sessionId, &cacheFolder, db.GenericsSaveData, sessions, nil, "")
	if err != nil {
		t.Fatal(err)
	}

	if !sessions[sessionId].AddedModels["memberBeneficiary"] {
		t.Fatalf("Test failed. Expected: true; Actual: %v",
			sessions[sessionId].AddedModels["memberBeneficiary"])
	}

	result, err := db.GenericModels["memberBeneficiary"].FilterBy("WHERE active=1")
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("Test failed. Expected: 2; Actual: %v", len(result))
	}
}

func TestHandleMemberDetails(t *testing.T) {
	phoneNumber := "0999888777"
	sourceFolder := filepath.Join("..", "database", "models", "fixtures", "cache", phoneNumber)
	cacheFolder := filepath.Join(".", "tmp5", "cache")

	sessionFolder := filepath.Join(cacheFolder, phoneNumber)

	os.MkdirAll(sessionFolder, 0755)

	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer func() {
		db.Close()

		os.RemoveAll(filepath.Join(".", "tmp5"))
	}()

	for _, file := range []string{
		"memberContact.json",
		"memberOccupation.json",
		"memberBeneficiary.json",
		"memberNominee.json",
	} {
		src, err := os.Open(filepath.Join(sourceFolder, file))
		if err != nil {
			t.Fatal(err)
			continue
		}
		defer src.Close()

		dst, err := os.Create(filepath.Join(sessionFolder, file))
		if err != nil {
			t.Fatal(err)
			continue
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			t.Fatal(err)
			continue
		}

		_, err = os.Stat(dst.Name())
		if os.IsNotExist(err) {
			t.Fatalf("Test failed. Failed to create %s", dst.Name())
		}
	}

	data := map[string]any{
		"dateOfBirth":       "1999-09-01",
		"phoneNumber":       phoneNumber,
		"fileNumber":        "",
		"firstName":         "Mary",
		"gender":            "Female",
		"id":                1,
		"lastName":          "Banda",
		"maritalStatus":     "Single",
		"nationalId":        "DHFYR8475",
		"oldFileNumber":     "",
		"otherName":         "",
		"title":             "Miss",
		"utilityBillNumber": "29383746",
		"utilityBillType":   "ESCOM",
	}

	var id int64
	sessionId := "sample"

	sessions := map[string]*parser.Session{
		sessionId: {
			MemberId:         &id,
			ActiveMemberData: map[string]any{},
			AddedModels: map[string]bool{
				"memberContact":     true,
				"memberNominee":     true,
				"memberOccupation":  true,
				"memberBeneficiary": true,
			},
		},
	}

	err := filehandling.HandleMemberDetails(data, &phoneNumber, &sessionId, &cacheFolder, db.GenericsSaveData, sessions, sessionFolder)
	if err != nil {
		t.Fatal(err)
	}

	result, err := db.MemberByPhoneNumber(phoneNumber, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	target := map[string]any{
		"dateOfBirth":   "1999-09-01",
		"firstName":     "Mary",
		"gender":        "Female",
		"id":            1,
		"lastName":      "Banda",
		"maritalStatus": "Single",
		"memberBeneficiary": []map[string]any{
			{
				"contact":    "0888777444",
				"id":         1,
				"memberId":   1,
				"name":       "John Phiri",
				"percentage": 10,
			},
			{
				"contact":    "07746635653",
				"id":         2,
				"memberId":   1,
				"name":       "Jean Banda",
				"percentage": 5,
			},
		},
		"memberContact": map[string]any{
			"homeDistrict":       "Lilongwe",
			"homeTA":             "Kalolo",
			"homeVillage":        "Kalulu",
			"id":                 1,
			"memberId":           1,
			"postalAddress":      "P.O. Box 1",
			"residentialAddress": "Area 49",
		},
		"memberNominee": map[string]any{
			"id":             1,
			"memberId":       1,
			"nomineeAddress": "P.O. Box 1",
			"nomineeName":    "John Phiri",
			"nomineePhone":   "0888444666",
		},
		"memberOccupation": map[string]any{
			"employerAddress":      "Kanengo",
			"employerName":         "SOBO",
			"employerPhone":        "01282373737",
			"grossPay":             100000,
			"highestQualification": "Secondary",
			"id":                   1,
			"jobTitle":             "Driver",
			"memberId":             1,
			"netPay":               90000,
			"periodEmployed":       36,
		},
		"nationalId":        "DHFYR8475",
		"phoneNumber":       "0999888777",
		"title":             "Miss",
		"utilityBillNumber": "29383746",
		"utilityBillType":   "ESCOM",
	}

	payloadResult, _ := json.MarshalIndent(result, "", "  ")
	payloadTarget, _ := json.MarshalIndent(target, "", "  ")

	if utils.CleanScript(payloadResult) != utils.CleanScript(payloadTarget) {
		t.Fatal("Test failed")
	}
}
