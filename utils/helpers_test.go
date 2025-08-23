package utils_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sacco/utils"
	"testing"
)

func TestLoadYaml(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(".", "fixtures", "newMember.yml"))
	if err != nil {
		t.Fatal(err)
	}

	result, err := utils.LoadYaml(string(content))
	if err != nil {
		t.Fatal(err)
	}

	target := map[string]any{}

	refData, err := os.ReadFile(filepath.Join(".", "fixtures", "newMember.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(refData, &target)
	if err != nil {
		t.Fatal(err)
	}

	compareObjects := func(obj1, obj2 map[string]any) bool {
		if len(obj1) != len(obj2) {
			return false
		}

		for key, val1 := range obj1 {
			val2, exists := obj2[key]
			if !exists || fmt.Sprintf("%v", val1) != fmt.Sprintf("%v", val2) {
				return false
			}
		}

		return true
	}

	if !compareObjects(target, result) {
		t.Fatal("Test failed")
	}
}

func TestLockFile(t *testing.T) {
	rootFolder := filepath.Join(".", "tmpFileLock")

	os.MkdirAll(rootFolder, 0755)
	defer func() {
		os.RemoveAll(rootFolder)
	}()

	filename := filepath.Join(rootFolder, "lock.txt")

	err := os.WriteFile(filename, []byte{}, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Remove(filename)
	}()

	lockFilename, err := utils.LockFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Remove(lockFilename)
	}()

	_, err = os.Stat(lockFilename)
	if os.IsNotExist(err) {
		t.Fatal("Test failed")
	}
}

func TestUnLockFile(t *testing.T) {
	rootFolder := filepath.Join(".", "tmpFileUnLock")

	os.MkdirAll(rootFolder, 0755)
	defer func() {
		os.RemoveAll(rootFolder)
	}()

	filename := filepath.Join(rootFolder, "lock.txt")
	lockFilename := fmt.Sprintf("%s.lock", filename)

	err := os.WriteFile(lockFilename, []byte{}, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Remove(lockFilename)
	}()

	err = utils.UnLockFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(lockFilename)
	if !os.IsNotExist(err) {
		t.Fatal("Test failed")
	}
}

func TestFileLocked(t *testing.T) {
	rootFolder := filepath.Join(".", "tmpFileLocked")

	os.MkdirAll(rootFolder, 0755)
	defer func() {
		os.RemoveAll(rootFolder)
	}()

	filename := filepath.Join(rootFolder, "lock.txt")
	lockFilename := fmt.Sprintf("%s.lock", filename)

	err := os.WriteFile(lockFilename, []byte{}, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := os.Stat(lockFilename); !os.IsNotExist(err) {
			os.Remove(lockFilename)
		}
	}()

	locked := utils.FileLocked(filename)
	if !locked {
		t.Fatalf("Test failed. Expected: true; Actual: %v", locked)
	}

	err = os.Remove(lockFilename)
	if err != nil {
		t.Fatal(err)
	}

	locked = utils.FileLocked(filename)
	if locked {
		t.Fatalf("Test failed. Expected: false; Actual: %v", locked)
	}
}

func TestIdentifierToLabel(t *testing.T) {
	result := utils.IdentifierToLabel("thisIsAString")

	target := "This Is A String"

	if result != target {
		t.Fatalf("Test failed. Expected: %s; Actual: %s", target, result)
	}
}
