package filehandling

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sacco/server/parser"
	"sacco/utils"
)

func RerunFailedSaves(phoneNumber, cacheFolder *string,
	saveFunc func(
		a map[string]any,
		b string,
		c int,
	) (*int64, error), sessions map[string]*parser.Session) error {
	sessionFolder := filepath.Join(*cacheFolder, *phoneNumber)

	_, err := os.Stat(sessionFolder)
	if !os.IsNotExist(err) {
		retryNumberLockfile := filepath.Join(sessionFolder, fmt.Sprintf("%s.number", *phoneNumber))

		if utils.FileLocked(retryNumberLockfile) {
			return nil
		}

		_, err = utils.LockFile(retryNumberLockfile)
		if err != nil {
			log.Printf("server.RerunFailedSaves.LockFile: %s", err.Error())
			return err
		}
		defer func() {
			err = utils.UnLockFile(retryNumberLockfile)
			if err != nil {
				log.Printf("server.RerunFailedSaves.UnLockFile: %s", err.Error())
				return
			}
		}()

		err := filepath.WalkDir(sessionFolder, func(fullpath string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}

			filename := filepath.Base(fullpath)

			if utils.FileLocked(filename) {
				return nil
			}

			re := regexp.MustCompile(`\.[a-z0-9-]+\.json$`)

			if !re.MatchString(filename) {
				return nil
			}

			model := re.ReplaceAllLiteralString(filename, "")

			content, err := os.ReadFile(fullpath)
			if err != nil {
				return err
			}

			if data := map[string]any{}; json.Unmarshal(content, &data) == nil {
				err := HandleNestedModel(data, &model, phoneNumber, cacheFolder, saveFunc, sessions, sessionFolder, nil)
				if err != nil {
					return err
				}
			} else if data := []map[string]any{}; json.Unmarshal(content, &data) == nil {
				for _, row := range data {
					err := HandleNestedModel(row, &model, phoneNumber, cacheFolder, saveFunc, sessions, sessionFolder, nil)
					if err != nil {
						return err
					}
				}
			}

			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}
