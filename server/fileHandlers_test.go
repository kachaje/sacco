package server_test

import (
	"io"
	"os"
	"path/filepath"
	"sacco/parser"
	"sacco/server"
	"testing"
)

func TestSaveData(t *testing.T) {
	phoneNumber := "0999888777"
	sourceFolder := filepath.Join(".", "database", "models", "fixtures", "cache", phoneNumber)
	cacheFolder := filepath.Join(".", "tmp", "cache")

	os.MkdirAll(filepath.Join(cacheFolder, phoneNumber), 0755)

	defer func() {
		if false {
			os.RemoveAll(filepath.Join(".", "tmp"))
		}
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

	err := server.SaveData(map[string]any{}, &model, &phoneNumber, &sessionId, &cacheFolder, nil, saveFunc, sessions)
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

	err = server.SaveData(data, &model, &phoneNumber, &sessionId, &cacheFolder, nil, saveFunc, sessions)
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
