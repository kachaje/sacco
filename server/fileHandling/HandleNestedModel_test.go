package filehandling_test

import (
	"io"
	"os"
	"path/filepath"
	"reflect"
	filehandling "sacco/server/fileHandling"
	"sacco/server/parser"
	"testing"
)

func TestSimpleNestedModel(t *testing.T) {
	phoneNumber := "0999888777"
	sourceFolder := filepath.Join("..", "database", "models", "fixtures", "cache", phoneNumber)
	cacheFolder := filepath.Join(".", "tmp5", "cache")

	sessionFolder := filepath.Join(cacheFolder, phoneNumber)

	os.MkdirAll(filepath.Join(cacheFolder, phoneNumber), 0755)

	defer func() {
		os.RemoveAll(filepath.Join(".", "tmp5"))
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

	err := filehandling.HandleNestedModel(data, &model, &phoneNumber, &cacheFolder, saveFunc, sessions, sessionFolder)
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

	if count != 2 {
		t.Fatalf("Test failed. Expected: 2; Actual: %v", count)
	}
}

func TestComplexNestedModel(t *testing.T) {
	phoneNumber := "0999888777"
	sourceFolder := filepath.Join("..", "database", "models", "fixtures", "cache", phoneNumber)
	cacheFolder := filepath.Join(".", "tmp1", "cache")

	sessionFolder := filepath.Join(cacheFolder, phoneNumber)

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

	err := filehandling.HandleNestedModel(data, &model, &phoneNumber, &cacheFolder, saveFunc, sessions, sessionFolder)
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
