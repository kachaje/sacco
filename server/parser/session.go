package parser

import (
	"fmt"
	"reflect"
	"slices"
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

	GlobalIds map[string]any

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
		GlobalIds:        map[string]any{},
	}

	if phoneNumber != nil {
		s.CurrentPhoneNumber = *phoneNumber
	}
	if sessionId != nil {
		s.SessionId = *sessionId
	}

	return s
}

func (s *Session) FlattenKeys(rawData any, seed map[string]any, parent *string) map[string]any {
	handleArrayValues := func(value any, seed map[string]any, parent *string) map[string]any {
		rows := []map[string]any{}

		if val, ok := value.([]map[string]any); ok {
			rows = val
		} else if val, ok := value.([]any); ok {
			for _, row := range val {
				if v, ok := row.(map[string]any); ok {
					rows = append(rows, v)
				}
			}
		}

		if parent != nil {
			for i, row := range rows {
				refKey := fmt.Sprintf("%s.%d", *parent, i)

				seed = s.FlattenKeys(row, seed, &refKey)
			}
		}

		return seed
	}

	var refKey string

	if data, ok := rawData.(map[string]any); ok {
		for key, value := range data {
			if value == nil {
				continue
			}

			if reflect.TypeOf(value).String() == "map[string]interface {}" {
				if val, ok := value.(map[string]any); ok {
					for k, v := range val {
						if v == nil {
							continue
						}

						if slices.Contains([]string{"[]map[string]interface {}", "[]interface {}", "map[string]interface {}"}, reflect.TypeOf(v).String()) {
							if parent != nil {
								refKey = fmt.Sprintf("%s.%s.%s", *parent, key, k)
							} else {
								refKey = fmt.Sprintf("%s.%s", key, k)
							}

							s.FlattenKeys(v, seed, &refKey)
						} else {
							if k == "id" {
								refKey = fmt.Sprintf("%sId", key)
							} else {
								refKey = fmt.Sprintf("%s.%s", key, k)
							}

							seed[refKey] = v
						}
					}
				}
			} else if slices.Contains([]string{"[]map[string]interface {}", "[]interface {}"}, reflect.TypeOf(value).String()) {
				if parent != nil {
					refKey = fmt.Sprintf("%s.%s", *parent, key)
				} else {
					refKey = key
				}

				handleArrayValues(value, seed, &refKey)
			} else {
				if parent != nil {
					refKey = fmt.Sprintf("%s.%s", *parent, key)
				} else {
					refKey = key
				}

				seed[refKey] = value
			}
		}
	} else if slices.Contains([]string{"[]map[string]interface {}", "[]interface {}"}, reflect.TypeOf(rawData).String()) && rawData != nil {
		handleArrayValues(rawData, seed, parent)
	}

	return seed
}

func (s *Session) LoadKeys(rawData any, seed map[string]any, parent *string) map[string]any {
	if seed == nil {
		seed = map[string]any{}
	}

	handleArrayValues := func(value any, seed map[string]any, parent *string) map[string]any {
		rows := []map[string]any{}

		if val, ok := value.([]map[string]any); ok {
			rows = val
		} else if val, ok := value.([]any); ok {
			for _, row := range val {
				if v, ok := row.(map[string]any); ok {
					rows = append(rows, v)
				}
			}
		}

		if parent != nil {
			seed[*parent] = []map[string]any{}

			for _, row := range rows {
				result := s.LoadKeys(row, map[string]any{}, parent)

				if len(result) > 0 {
					seed[*parent] = append(seed[*parent].([]map[string]any), result)
				}
			}
		}

		return seed
	}

	if data, ok := rawData.(map[string]any); ok {
		for key, value := range data {
			if value == nil {
				continue
			}

			if key == "id" {
				if parent != nil {
					seed[fmt.Sprintf("%vId", *parent)] = fmt.Sprintf("%v", value)
				} else {
					seed[key] = fmt.Sprintf("%v", value)
				}
			} else if reflect.TypeOf(value).String() == "map[string]interface {}" {
				if val, ok := value.(map[string]any); ok {
					for k, v := range val {
						if v == nil {
							continue
						}

						if k == "id" {
							seed[fmt.Sprintf("%vId", key)] = fmt.Sprintf("%v", v)
						} else {
							seed = s.LoadKeys(v, seed, &k)
						}
					}
				}
			} else if slices.Contains([]string{"[]map[string]interface {}", "[]interface {}"}, reflect.TypeOf(value).String()) {
				seed = handleArrayValues(value, seed, &key)
			}
		}
	} else if slices.Contains([]string{"[]map[string]interface {}", "[]interface {}"}, reflect.TypeOf(rawData).String()) && rawData != nil {
		seed = handleArrayValues(rawData, seed, parent)
	}

	return seed
}

func (s *Session) UpdateSessionFlags(model *string) error {
	if model == nil {
		defaultModel := "member"

		model = &defaultModel
	}

	data := s.LoadKeys(s.ActiveData, map[string]any{}, model)

	s.GlobalIds = data

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
	s.GlobalIds = map[string]any{}
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
