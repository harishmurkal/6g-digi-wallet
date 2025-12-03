package storage

import (
	"os"
	"reflect"
	"testing"
)

// --- Test helpers ---

type dummyStruct struct {
	Name string
	Age  int
}

// --- Tests for MemoryStore ---

func TestMemoryStore_SaveLoad(t *testing.T) {
	store := NewMemoryStore()

	input := dummyStruct{Name: "Alice", Age: 30}
	err := store.Save("user:1", input)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	var output dummyStruct
	err = store.Load("user:1", &output)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if !reflect.DeepEqual(input, output) {
		t.Errorf("expected %+v, got %+v", input, output)
	}
}

func TestMemoryStore_SaveEmptyKey(t *testing.T) {
	store := NewMemoryStore()
	err := store.Save("", dummyStruct{})
	if err == nil {
		t.Error("expected error for empty key, got nil")
	}
}

func TestMemoryStore_ListKeys(t *testing.T) {
	store := NewMemoryStore()
	store.Save("user:1", dummyStruct{})
	store.Save("user:2", dummyStruct{})
	store.Save("config:sys", dummyStruct{})

	keys, err := store.ListKeys("user:")
	if err != nil {
		t.Fatalf("ListKeys failed: %v", err)
	}

	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}

func TestMemoryStore_LoadNonexistent(t *testing.T) {
	store := NewMemoryStore()
	var v dummyStruct
	err := store.Load("missing", &v)
	if err == nil {
		t.Error("expected error for non-existent key")
	}
}

// --- Tests for FileStore ---

func TestFileStore_SaveLoad(t *testing.T) {
	tmpFile := "test_store.jsonl"
	defer os.Remove(tmpFile)

	store, err := NewFileStore(tmpFile)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}

	input := dummyStruct{Name: "Bob", Age: 42}
	err = store.Save("user:bob", input)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	var output dummyStruct
	err = store.Load("user:bob", &output)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if !reflect.DeepEqual(input, output) {
		t.Errorf("expected %+v, got %+v", input, output)
	}

	// Test ListKeys
	keys, err := store.ListKeys("user:")
	if err != nil {
		t.Fatalf("ListKeys failed: %v", err)
	}
	if len(keys) != 1 {
		t.Errorf("expected 1 key, got %d", len(keys))
	}
}

func TestFileStore_PersistenceAcrossInstances(t *testing.T) {
	tmpFile := "test_persist.jsonl"
	defer os.Remove(tmpFile)

	// First instance: save
	store1, err := NewFileStore(tmpFile)
	if err != nil {
		t.Fatalf("NewFileStore failed: %v", err)
	}
	store1.Save("persist:key", dummyStruct{Name: "Eve", Age: 28})

	// Second instance: load
	store2, err := NewFileStore(tmpFile)
	if err != nil {
		t.Fatalf("reopen store failed: %v", err)
	}

	var out dummyStruct
	err = store2.Load("persist:key", &out)
	if err != nil {
		t.Fatalf("Load after reopen failed: %v", err)
	}

	if out.Name != "Eve" || out.Age != 28 {
		t.Errorf("expected persisted value Eve/28, got %+v", out)
	}
}

func TestFileStore_LoadNonexistent(t *testing.T) {
	tmpFile := "test_nonexistent.jsonl"
	defer os.Remove(tmpFile)

	store, _ := NewFileStore(tmpFile)
	var v dummyStruct
	err := store.Load("nope", &v)
	if err == nil {
		t.Error("expected error for non-existent key")
	}
}
