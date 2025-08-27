package utils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
)

func CleanScript(content []byte) string {
	stage1 := regexp.MustCompile(`\n|\r`).ReplaceAllLiteralString(string(content), " ")

	return strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllLiteralString(stage1, " "))
}

func CleanString(content string) string {
	stage1 := regexp.MustCompile(`\n|\r`).ReplaceAllLiteralString(string(content), " ")

	return strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllLiteralString(stage1, " "))
}

func DumpYaml(data map[string]any) (*string, error) {
	var result string

	payload, err := yaml.Marshal(data)
	if err != nil {
		return nil, err
	}

	result = string(payload)

	return &result, nil
}

func LoadYaml(yamlData string) (map[string]any, error) {
	var data map[string]any

	err := yaml.Unmarshal([]byte(yamlData), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func WaitForPort(host string, port string, timeout time.Duration, retryInterval time.Duration, debug bool) error {
	address := net.JoinHostPort(host, port)
	startTime := time.Now()

	for {
		conn, err := net.DialTimeout("tcp", address, retryInterval)
		if err == nil {
			conn.Close()
			if debug {
				fmt.Printf("Port %s on %s is open.\n", port, host)
			}
			return nil
		}

		if time.Since(startTime) >= timeout {
			return fmt.Errorf("timeout waiting for port %s on %s: %w", port, host, err)
		}

		if debug {
			fmt.Printf("Waiting for port %s on %s... Retrying in %v\n", port, host, retryInterval)
		}

		time.Sleep(retryInterval)
	}
}

func GetFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "0.0.0.0:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}

func LockFile(filename string) (string, error) {
	lockFilename := fmt.Sprintf("%s.lock", filename)

	return lockFilename, os.WriteFile(lockFilename, []byte{}, 0644)
}

func UnLockFile(filename string) error {
	lockFilename := fmt.Sprintf("%s.lock", filename)

	return os.Remove(lockFilename)
}

func FileLocked(filename string) bool {
	lockFilename := fmt.Sprintf("%s.lock", filename)

	_, err := os.Stat(lockFilename)

	return !os.IsNotExist(err)
}

func QueryWithRetry(db *sql.DB, ctx context.Context, retries int, query string, args ...any) (sql.Result, error) {
	time.Sleep(time.Duration(retries) * time.Second)

	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		if regexp.MustCompile("SQL logic error: no such table").MatchString(err.Error()) {
			if retries < 3 {
				retries++

				log.Printf("models.QueryWithRetry.retry: %d\n", retries)

				return QueryWithRetry(db, ctx, retries, query, args...)
			}
		}
		return nil, fmt.Errorf("models.QueryWithRetry.1: %s", err.Error())
	}

	return result, nil
}

func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func IdentifierToLabel(identifier string) string {
	re := regexp.MustCompile("([A-Z][a-z]*)")

	parts := re.FindAllString(CapitalizeFirstLetter(identifier), -1)

	return strings.Join(parts, " ")
}

func Index[T comparable](s []T, item T) int {
	for i, v := range s {
		if v == item {
			return i
		}
	}
	return -1
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
