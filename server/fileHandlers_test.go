package server_test

import (
	"io"
	"os"
	"path/filepath"
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
}
