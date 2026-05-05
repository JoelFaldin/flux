package store

import (
	"flux/internal/loader"
	"sync"
	"time"
)

type entry struct {
	value     string
	expiresAt *time.Time
}

type Store struct {
	mu   sync.RWMutex
	data map[string]entry
}

func NewStore(data *loader.Data) *Store {
	s := &Store{
		data: make(map[string]entry),
	}

	for k, v := range data.Storage {
		newEntry := entry{
			value:     v,
			expiresAt: nil,
		}

		s.data[k] = newEntry
	}

	return s
}

func (s *Store) SetValue(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	temp := s.data[key]
	temp.value = value
	s.data[key] = temp
}

func (s *Store) GetValue(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.data[key]

	return val.value, ok
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
	temp.value = value
	s.data[key] = temp

	expiration := time.Now().Add(t * time.Second)
	temp.expiresAt = &expiration
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
				if value.expiresAt != nil && time.Now().After(*value.expiresAt) {
					delete(s.data, key)
				}
			}
		}
	}()
}
