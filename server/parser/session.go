package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sacco/server/database"
	"slices"
	"strconv"
	"sync"
	"time"
)

type Session struct {
	CurrentMenu string
	Data        map[string]string

	PreferredLanguage string
	MemberId          *int64
	SessionId         string
	PhoneNumber       string

	WorkflowsMapping map[string]*WorkFlow

	AddedModels map[string]bool

	ActiveData map[string]any

	QueryFn    func(string, []string, []string) (map[string]any, error)
	SkipFields []string

	Mu *sync.Mutex

	SessionToken *string
	SessionUser  *string

	Cache      map[string]string
	LastPrompt string
}

func NewSession(
	queryFn func(string, []string, []string) (map[string]any, error),
	phoneNumber, sessionId *string,
) *Session {
	s := &Session{
		QueryFn:          queryFn,
		Mu:               &sync.Mutex{},
		AddedModels:      map[string]bool{},
		ActiveData:       map[string]any{},
		Data:             map[string]string{},
		SkipFields:       []string{"active"},
		CurrentMenu:      "main",
		WorkflowsMapping: map[string]*WorkFlow{},
		Cache:            map[string]string{},
		LastPrompt:       "",
	}

	if phoneNumber != nil {
		s.PhoneNumber = *phoneNumber
	}
	if sessionId != nil {
		s.SessionId = *sessionId
	}

	return s
}

func (s *Session) UpdateSessionFlags() error {
	for _, model := range database.MemberSingleChildren {
		data := s.ReadFromMap(model, 0)
		if data != nil {
			val, ok := data.(map[string]any)
			if ok && len(val) > 0 {
				s.AddedModels[model] = true
			}
		}
	}

	for _, model := range database.MemberArrayChildren {
		data := s.ReadFromMap(model, 0)
		if data != nil {
			if val, ok := data.([]any); ok && len(val) > 0 {
				s.AddedModels[model] = true
			} else if val, ok := data.([]map[string]any); ok && len(val) > 0 {
				s.AddedModels[model] = true
			}
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

	nameData := s.ReadFromMap("firstName", 0)
	if nameData != nil {
		s.AddedModels["member"] = true
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

	s.ActiveData = data
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

	if s.ActiveData == nil {
		s.ActiveData = map[string]any{}
	}

	s.ActiveData[key] = value
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

	return s.ActiveData[key]
}

func (s *Session) LoadMemberCache(phoneNumber, cacheFolder string) error {
	sessionFolder := filepath.Join(cacheFolder, phoneNumber)

	_, err := os.Stat(sessionFolder)
	if !os.IsNotExist(err) {
		os.MkdirAll(sessionFolder, 0755)
	}

	models := []string{}

	models = append(models, database.MemberArrayChildren...)

	models = append(models, database.MemberSingleChildren...)

	for _, key := range models {
		filename := filepath.Join(sessionFolder, fmt.Sprintf("%s.json", key))

		_, err := os.Stat(filename)
		if os.IsNotExist(err) {
			continue
		}

		content, err := os.ReadFile(filename)
		if err != nil {
			continue
		}

		if slices.Contains(database.MemberArrayChildren, key) {
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
			return s.ActiveData, err
		}

		s.UpdateActiveMemberData(data, 0)

		return data, nil
	}
	return s.ActiveData, nil
}
