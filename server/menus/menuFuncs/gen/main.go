package main

import (
	"fmt"
	"go/format"
	"io/fs"
	"os"
	"path/filepath"
	"sacco/utils"
	"strings"
)

func main() {
	folder := filepath.Join("..")

	rows := []string{}

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

		cFnName := utils.CapitalizeFirstLetter(fnName)

		row := fmt.Sprintf(`FunctionsMap["%s"] = %s`, fnName, cFnName)

		rows = append(rows, row)

		return nil
	})
	if err != nil {
		panic(err)
	}

	content := fmt.Sprintf(`package menufuncs

import (
	"sacco/server/database"
	"sacco/server/parser"
)

var (
	DB       *database.Database
	Sessions = map[string]*parser.Session{}

	FunctionsMap = map[string]func(
		func(
			string, *parser.Session,
			string, string, string, string,
		) string,
		map[string]any,
	) string{}
)

func init() {
%s
}
`, strings.Join(rows, "\n"))

	filename := filepath.Join("..", "menufuncs.go")

	rawData, err := format.Source([]byte(content))
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(filename, rawData, 0644)
	if err != nil {
		panic(err)
	}
}
