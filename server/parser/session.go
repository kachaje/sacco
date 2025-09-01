package parser

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

type Session struct {
	CurrentMenu string
	Data        map[string]string

	PreferredLanguage  string
	SessionId          string
	PhoneNumber        string
	CurrentPhoneNumber string

	GlobalIds map[string]int64

	WorkflowsMapping map[string]*WorkFlow

	AddedModels map[string]bool

	ActiveData map[string]any

	QueryFn    func(string, []string) (map[string]any, error)
	SkipFields []string

	Mu *sync.Mutex

	SessionToken    *string
	SessionUser     *string
	SessionUserRole *string
	SessionUserId   *int64

	Cache      map[string]string
	LastPrompt string
}

func NewSession(
	queryFn func(string, []string) (map[string]any, error),
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
		s.CurrentPhoneNumber = *phoneNumber
	}
	if sessionId != nil {
		s.SessionId = *sessionId
	}

	return s
}

func (s *Session) LoadKeys(rawData any, seed map[string]any, parent *string) map[string]any {
	if seed == nil {
		seed = map[string]any{}
	}

	if data, ok := rawData.(map[string]any); ok {
		for key, value := range data {
			if key == "id" {
				if parent != nil {
					seed[fmt.Sprintf("%vId", *parent)] = fmt.Sprintf("%v", value)
				}
			} else if reflect.TypeOf(value).String() == "map[string]interface {}" && value != nil {
				if val, ok := value.(map[string]any); ok {
					for k, v := range val {
						if k == "id" {
							seed[fmt.Sprintf("%vId", key)] = fmt.Sprintf("%v", v)
						} else {
							seed = s.LoadKeys(v, seed, &k)
						}
					}
				}
			}
		}
	} else if reflect.TypeOf(rawData).String() == "[]map[string]interface {}" && rawData != nil {
		if val, ok := rawData.([]map[string]any); ok {
			if parent != nil {
				seed[*parent] = []map[string]any{}

				for _, row := range val {
					result := s.LoadKeys(row, map[string]any{}, parent)

					if len(result) > 0 {
						seed[*parent] = append(seed[*parent].([]map[string]any), result)
					}
				}
			}
		}
	}

	return seed
}

func (s *Session) UpdateSessionFlags() error {
	data := s.LoadKeys(s.ActiveData, map[string]any{}, nil)

	fmt.Println(data)

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

func (s *Session) ClearSession() {
	s.ActiveData = map[string]any{}
	s.Data = map[string]string{}
	s.AddedModels = map[string]bool{}
	s.GlobalIds = map[string]int64{}
}

func (s *Session) RefreshSession() (map[string]any, error) {
	if s.CurrentPhoneNumber != "" && s.QueryFn != nil {
		data, err := s.QueryFn(s.CurrentPhoneNumber, s.SkipFields)
		if err != nil {
			return nil, err
		}

		s.UpdateActiveData(data, 0)

		return data, nil
	}
	return s.ActiveData, nil
}
