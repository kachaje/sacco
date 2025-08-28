package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

func main() {
	folder := filepath.Join("..")

	err := filepath.WalkDir(folder, func(file string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		basename := filepath.Base(file)

		if !strings.HasSuffix(basename, "_fn.go") {
			return nil
		}

		fnName := strings.TrimSuffix(basename, "_fn.go")

		fmt.Println(fnName)

		return nil
	})
	if err != nil {
		panic(err)
	}
}
