package storage

import "fmt"

type BackendType string

const (
	BackendMemory BackendType = "memory"
	BackendFile   BackendType = "file"
	BackendRedis  BackendType = "redis"
)

func NewStore(backend BackendType, opts map[string]string) (Store, error) {
	switch backend {
	case BackendMemory:
		return NewMemoryStore(), nil
	case BackendFile:
		path := opts["path"]
		if path == "" {
			path = "./data/store.jsonl"
		}
		return NewFileStore(path)
	case BackendRedis:
		return nil, fmt.Errorf("RedisStore not implemented yet")
	default:
		return nil, fmt.Errorf("unknown backend: %s", backend)
	}
}
