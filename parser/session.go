package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type Session struct {
	CurrentMenu           string
	Data                  map[string]string
	PIWorkflow            *WorkFlow
	LanguageWorkflow      *WorkFlow
	OccupationWorkflow    *WorkFlow
	ContactsWorkflow      *WorkFlow
	NomineeWorkflow       *WorkFlow
	BeneficiariesWorkflow *WorkFlow
	PreferredLanguage     string
	MemberId              *int64
	SessionId             string
	PhoneNumber           string

	ContactsAdded      bool
	NomineeAdded       bool
	OccupationAdded    bool
	BeneficiariesAdded bool
	ActiveMemberData   map[string]any

	QueryFn func(string) (map[string]any, error)

	Mu *sync.Mutex
}

func NewSession(queryFn func(string) (map[string]any, error)) *Session {
	return &Session{
		QueryFn: queryFn,
		Mu:      &sync.Mutex{},
	}
}

func (s *Session) UpdateSessionFlags() error {
	beneficiariesData := s.ReadFromMap("beneficiaries")
	if beneficiariesData != nil {
		val, ok := beneficiariesData.([]any)
		if ok && len(val) > 0 {
			s.BeneficiariesAdded = true
		} else {
			val, ok := beneficiariesData.([]map[string]any)
			if ok && len(val) > 0 {
				s.BeneficiariesAdded = true
			}
		}
	}

	contactDetailsData := s.ReadFromMap("contactDetails")
	if contactDetailsData != nil {
		val, ok := contactDetailsData.(map[string]any)
		if ok && len(val) > 0 {
			s.ContactsAdded = true
		}
	}

	nomineeDetailsData := s.ReadFromMap("nomineeDetails")
	if nomineeDetailsData != nil {
		val, ok := nomineeDetailsData.(map[string]any)
		if ok && len(val) > 0 {
			s.NomineeAdded = true
		}
	}

	occupationDetailsData := s.ReadFromMap("occupationDetails")
	if occupationDetailsData != nil {
		val, ok := occupationDetailsData.(map[string]any)
		if ok && len(val) > 0 {
			s.OccupationAdded = true
		}
	}

	idData := s.ReadFromMap("id")
	if idData != nil {
		val := fmt.Sprintf("%v", idData)

		id, err := strconv.ParseInt(val, 10, 64)
		if err == nil {
			s.MemberId = &id
		}
	}

	return nil
}

func (s *Session) WriteToMap(key string, value any) {
	retries := 0

RETRY:
	time.Sleep(time.Duration(retries) * time.Second)

	if s.Mu == nil {
		s.Mu = &sync.Mutex{}
	}

	done := s.Mu.TryLock()
	if !done {
		if retries < 3 {
			retries++
			goto RETRY
		}
	}
	defer s.Mu.Unlock()

	if s.ActiveMemberData == nil {
		s.ActiveMemberData = map[string]any{}
	}

	s.ActiveMemberData[key] = value
}

func (s *Session) ReadFromMap(key string) any {
	retries := 0

RETRY:
	time.Sleep(time.Duration(retries) * time.Second)

	if s.Mu == nil {
		s.Mu = &sync.Mutex{}
	}

	done := s.Mu.TryLock()
	if !done {
		if retries < 3 {
			retries++
			goto RETRY
		}
	}
	defer s.Mu.Unlock()

	return s.ActiveMemberData[key]
}

func (s *Session) LoadMemberCache(phoneNumber, cacheFolder string) error {
	sessionFolder := filepath.Join(cacheFolder, phoneNumber)

	_, err := os.Stat(sessionFolder)
	if os.IsNotExist(err) {
		return err
	}

	for _, key := range []string{"contactDetails", "nomineeDetails", "occupationDetails", "beneficiaries"} {
		filename := filepath.Join(sessionFolder, fmt.Sprintf("%s.json", key))

		_, err := os.Stat(filename)
		if os.IsNotExist(err) {
			continue
		}

		content, err := os.ReadFile(filename)
		if err != nil {
			continue
		}

		if key == "beneficiaries" {
			data := []map[string]any{}
			err = json.Unmarshal(content, &data)
			if err != nil {
				continue
			}

			s.WriteToMap(key, data)
		} else {
			data := map[string]any{}
			err = json.Unmarshal(content, &data)
			if err != nil {
				continue
			}

			s.WriteToMap(key, data)
		}
	}

	s.UpdateSessionFlags()

	return nil
}

func (s *Session) RefreshSession() (map[string]any, error) {
	if s.PhoneNumber != "" && s.QueryFn != nil {
		data, err := s.QueryFn(s.PhoneNumber)
		if err != nil {
			return s.ActiveMemberData, err
		}

		s.ActiveMemberData = data

		return data, nil
	}
	return s.ActiveMemberData, nil
}
