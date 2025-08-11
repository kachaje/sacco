package chatbot

import (
	"fmt"
	"strings"
)

type State interface {
	Handle(input string, data *ChatData) (State, string)
}

type ChatData struct {
	UserName  string
	UserEmail string
}

type WelcomeState struct{}

func (s *WelcomeState) Handle(input string, data *ChatData) (State, string) {
	return &AskingForNameState{}, "Hello! What's your name?"
}

type AskingForNameState struct{}

func (s *AskingForNameState) Handle(input string, data *ChatData) (State, string) {
	data.UserName = strings.TrimSpace(input)
	return &AskingForEmailState{}, fmt.Sprintf("Nice to meet you, %s! What's your email?", data.UserName)
}

type AskingForEmailState struct{}

func (s *AskingForEmailState) Handle(input string, data *ChatData) (State, string) {
	if !strings.Contains(input, "@") {
		return s, "That doesn't look like a valid email. Please try again."
	}
	data.UserEmail = strings.TrimSpace(input)
	return &ConfirmedState{}, "Thanks! Your information has been recorded."
}

type ConfirmedState struct{}

func (s *ConfirmedState) Handle(input string, data *ChatData) (State, string) {
	return s, "We've already confirmed your details. Is there anything else I can help with?"
}

type Chatbot struct {
	CurrentState State
	Data         *ChatData
}

func NewChatbot() *Chatbot {
	return &Chatbot{
		CurrentState: &WelcomeState{},
		Data:         &ChatData{},
	}
}

func (cb *Chatbot) ProcessInput(input string) string {
	var response string
	cb.CurrentState, response = cb.CurrentState.Handle(input, cb.Data)
	return response
}
