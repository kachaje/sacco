package parser_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sacco/parser"
	"sacco/utils"
	"testing"
)

var content []byte
var err error
var data map[string]any

func init() {
	content, err = os.ReadFile(filepath.Join(".", "fixtures", "sample.json"))
	if err != nil {
		panic(err)
	}

	data = map[string]any{}

	err = json.Unmarshal(content, &data)
	if err != nil {
		panic(err)
	}
}

func TestGetNode(t *testing.T) {
	wf := parser.NewWorkflow(data)

	result := wf.GetNode("enterLanguage")

	if result == nil {
		t.Fatal("Test failed")
	}

	for _, key := range []string{"type", "text", "options", "inputIdentifier", "nextScreen"} {
		if result[key] == nil {
			t.Fatalf("Test failed on key %s", key)
		}
	}
}

func TestInputIncluded(t *testing.T) {
	targetRoute := "enterOtherName"

	options := []any{
		map[string]any{
			"position": 1,
			"label": map[string]any{
				"en": "Yes",
				"ny": "Inde",
			},
			"nextScreen": targetRoute,
		},
		map[string]any{
			"position": 2,
			"label": map[string]any{
				"en": "No",
				"ny": "Ayi",
			},
			"nextScreen": "enterGender",
		},
	}

	wf := parser.NewWorkflow(data)

	defaultRoute := "enterAskOtherName"

	wf.CurrentScreen = defaultRoute

	result, nextRoute := wf.InputIncluded("3", options)

	if result {
		t.Fatalf("Test failed. Expected: false; Actual: %v", result)
	}
	if nextRoute != "" {
		t.Fatalf("Test failed. Expected: %s; Actual: %s", defaultRoute, nextRoute)
	}

	result, nextRoute = wf.InputIncluded("1", options)

	if !result {
		t.Fatalf("Test failed. Expected: true; Actual: %v", result)
	}
	if nextRoute != targetRoute {
		t.Fatalf("Test failed. Expected: %s; Actual: %s", targetRoute, nextRoute)
	}

	wf.CurrentScreen = defaultRoute

	result, nextRoute = wf.InputIncluded("2", options)

	if !result {
		t.Fatalf("Test failed. Expected: true; Actual: %v", result)
	}
	if nextRoute != "enterGender" {
		t.Fatalf("Test failed. Expected: enterGender; Actual: %s", nextRoute)
	}
}

func TestNodeOptions(t *testing.T) {
	wf := parser.NewWorkflow(data)

	result := wf.NodeOptions("enterLanguage")

	if len(result) != 2 {
		t.Fatalf("Test failed. Expected: 2; Actual: %v", len(result))
	}

	for i, entry := range []string{"1. English", "2. Chichewa"} {
		if result[i] != entry {
			t.Fatalf("Test failed. Expected: %s; Actual: %s", entry, result[i])
		}
	}

	wf.CurrentLanguage = "ny"

	result = wf.NodeOptions("enterMaritalStatus")

	if len(result) != 4 {
		t.Fatalf("Test failed. Expected: 4; Actual: %v", len(result))
	}

	for i, entry := range []string{"1. Inde", "2. Ayi", "3. Woferedwa", "4. Osudzulidwa"} {
		if result[i] != entry {
			t.Fatalf("Test failed. Expected: %s; Actual: %s", entry, result[i])
		}
	}

	wf.CurrentLanguage = "en"

	result = wf.NodeOptions("enterMaritalStatus")

	if len(result) != 4 {
		t.Fatalf("Test failed. Expected: 4; Actual: %v", len(result))
	}

	for i, entry := range []string{"1. Married", "2. Single", "3. Widowed", "4. Divorced"} {
		if result[i] != entry {
			t.Fatalf("Test failed. Expected: %s; Actual: %s", entry, result[i])
		}
	}
}

func TestNextNode(t *testing.T) {
	wf := parser.NewWorkflow(data)

	result := wf.NextNode("")

	if result == nil {
		t.Fatal("Test failed")
	}

	for _, key := range []string{"type", "text", "options", "inputIdentifier", "nextScreen"} {
		if result[key] == nil {
			t.Fatalf("Test failed on key %s", key)
		}
	}

	if wf.CurrentScreen != "enterLanguage" {
		t.Fatalf("Test failed. Expected: 'enterLanguage'; Actual: '%v'", wf.CurrentScreen)
	}

	if wf.PreviousScreen != "initialScreen" {
		t.Fatalf("Test failed. Expected: 'initialScreen'; Actual: '%v'", wf.PreviousScreen)
	}

	result = wf.NextNode("3")

	for _, key := range []string{"type", "text", "options", "inputIdentifier", "nextScreen"} {
		if result[key] == nil {
			t.Fatalf("Test failed on key %s", key)
		}
	}

	if wf.CurrentScreen != "enterLanguage" {
		t.Fatalf("Test failed. Expected: 'enterLanguage'; Actual: '%v'", wf.CurrentScreen)
	}

	if wf.PreviousScreen != "initialScreen" {
		t.Fatalf("Test failed. Expected: 'initialScreen'; Actual: '%v'", wf.PreviousScreen)
	}

	wf.CurrentScreen = "enterDateOfBirth"

	wf.NextNode("1999")

	if wf.CurrentScreen != "enterDateOfBirth" {
		t.Fatalf("Test failed. Expected: 'enterDateOfBirth'; Actual: '%v'", wf.CurrentScreen)
	}

	wf.NextNode("1999-09-01")

	if wf.CurrentScreen != "enterMaritalStatus" {
		t.Fatalf("Test failed. Expected: 'enterMaritalStatus'; Actual: '%v'", wf.CurrentScreen)
	}

	wf.CurrentScreen = "enterLanguage"

	result = wf.NextNode("1")

	for _, key := range []string{"type", "text", "inputIdentifier", "nextScreen"} {
		if result[key] == nil {
			t.Fatalf("Test failed on key %s", key)
		}
	}

	if wf.CurrentScreen != "enterFirstName" {
		t.Fatalf("Test failed. Expected: 'enterFirstName'; Actual: '%v'", wf.CurrentScreen)
	}

	if wf.PreviousScreen != "enterLanguage" {
		t.Fatalf("Test failed. Expected: 'enterLanguage'; Actual: '%v'", wf.PreviousScreen)
	}

	if wf.Data["dateOfBirth"] == nil || fmt.Sprintf("%v", wf.Data["dateOfBirth"]) != "1999-09-01" {
		t.Fatalf("Test failed. Expected: '1999-09-01'; Actual: %v", wf.Data["dateOfBirth"])
	}

	if wf.Data["language"] == nil || fmt.Sprintf("%v", wf.Data["language"]) != "1" {
		t.Fatalf("Test failed. Expected: '1'; Actual: %v", wf.Data["dateOfBirth"])
	}
}

func TestOptionValue(t *testing.T) {
	wf := parser.NewWorkflow(data)

	wf.CurrentLanguage = "2"

	options := []any{
		map[string]any{
			"position": 1,
			"label": map[string]any{
				"en": "Yes",
				"ny": "Inde",
			},
			"nextScreen": "",
		},
		map[string]any{
			"position": 2,
			"label": map[string]any{
				"all": "No",
			},
			"nextScreen": "enterGender",
		},
	}

	result := wf.OptionValue(options, "2")

	if result != "No" {
		t.Fatalf("Test failed. Expected: No; Actual: %v", result)
	}

	result = wf.OptionValue(options, "1")

	if result != "Yes" {
		t.Fatalf("Test failed. Expected: Yes; Actual: %v", result)
	}
}

func TestResolveData(t *testing.T) {
	wf := parser.NewWorkflow(data)

	result := wf.ResolveData(map[string]any{
		"language":      "1",
		"firstName":     "Mary",
		"lastName":      "Banda",
		"askOtherName":  "2",
		"dateOfBirth":   "1999-09-01",
		"maritalStatus": "2",
	})

	target := map[string]any{
		"language":      "English",
		"firstName":     "Mary",
		"lastName":      "Banda",
		"askOtherName":  "No",
		"dateOfBirth":   "1999-09-01",
		"maritalStatus": "Single",
	}

	fmt.Println(result)

	if result == nil {
		t.Fatal("Test failed")
	}

	for key, val := range target {
		if result[key] == nil || fmt.Sprintf("%v", result[key]) != fmt.Sprintf("%v", val) {
			t.Fatalf("Test failed. Expected: %v; Actual: %v", val, result[key])
		}
	}
}

func TestGetLabel(t *testing.T) {
	wf := parser.NewWorkflow(data)

	node := wf.NextNode("")

	if node == nil {
		t.Fatal("Test failed")
	}

	target := "Language: 1. English 2. Chichewa 99. Cancel"

	result := wf.GetLabel(node, wf.CurrentScreen)

	label := utils.CleanScript([]byte(result))

	if label != target {
		t.Fatalf("Test failed. Expected: %s; Actual: %s", target, label)
	}

	wf.CurrentLanguage = "2"

	target = "Language: 1. English 2. Chichewa 99. Basi"

	result = wf.GetLabel(node, wf.CurrentScreen)

	label = utils.CleanScript([]byte(result))

	if label != target {
		t.Fatalf("Test failed. Expected: %s; Actual: %s", target, label)
	}

	wf.Data = map[string]any{
		"language":      "1",
		"firstName":     "Mary",
		"lastName":      "Banda",
		"askOtherName":  "2",
		"dateOfBirth":   "1999-09-01",
		"maritalStatus": "2",
	}
}
