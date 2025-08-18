package server_test

import (
	"io"
	"os"
	"path/filepath"
	"sacco/server"
	"sacco/server/menus"
	"testing"
)

func TestSaveData(t *testing.T) {
	phoneNumber := "0999888777"
	sourceFolder := filepath.Join(".", "database", "models", "fixtures", "cache", phoneNumber)
	cacheFolder := filepath.Join(".", "tmp", "cache", phoneNumber)

	os.MkdirAll(cacheFolder, 0755)

	for _, file := range []string{"contactDetails.json"} {
		src, err := os.Open(filepath.Join(sourceFolder, file))
		if err != nil {
			continue
		}
		defer src.Close()

		dst, err := os.Create(filepath.Join(cacheFolder, file))
		if err != nil {
			continue
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			continue
		}
	}

	session := &menus.Session{
		ContactsAdded: true,
	}

	sessionId := "sample"

	sessions := map[string]*menus.Session{
		sessionId: session,
	}

	saveFunc := func(
		a map[string]any,
		b map[string]any,
		c map[string]any,
		d map[string]any,
		e []map[string]any,
		f *int64,
	) (*int64, error) {
		_ = sessions

		return nil, nil
	}

	data := map[string]any{}
	model := "memberDetails"

	err := server.SaveData(data, &model, &phoneNumber, &sessionId, &cacheFolder, nil, saveFunc)
	if err != nil {
		t.Fatal(err)
	}
}
