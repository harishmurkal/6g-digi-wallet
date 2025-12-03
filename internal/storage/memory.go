package storage

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

type MemoryStore struct {
	// RWMutex is correct for concurrency (multiple readers, single writer)
	mu   sync.RWMutex
	data map[string][]byte
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string][]byte),
	}
}

// Save will ALWAYS overwrite if the key already exists, preventing duplicates.
func (m *MemoryStore) Save(key string, value any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check for empty key is a good practice, though not strictly necessary
	if key == "" {
		logError("Key(%s) to save is empty", key)
		return fmt.Errorf("key cannot be empty")
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// If key exists, m.data[key] = data overwrites it. No duplicate entry is created.
	m.data[key] = data
	logInfo("Added %s", key)
	return nil
}

func (m *MemoryStore) Load(key string, out any) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data, ok := m.data[key]
	if !ok {
		return fmt.Errorf("key not found: %s", key)
	}
	return json.Unmarshal(data, out)
}

// ListKeys iterates over the map keys, which are guaranteed to be unique.
func (m *MemoryStore) ListKeys(prefix string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Optimization: Allocate capacity based on the total map size to avoid reallocations
	keys := make([]string, 0, len(m.data))

	for k := range m.data {
		if prefix == "" || strings.HasPrefix(k, prefix) {
			keys = append(keys, k)
		}
	}
	return keys, nil
}

func (m *MemoryStore) String() string {
	return "MemoryStore"
}
