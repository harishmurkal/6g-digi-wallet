package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

type fileRecord struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

type FileStore struct {
	mu   sync.RWMutex
	path string
	data map[string]json.RawMessage
}

func NewFileStore(path string) (*FileStore, error) {
	store := &FileStore{
		path: path,
		data: make(map[string]json.RawMessage),
	}
	// Load existing data if file exists
	file, err := os.Open(path)
	if err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			var rec fileRecord
			if json.Unmarshal(scanner.Bytes(), &rec) == nil {
				store.data[rec.Key] = rec.Value
			}
		}
	}
	return store, nil
}

func (f *FileStore) persist() error {
	file, err := os.Create(f.path)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for k, v := range f.data {
		rec := fileRecord{Key: k, Value: v}
		line, _ := json.Marshal(rec)
		writer.WriteString(string(line) + "\n")
	}
	return writer.Flush()
}

func (f *FileStore) Save(key string, value any) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	f.data[key] = data
	return f.persist()
}

func (f *FileStore) Load(key string, out any) error {
	f.mu.RLock()
	defer f.mu.RUnlock()
	data, ok := f.data[key]
	if !ok {
		return fmt.Errorf("key not found: %s", key)
	}
	return json.Unmarshal(data, out)
}

func (f *FileStore) ListKeys(prefix string) ([]string, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	var keys []string
	for k := range f.data {
		if prefix == "" || strings.HasPrefix(k, prefix) {
			keys = append(keys, k)
		}
	}
	return keys, nil
}

func (f *FileStore) String() string {
	return fmt.Sprintf("FileStore[%s]", f.path)
}
