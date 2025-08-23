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
	BusinessInfoWorkflow  *WorkFlow
	PreferredLanguage     string
	MemberId              *int64
	SessionId             string
	PhoneNumber           string

	AddedModels map[string]bool

	ContactsAdded      bool
	NomineeAdded       bool
	OccupationAdded    bool
	BeneficiariesAdded bool
	BusinessInfoAdded  bool
	ActiveMemberData   map[string]any

	QueryFn    func(string, []string, []string) (map[string]any, error)
	SkipFields []string

	Mu *sync.Mutex
}

func NewSession(queryFn func(string, []string, []string) (map[string]any, error)) *Session {
	return &Session{
		QueryFn:     queryFn,
		Mu:          &sync.Mutex{},
		AddedModels: map[string]bool{},
	}
}

func (s *Session) UpdateSessionFlags() error {
	beneficiariesData := s.ReadFromMap("memberBeneficiary", 0)
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

	contactDetailsData := s.ReadFromMap("memberContact", 0)
	if contactDetailsData != nil {
		val, ok := contactDetailsData.(map[string]any)
		if ok && len(val) > 0 {
			s.ContactsAdded = true
		}
	}

	nomineeDetailsData := s.ReadFromMap("memberNominee", 0)
	if nomineeDetailsData != nil {
		val, ok := nomineeDetailsData.(map[string]any)
		if ok && len(val) > 0 {
			s.NomineeAdded = true
		}
	}

	occupationDetailsData := s.ReadFromMap("memberOccupation", 0)
	if occupationDetailsData != nil {
		val, ok := occupationDetailsData.(map[string]any)
		if ok && len(val) > 0 {
			s.OccupationAdded = true
		}
	}

	idData := s.ReadFromMap("id", 0)
	if idData != nil {
		val := fmt.Sprintf("%v", idData)

		id, err := strconv.ParseInt(val, 10, 64)
		if err == nil {
			s.MemberId = &id
		}
	}

	return nil
}

func (s *Session) UpdateActiveMemberData(data map[string]any, retries int) {
	time.Sleep(time.Duration(retries) * time.Second)

	if s.Mu == nil {
		s.Mu = &sync.Mutex{}
	}

	done := s.Mu.TryLock()
	if !done {
		if retries < 3 {
			retries++
			s.UpdateActiveMemberData(data, retries)
			return
		}
	}
	defer s.Mu.Unlock()

	s.ActiveMemberData = data
}

func (s *Session) WriteToMap(key string, value any, retries int) {
	time.Sleep(time.Duration(retries) * time.Second)

	if s.Mu == nil {
		s.Mu = &sync.Mutex{}
	}

	done := s.Mu.TryLock()
	if !done {
		if retries < 3 {
			retries++
			s.WriteToMap(key, value, retries)
			return
		}
	}
	defer s.Mu.Unlock()

	if s.ActiveMemberData == nil {
		s.ActiveMemberData = map[string]any{}
	}

	s.ActiveMemberData[key] = value
}

func (s *Session) ReadFromMap(key string, retries int) any {
	time.Sleep(time.Duration(retries) * time.Second)

	if s.Mu == nil {
		s.Mu = &sync.Mutex{}
	}

	done := s.Mu.TryLock()
	if !done {
		if retries < 3 {
			retries++
			return s.ReadFromMap(key, retries)
		}
	}
	defer s.Mu.Unlock()

	return s.ActiveMemberData[key]
}

func (s *Session) LoadMemberCache(phoneNumber, cacheFolder string) error {
	sessionFolder := filepath.Join(cacheFolder, phoneNumber)

	_, err := os.Stat(sessionFolder)
	if !os.IsNotExist(err) {
		os.MkdirAll(sessionFolder, 0755)
	}

	for _, key := range []string{"memberContact", "memberNominee", "memberOccupation", "memberBeneficiary"} {
		filename := filepath.Join(sessionFolder, fmt.Sprintf("%s.json", key))

		_, err := os.Stat(filename)
		if os.IsNotExist(err) {
			continue
		}

		content, err := os.ReadFile(filename)
		if err != nil {
			continue
		}

		if key == "memberBeneficiary" {
			data := []map[string]any{}
			err = json.Unmarshal(content, &data)
			if err != nil {
				continue
			}

			s.WriteToMap(key, data, 0)
		} else {
			data := map[string]any{}
			err = json.Unmarshal(content, &data)
			if err != nil {
				continue
			}

			s.WriteToMap(key, data, 0)
		}
	}

	s.UpdateSessionFlags()

	return nil
}

func (s *Session) RefreshSession() (map[string]any, error) {
	if s.PhoneNumber != "" && s.QueryFn != nil {
		data, err := s.QueryFn(s.PhoneNumber, nil, s.SkipFields)
		if err != nil {
			return s.ActiveMemberData, err
		}

		s.UpdateActiveMemberData(data, 0)

		return data, nil
	}
	return s.ActiveMemberData, nil
}
