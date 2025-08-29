package parser

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sacco/server/database"
	"slices"
	"sync"
	"time"
)

type Session struct {
	CurrentMenu string
	Data        map[string]string

	PreferredLanguage string
	SessionId         string
	PhoneNumber       string

	GlobalIds map[string]int64

	WorkflowsMapping map[string]*WorkFlow

	AddedModels map[string]bool

	ActiveData map[string]any

	QueryFn    func(string, []string, []string) (map[string]any, error)
	SkipFields []string

	Mu *sync.Mutex

	SessionToken  *string
	SessionUser   *string
	SessionUserId *int64

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
		GlobalIds:        map[string]int64{},
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
	for _, group := range database.SingleChildren {
		for _, model := range group {
			data := s.ReadFromMap(model, 0)
			if data != nil {
				if val, ok := data.(map[string]any); ok && len(val) > 0 {
					s.AddedModels[model] = true
				}
			}
		}
	}

	for _, group := range database.ArrayChildren {
		for _, model := range group {
			data := s.ReadFromMap(model, 0)
			if data != nil {
				if val, ok := data.([]any); ok && len(val) > 0 {
					s.AddedModels[model] = true
				} else if val, ok := data.([]map[string]any); ok && len(val) > 0 {
					s.AddedModels[model] = true
				} else if val, ok := data.(map[string]any); ok && len(val) > 0 {
					s.AddedModels[model] = true
				}
			}
		}
	}

	nameData := s.ReadFromMap("firstName", 0)
	if nameData != nil {
		s.AddedModels["member"] = true
	}

	return nil
}

func (s *Session) UpdateActiveData(data map[string]any, retries int) {
	time.Sleep(time.Duration(retries) * time.Second)

	if s.Mu == nil {
		s.Mu = &sync.Mutex{}
	}

	done := s.Mu.TryLock()
	if !done {
		if retries < 3 {
			retries++
			s.UpdateActiveData(data, retries)
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

func (s *Session) LoadCacheData(phoneNumber, cacheFolder string) error {
	sessionFolder := filepath.Join(cacheFolder, phoneNumber)

	_, err := os.Stat(sessionFolder)
	if !os.IsNotExist(err) {
		os.MkdirAll(sessionFolder, 0755)
	}

	arraysModels := []string{}

	for _, group := range database.ArrayChildren {
		arraysModels = append(arraysModels, group...)
	}

	err = filepath.WalkDir(sessionFolder, func(fullpath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		filename := filepath.Base(fullpath)

		re := regexp.MustCompile(`\.[a-z0-9-]+\.json$`)

		if !re.MatchString(filename) {
			return nil
		}

		model := re.ReplaceAllLiteralString(filename, "")

		content, err := os.ReadFile(fullpath)
		if err != nil {
			return err
		}

		if slices.Contains(arraysModels, model) {
			data := []map[string]any{}
			err = json.Unmarshal(content, &data)
			if err != nil {
				data := map[string]any{}
				err = json.Unmarshal(content, &data)
				if err != nil {
					return err
				}
				s.WriteToMap(model, data, 0)
			} else {
				s.WriteToMap(model, data, 0)
			}
		} else {
			data := map[string]any{}
			err = json.Unmarshal(content, &data)
			if err != nil {
				return err
			}

			s.WriteToMap(model, data, 0)
		}

		return nil
	})
	if err != nil {
		return err
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

		s.UpdateActiveData(data, 0)

		return data, nil
	}
	return s.ActiveData, nil
}
