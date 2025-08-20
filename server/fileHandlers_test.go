package server_test

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sacco/parser"
	"sacco/server"
	"testing"
)

func TestSaveDataAll(t *testing.T) {
	phoneNumber := "0999888777"
	sourceFolder := filepath.Join(".", "database", "models", "fixtures", "cache", phoneNumber)
	cacheFolder := filepath.Join(".", "tmp1", "cache")

	os.MkdirAll(filepath.Join(cacheFolder, phoneNumber), 0755)

	defer func() {
		os.RemoveAll(filepath.Join(".", "tmp1"))
	}()

	for _, file := range []string{
		"contactDetails.json",
		"occupationDetails.json",
		"beneficiaries.json",
		"nomineeDetails.json",
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
		b map[string]any,
		c map[string]any,
		d map[string]any,
		e []map[string]any,
		f *int64,
	) (*int64, error) {
		var id int64 = 13

		return &id, nil
	}

	model := "memberDetails"

	err := server.SaveData(map[string]any{}, &model, &phoneNumber, &sessionId, &cacheFolder, nil, saveFunc, sessions, nil)
	if err != nil {
		t.Fatal(err)
	}

	sessions[sessionId].ContactsAdded = true

	data := map[string]any{
		"dateOfBirth":        "1999-09-01",
		"defaultPhoneNumber": "09999999999",
		"fileNumber":         "",
		"firstName":          "Mary",
		"gender":             "Female",
		"id":                 1,
		"lastName":           "Banda",
		"maritalStatus":      "Single",
		"nationalId":         "DHFYR8475",
		"oldFileNumber":      "",
		"otherName":          "",
		"title":              "Miss",
		"utilityBillNumber":  "29383746",
		"utilityBillType":    "ESCOM",
	}

	err = server.SaveData(data, &model, &phoneNumber, &sessionId, &cacheFolder, nil, saveFunc, sessions, nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range []string{
		"contactDetails.json",
		"occupationDetails.json",
		"beneficiaries.json",
		"nomineeDetails.json",
	} {
		filename := filepath.Join(cacheFolder, phoneNumber, file)

		_, err = os.Stat(filename)
		if !os.IsNotExist(err) {
			t.Fatalf("Test failed. Expected file %s to be deleted by now", filename)
		}
	}
}

func TestSaveDataOne(t *testing.T) {
	phoneNumber := "0999888777"
	sourceFolder := filepath.Join(".", "database", "models", "fixtures", "cache", phoneNumber)
	cacheFolder := filepath.Join(".", "tmp2", "cache")

	os.MkdirAll(filepath.Join(cacheFolder, phoneNumber), 0755)

	defer func() {
		os.RemoveAll(filepath.Join(".", "tmp2"))
	}()

	for _, file := range []string{
		"occupationDetails.json",
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
		b map[string]any,
		c map[string]any,
		d map[string]any,
		e []map[string]any,
		f *int64,
	) (*int64, error) {
		var id int64 = 13

		return &id, nil
	}

	model := "memberDetails"

	err := server.SaveData(map[string]any{}, &model, &phoneNumber, &sessionId, &cacheFolder, nil, saveFunc, sessions, nil)
	if err != nil {
		t.Fatal(err)
	}

	sessions[sessionId].OccupationAdded = true

	data := map[string]any{
		"dateOfBirth":        "1999-09-01",
		"defaultPhoneNumber": "09999999999",
		"fileNumber":         "",
		"firstName":          "Mary",
		"gender":             "Female",
		"id":                 1,
		"lastName":           "Banda",
		"maritalStatus":      "Single",
		"nationalId":         "DHFYR8475",
		"oldFileNumber":      "",
		"otherName":          "",
		"title":              "Miss",
		"utilityBillNumber":  "29383746",
		"utilityBillType":    "ESCOM",
	}

	err = server.SaveData(data, &model, &phoneNumber, &sessionId, &cacheFolder, nil, saveFunc, sessions, nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range []string{
		"contactDetails.json",
		"occupationDetails.json",
		"beneficiaries.json",
		"nomineeDetails.json",
	} {
		filename := filepath.Join(cacheFolder, phoneNumber, file)

		_, err = os.Stat(filename)
		if !os.IsNotExist(err) {
			t.Fatalf("Test failed. Expected file %s to be deleted by now", filename)
		}
	}
}

func TestRerunFailedSaves(t *testing.T) {
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

	err = os.WriteFile(filepath.Join(sessionFolder, "beneficiaries.json"), payload, 0644)
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

	err = os.WriteFile(filepath.Join(sessionFolder, "contactDetails.json"), payload, 0644)
	if err != nil {
		t.Fatal(err)
	}

	count := 0

	saveFunc := func(
		a map[string]any,
		b map[string]any,
		c map[string]any,
		d map[string]any,
		e []map[string]any,
		f *int64,
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

	err = server.RerunFailedSaves(&phoneNumber, &sessionId, &cacheFolder, saveFunc, sessions)
	if err != nil {
		t.Fatal(err)
	}

	if count != 2 {
		t.Fatalf("Test failed. Expected: 2; Actual: %v", count)
	}
}
