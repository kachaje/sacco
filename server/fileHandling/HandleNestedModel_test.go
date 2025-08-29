package filehandling_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sacco/server/database"
	filehandling "sacco/server/fileHandling"
	"sacco/server/parser"
	"testing"
)

func TestSimpleNestedModel(t *testing.T) {
	phoneNumber := "0999888777"
	sourceFolder := filepath.Join("..", "database", "models", "fixtures", "cache", phoneNumber)
	cacheFolder := filepath.Join(".", "tmp15", "cache")

	sessionFolder := filepath.Join(cacheFolder, phoneNumber)

	os.MkdirAll(filepath.Join(cacheFolder, phoneNumber), 0755)

	defer func() {
		os.RemoveAll(filepath.Join(".", "tmp15"))
	}()

	for _, file := range []string{
		"memberOccupation.27395048-84f4-11f0-9d0e-1e4d4999250c.json",
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
		GlobalIds:   map[string]int64{},
	}

	sessions := make(map[string]*parser.Session)

	sessions[phoneNumber] = session

	count := 0

	saveFunc := func(
		a map[string]any,
		b string,
		c int,
	) (*int64, error) {
		var id int64 = 13

		count++

		return &id, nil
	}

	model := "member"

	sessions[phoneNumber].AddedModels["memberContact"] = true

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

	err := filehandling.HandleNestedModel(data, &model, &phoneNumber, &cacheFolder, saveFunc, sessions, sessionFolder, nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range []string{
		"memberContact.158a2d54-84f4-11f0-8e0d-1e4d4999250c.json",
		"memberOccupation.27395048-84f4-11f0-9d0e-1e4d4999250c.json",
		"memberBeneficiary.fd40d7de-84f3-11f0-9b12-1e4d4999250c.json",
		"memberNominee.1efda9a6-84f4-11f0-8797-1e4d4999250c.json",
	} {
		filename := filepath.Join(cacheFolder, phoneNumber, file)

		_, err = os.Stat(filename)
		if !os.IsNotExist(err) {
			t.Fatalf("Test failed. Expected file %s to be deleted by now", filename)
		}
	}

	if count != 2 {
		t.Fatalf("Test failed. Expected: 2; Actual: %v", count)
	}
}

func TestComplexNestedModel(t *testing.T) {
	phoneNumber := "0999888777"
	sourceFolder := filepath.Join("..", "database", "models", "fixtures", "cache", phoneNumber)
	cacheFolder := filepath.Join(".", "tmp10", "cache")

	sessionFolder := filepath.Join(cacheFolder, phoneNumber)

	os.MkdirAll(filepath.Join(cacheFolder, phoneNumber), 0755)

	defer func() {
		os.RemoveAll(filepath.Join(".", "tmp10"))
	}()

	for _, file := range []string{
		"memberContact.158a2d54-84f4-11f0-8e0d-1e4d4999250c.json",
		"memberOccupation.27395048-84f4-11f0-9d0e-1e4d4999250c.json",
		"memberBeneficiary.fd40d7de-84f3-11f0-9b12-1e4d4999250c.json",
		"memberNominee.1efda9a6-84f4-11f0-8797-1e4d4999250c.json",
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
		AddedModels: map[string]bool{
			"memberContact":     true,
			"memberOccupation":  true,
			"memberBeneficiary": true,
			"memberNominee":     true,
		},
		GlobalIds: map[string]int64{},
	}

	sessions := make(map[string]*parser.Session)

	sessions[phoneNumber] = session

	count := 0

	saveFunc := func(
		a map[string]any,
		b string,
		c int,
	) (*int64, error) {
		count++

		var id int64 = int64(count)

		a["id"] = id

		return &id, nil
	}

	model := "member"

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

	err := filehandling.HandleNestedModel(data, &model, &phoneNumber, &cacheFolder, saveFunc, sessions, sessionFolder, nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range []string{
		"memberContact.158a2d54-84f4-11f0-8e0d-1e4d4999250c.json",
		"memberOccupation.27395048-84f4-11f0-9d0e-1e4d4999250c.json",
		"memberBeneficiary.fd40d7de-84f3-11f0-9b12-1e4d4999250c.json",
		"memberNominee.1efda9a6-84f4-11f0-8797-1e4d4999250c.json",
	} {
		filename := filepath.Join(cacheFolder, phoneNumber, file)

		_, err = os.Stat(filename)
		if !os.IsNotExist(err) {
			t.Fatalf("Test failed. Expected file %s to be deleted by now", filename)
		}
	}

	if count != 6 {
		t.Fatalf("Test failed. Expected: 6; Actual: %v", count)
	}

	target := map[string]int64{
		"memberBeneficiaryId": 6,
		"memberContactId":     2,
		"memberId":            1,
		"memberNomineeId":     3,
		"memberOccupationId":  4,
	}

	if !reflect.DeepEqual(target, session.GlobalIds) {
		t.Fatalf("Test failed")
	}
}

func TestChildNestedModel(t *testing.T) {
	phoneNumber := "0999888777"
	cacheFolder := filepath.Join(".", "tmp11", "cache")

	sessionFolder := filepath.Join(cacheFolder, phoneNumber)

	os.MkdirAll(filepath.Join(cacheFolder, phoneNumber), 0755)

	defer func() {
		os.RemoveAll(filepath.Join(".", "tmp11"))
	}()

	session := &parser.Session{
		AddedModels: map[string]bool{},
		GlobalIds: map[string]int64{
			"memberId":     16,
			"memberLoanId": 13,
		},
	}

	sessions := make(map[string]*parser.Session)

	sessions[phoneNumber] = session

	count := 0

	saveFunc := func(
		data map[string]any,
		model string,
		retries int,
	) (*int64, error) {
		if data["memberId"] == nil {
			return nil, fmt.Errorf("missing required field memberId")
		}
		if data["memberLoanId"] == nil {
			return nil, fmt.Errorf("missing required field memberLoanId")
		}

		count++

		var id int64 = int64(count)

		data["id"] = id

		return &id, nil
	}

	model := "memberOccupation"

	sessions[phoneNumber].AddedModels["member"] = true

	data := map[string]any{
		"employerAddress":        "Kanengo",
		"employerName":           "SOBO",
		"employerPhone":          "01282373737",
		"grossPay":               100000,
		"highestQualification":   "Secondary",
		"jobTitle":               "Driver",
		"netPay":                 90000,
		"periodEmployedInMonths": "36",
	}

	err := filehandling.HandleNestedModel(data, &model, &phoneNumber, &cacheFolder, saveFunc, sessions, sessionFolder, nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range []string{
		"memberContact.158a2d54-84f4-11f0-8e0d-1e4d4999250c.json",
		"memberOccupation.27395048-84f4-11f0-9d0e-1e4d4999250c.json",
		"memberBeneficiary.fd40d7de-84f3-11f0-9b12-1e4d4999250c.json",
		"memberNominee.1efda9a6-84f4-11f0-8797-1e4d4999250c.json",
	} {
		filename := filepath.Join(cacheFolder, phoneNumber, file)

		_, err = os.Stat(filename)
		if !os.IsNotExist(err) {
			t.Fatalf("Test failed. Expected file %s to be deleted by now", filename)
		}
	}

	if count != 1 {
		t.Fatalf("Test failed. Expected: 1; Actual: %v", count)
	}

	target := map[string]int64{
		"memberOccupationId": 1,
		"memberId":           16,
		"memberLoanId":       13,
	}

	if !reflect.DeepEqual(target, session.GlobalIds) {
		t.Fatalf("Test failed")
	}
}

func TestUnpackData(t *testing.T) {
	data := map[string]any{}
	target := []map[string]any{}

	for i := range 4 {
		row := map[string]any{}

		for _, key := range []string{"id", "name", "value"} {
			label := fmt.Sprintf("%s%d", key, i+1)
			value := fmt.Sprintf("%s%d", key, i+1)

			data[label] = value
			row[key] = value
		}

		target = append(target, row)
	}

	result := filehandling.UnpackData(data)

	if len(result) != len(target) {
		t.Fatalf("Test failed. Expected: %v; Actual: %v", len(target), len(result))
	}

	if len(result[0]) != len(target[0]) {
		t.Fatalf("Test failed. Expected: %v; Actual: %v", len(target[0]), len(result[0]))
	}

	data = map[string]any{
		"id":    "1",
		"name":  "test",
		"value": "something",
	}
	target = []map[string]any{}

	target = append(target, data)

	result = filehandling.UnpackData(data)

	if !reflect.DeepEqual(target, result) {
		t.Fatal("Test failed")
	}
}

func TestArrayChildData(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer func() {
		db.Close()
	}()

	phoneNumber := "0999888777"
	cacheFolder := "./tmpArrData"

	sessionFolder := filepath.Join(cacheFolder, phoneNumber)

	os.MkdirAll(filepath.Join(cacheFolder, phoneNumber), 0755)

	defer func() {
		os.RemoveAll(filepath.Join(cacheFolder))
	}()

	sessions := map[string]*parser.Session{
		phoneNumber: {
			GlobalIds:   map[string]int64{},
			ActiveData:  map[string]any{},
			AddedModels: map[string]bool{},
		},
	}

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

	model := "member"

	err := filehandling.HandleNestedModel(data, &model, &phoneNumber, &cacheFolder, db.GenericsSaveData, sessions, sessionFolder, nil)
	if err != nil {
		t.Fatal(err)
	}

	data = map[string]any{
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

	model = "memberBeneficiary"

	err = filehandling.HandleNestedModel(data, &model, &phoneNumber, &cacheFolder, db.GenericsSaveData, sessions, sessionFolder, nil)
	if err != nil {
		t.Fatal(err)
	}

	if !sessions[phoneNumber].AddedModels["memberBeneficiary"] {
		t.Fatalf("Test failed. Expected: true; Actual: %v",
			sessions[phoneNumber].AddedModels["memberBeneficiary"])
	}

	result, err := db.GenericModels["memberBeneficiary"].FilterBy("WHERE active=1")
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("Test failed. Expected: 2; Actual: %v", len(result))
	}
}

func TestGetSkippedRefIds(t *testing.T) {
	refData := []map[string]any{
		{
			"contact":    "P.O. Box 2",
			"id":         2,
			"memberId":   1,
			"name":       "Benefator 2",
			"percentage": 8,
		},
		{
			"contact":    "P.O. Box 3",
			"id":         3,
			"memberId":   1,
			"name":       "Benefator 3",
			"percentage": 5,
		},
		{
			"contact":    "P.O. Box 4",
			"id":         4,
			"memberId":   1,
			"name":       "Benefator 4",
			"percentage": 2,
		},
		{
			"contact":    "P.O. Box 1",
			"id":         1,
			"memberId":   1,
			"name":       "Benefator 1",
			"percentage": 10,
		},
	}
	data := []map[string]any{
		{
			"contact":    "P.O. Box 5678",
			"id":         2,
			"memberId":   1,
			"name":       "Benefator 2",
			"percentage": 25,
		},
		{
			"contact":    "P.O. Box 1234",
			"id":         1,
			"memberId":   1,
			"name":       "Benefator 1",
			"percentage": 35,
		},
	}

	result := filehandling.GetSkippedRefIds(data, refData)

	target := []map[string]any{
		{
			"contact":    "P.O. Box 3",
			"id":         3,
			"memberId":   1,
			"name":       "Benefator 3",
			"percentage": 5},
		{
			"contact":    "P.O. Box 4",
			"id":         4,
			"memberId":   1,
			"name":       "Benefator 4",
			"percentage": 2,
		},
	}

	if !reflect.DeepEqual(target, result) {
		t.Fatal("Test failed")
	}
}
