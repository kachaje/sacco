package main

import (
	"fmt"
	"sacco/chatbot"
)

func main() {
	bot := chatbot.NewChatbot()

	fmt.Println(bot.ProcessInput("")) // Initial greeting

	fmt.Println(bot.ProcessInput("Alice")) // User provides name

	fmt.Println(bot.ProcessInput("alice@example.com")) // User provides email

	fmt.Println(bot.ProcessInput("Hello again")) // After confirmation
}
