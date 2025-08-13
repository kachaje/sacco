package utils

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

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
