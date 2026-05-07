package store

import (
	"flux/internal/loader"
	"flux/internal/models"
	"fmt"
	"sync"
	"time"
)

type Store struct {
	mu   sync.RWMutex
	data map[string]models.Entry
}

func NewStore(data models.Data) *Store {
	s := &Store{
		data: make(map[string]models.Entry),
	}

	for k, v := range data.Storage {
		newEntry := models.Entry{
			Value:     v.Value,
			ExpiresAt: v.ExpiresAt,
		}

		s.data[k] = newEntry
	}

	return s
}

func (s *Store) SetValue(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.set(key, value)
}

func (s *Store) GetValue(key string) (any, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.data[key]
	entry, isString := val.Value.(string)

	if isString {
		return entry, ok
	}

	return fmt.Sprintf("%v", val.Value), ok
}

func (s *Store) DeleteValue(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
}

func (s *Store) SetTemporalValue(key, value string, t time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	temp := s.data[key]
	temp.Value = value
	s.data[key] = temp

	expiration := time.Now().Add(t)
	temp.ExpiresAt = &expiration
	s.data[key] = temp
}

func (s *Store) StartCleaner(globalConfig *models.Data, interval time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	go func() {
		ticker := time.NewTicker(interval * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			for key, value := range s.data {
				if value.ExpiresAt != nil && time.Now().After(*value.ExpiresAt) {
					delete(s.data, key)
				}
			}

			r := s.GetAllValues()
			loader.WriteData(globalConfig, r)
		}
	}()
}

func (s *Store) GetAllValues() map[string]models.Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]models.Entry)

	for key, entry := range s.data {
		result[key] = models.Entry{
			Value:     entry.Value,
			ExpiresAt: entry.ExpiresAt,
		}
	}

	return result
}

func (s *Store) LPush(key string, items []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var combined []any

	for _, item := range items {
		combined = append(combined, item)
	}

	if entry, exists := s.data[key]; exists {
		if existingList, ok := entry.Value.([]any); ok {
			combined = append(combined, existingList...)
		} else if existingList, ok := entry.Value.([]string); ok {
			for _, v := range existingList {
				combined = append(combined, v)
			}
		}
	}

	s.set(key, combined)
}

func (s *Store) set(key string, cm any) {
	s.data[key] = models.Entry{
		Value: cm,
	}
}

func convertToAny(in []string) []any {
	out := make([]any, len(in))

	for i, v := range in {
		out[i] = v
	}

	return out
}
