package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sacco/forms"
)

func main() {
	bot := forms.NewMembershipChatBot()
	scanner := bufio.NewScanner(os.Stdin)

	var input string

	for {
		fmt.Print("\033[H\033[2J")

		question := bot.ProcessInput(input)

		if question == "" {
			break
		}

		fmt.Println(question)
		scanner.Scan()

		input = scanner.Text()
	}

	payload, _ := json.MarshalIndent(bot.Data, "", "  ")

	fmt.Printf("%s\n", payload)
}
