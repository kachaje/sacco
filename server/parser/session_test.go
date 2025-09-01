package parser_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sacco/server/parser"
	"testing"
)

func TestLoadKeys(t *testing.T) {
	data := map[string]any{
		"child": map[string]any{
			"id":    "1",
			"field": "value",
			"child1": map[string]any{
				"id":    "2",
				"field": "value",
				"child1_1": map[string]any{
					"id":    "3",
					"field": "value",
					"child1_1_1": map[string]any{
						"id":    "4",
						"field": "value",
					},
					"child1_1_2": []map[string]any{
						{
							"id":    "5",
							"field": "value",
						},
						{
							"id":    "6",
							"field": "value",
						},
					},
					"child1_1_3": map[string]any{
						"id":    "7",
						"field": "value",
						"child1_1_3_1": map[string]any{
							"id":    "8",
							"field": "value",
						},
					},
				},
			},
			"child2": []map[string]any{
				{
					"id":    "9",
					"field": "value",
					"child2_1": map[string]any{
						"id":    "10",
						"field": "value",
						"child2_1_1": []map[string]any{
							{
								"id":    "11",
								"field": "value",
							},
							{
								"id":    "12",
								"field": "value",
							},
						},
					},
				},
				{
					"id":    "13",
					"field": "value",
				},
			},
		},
	}
	target := map[string]any{
		"child1Id":     "2",
		"child1_1Id":   "3",
		"child1_1_1Id": "4",
		"child1_1_2": []map[string]any{
			{
				"child1_1_2Id": "5",
			},
			{
				"child1_1_2Id": "6",
			},
		},
		"child1_1_3Id":   "7",
		"child1_1_3_1Id": "8",
		"child2": []map[string]any{
			{
				"child2Id":   "9",
				"child2_1Id": "10",
				"child2_1_1": []map[string]any{
					{
						"child2_1_1Id": "11",
					},
					{
						"child2_1_1Id": "12",
					},
				},
			},
			{
				"child2Id": "13",
			},
		},
		"childId": "1",
	}

	session := parser.NewSession(nil, nil, nil)

	result := session.LoadKeys(data, map[string]any{}, nil)

	if !reflect.DeepEqual(target, result) {
		t.Fatal("Test failed")
	}
}

func TestUpdateSessionFlags(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("..", "database", "models", "fixtures", "sample.json"))
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]any{}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	session := parser.NewSession(nil, nil, nil)
	session.UpdateActiveData(data, 0)

	err = session.UpdateSessionFlags()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(session.GlobalIds)
}
