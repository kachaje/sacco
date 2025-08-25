package menus_test

import (
	"fmt"
	"sacco/server/menus"
	"testing"
)

func TestLoadMenus(t *testing.T) {
	m := menus.NewMenus()

	result := m.LoadMenu("main", nil, "", "", "", "")

	fmt.Println(result)
}
