package store

import (
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

func (s *Store) SetValue(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	temp := s.data[key]
	temp.Value = value
	s.data[key] = temp
}

func (s *Store) GetValue(key string) (string, bool) {
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

	expiration := time.Now().Add(t * time.Second)
	temp.ExpiresAt = &expiration
	s.data[key] = temp
}

func (s *Store) StartCleaner(interval time.Duration) {
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
		}
	}()
}

func (s *Store) GetAllValues() map[string]any {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]any)

	// Type assertion:
	for key, entry := range s.data {
		result[key] = entry.Value
	}

	return result
}

func (s *Store) LPush(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// fmt.Printf("LPUSH %s %s\n", key, value)
}
