// internal/storage/store.go
package storage

// Store defines the contract for storage backends.
type Store interface {
	Save(key string, value any) error
	Load(key string, out any) error
	ListKeys(prefix string) ([]string, error)
}
