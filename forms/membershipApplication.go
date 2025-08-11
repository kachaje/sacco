package forms

import (
	"regexp"
	"slices"
	"strings"
)

type State interface {
	Handle(input string, data *MemberData) (State, string)
}

type MemberData struct {
	Surname         string
	FirstName       string
	OtherNames      string
	Gender          string
	Title           string
	MaritalStatus   string
	DateOfBirth     string
	NationalID      string
	UtilityBillNo   string
	UtilityBillType string
}

type StartState struct{}
type SurnameState struct{}
type FirstNameState struct{}
type OtherNamesState struct{}
type GenderState struct{}
type TitleState struct{}
type MaritalStatusState struct{}
type DateOfBirthState struct{}
type NationalIDState struct{}
type UtilityBillNoState struct{}
type UtilityBillTypeState struct{}
type DoneState struct{}

func (s *StartState) Handle(input string, data *MemberData) (State, string) {
	return &FirstNameState{}, "First Name:"
}

func (s *FirstNameState) Handle(input string, data *MemberData) (State, string) {
	if input == "" {
		return &FirstNameState{}, "First Name:"
	} else {
		data.Surname = input

		return &SurnameState{}, "Surname:"
	}
}

func (s *SurnameState) Handle(input string, data *MemberData) (State, string) {
	if input == "" {
		return &SurnameState{}, "Surname:"
	} else {
		data.FirstName = input

		return &OtherNamesState{}, "Other Name(s):"
	}
}

func (s *OtherNamesState) Handle(input string, data *MemberData) (State, string) {
	if input != "" {
		data.OtherNames = input
	}

	return &GenderState{}, `Gender:
1. Female
2. Male`
}

func (s *GenderState) Handle(input string, data *MemberData) (State, string) {
	if !slices.Contains([]string{"1", "2"}, input) {
		return &MaritalStatusState{}, `Gender:
1. Female
2. Male`
	} else {
		data.Gender = map[string]string{
			"1": "Female",
			"2": "Male",
		}[input]

		return &TitleState{}, `Title:
1. Mr
2. Mrs
3. Miss
4. Dr
5. Prof
6. Rev
7. Other`
	}
}

func (s *TitleState) Handle(input string, data *MemberData) (State, string) {
	if !slices.Contains([]string{"1", "2", "3", "4", "5", "6", "7"}, input) {
		data.MaritalStatus = map[string]string{
			"1": "Mr",
			"2": "Mrs",
			"3": "Miss",
			"4": "Dr",
			"5": "Prof",
			"6": "Rev",
			"7": "Other",
		}[input]
	}

	return &MaritalStatusState{}, `Marital Status:
1. Married
2. Single
3. Widowed
4. Divorced`
}

func (s *MaritalStatusState) Handle(input string, data *MemberData) (State, string) {
	if !slices.Contains([]string{"1", "2", "3", "4"}, input) {
		return &MaritalStatusState{}, `Marital Status:
1. Married
2. Single
3. Widowed
4. Divorced`
	} else {
		data.Title = map[string]string{
			"1": "Married",
			"2": "Single",
			"3": "Widowed",
			"4": "Divorced",
		}[input]

		return &DateOfBirthState{}, "Date of Birth (YYYY-MM-DD):"
	}
}

func (s *DateOfBirthState) Handle(input string, data *MemberData) (State, string) {
	if !regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`).MatchString(input) {
		return &DateOfBirthState{}, "Date of Birth (YYYY-MM-DD):"
	} else {
		data.DateOfBirth = input

		return &NationalIDState{}, "National ID:"
	}
}

func (s *NationalIDState) Handle(input string, data *MemberData) (State, string) {
	if input == "" {
		return &NationalIDState{}, "National ID:"
	} else {
		data.NationalID = input

		return &UtilityBillNoState{}, "Utility Bill Number:"
	}
}

func (s *UtilityBillNoState) Handle(input string, data *MemberData) (State, string) {
	if input == "" {
		return &UtilityBillNoState{}, "Utility Bill Number:"
	} else {
		data.UtilityBillNo = input

		return &UtilityBillTypeState{}, `Utility Bill Type:
1. Escom
2. Water Board`
	}
}

func (s *UtilityBillTypeState) Handle(input string, data *MemberData) (State, string) {
	if !slices.Contains([]string{"1", "2"}, input) {
		return &UtilityBillTypeState{}, `Utility Bill Type:
1. Escom
2. Water Board`
	} else {
		data.UtilityBillType = map[string]string{
			"1": "Escom",
			"2": "Water Board",
		}[input]

		return &DoneState{}, ""
	}
}

func (s *DoneState) Handle(input string, data *MemberData) (State, string) {
	return nil, ""
}

type MembershipChatbot struct {
	CurrentState State
	Data         *MemberData
}

func NewMembershipChatBot() *MembershipChatbot {
	return &MembershipChatbot{
		CurrentState: &StartState{},
		Data:         &MemberData{},
	}
}

func (cb *MembershipChatbot) ProcessInput(input string) string {
	var response string

	input = strings.TrimSpace(input)

	cb.CurrentState, response = cb.CurrentState.Handle(input, cb.Data)

	if cb.CurrentState != nil {
		return response
	} else {
		return ""
	}
}
