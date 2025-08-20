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
	if s.ActiveMemberData != nil {
		if s.ActiveMemberData["beneficiaries"] != nil {
			val, ok := s.ActiveMemberData["beneficiaries"].([]any)
			if ok && len(val) > 0 {
				s.BeneficiariesAdded = true
			} else {
				val, ok := s.ActiveMemberData["beneficiaries"].([]map[string]any)
				if ok && len(val) > 0 {
					s.BeneficiariesAdded = true
				}
			}
		}
		if s.ActiveMemberData["contactDetails"] != nil {
			val, ok := s.ActiveMemberData["contactDetails"].(map[string]any)
			if ok && len(val) > 0 {
				s.ContactsAdded = true
			}
		}
		if s.ActiveMemberData["nomineeDetails"] != nil {
			val, ok := s.ActiveMemberData["nomineeDetails"].(map[string]any)
			if ok && len(val) > 0 {
				s.NomineeAdded = true
			}
		}
		if s.ActiveMemberData["occupationDetails"] != nil {
			val, ok := s.ActiveMemberData["occupationDetails"].(map[string]any)
			if ok && len(val) > 0 {
				s.OccupationAdded = true
			}
		}
		if s.ActiveMemberData["id"] != nil {
			val := fmt.Sprintf("%v", s.ActiveMemberData["id"])

			id, err := strconv.ParseInt(val, 10, 64)
			if err == nil {
				s.MemberId = &id
			}
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
	s.Mu.TryLock()
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
