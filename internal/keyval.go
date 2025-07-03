package internal

import "sync"

// KV is an in-memory key-value store with thread-safe access.
type KV struct {
	mu   sync.RWMutex
	data map[string][]byte
}

// NewKV creates and returns a new instance of KV.
func NewKV() *KV {
	return &KV{
		data: map[string][]byte{},
	}
}

// Set stores the given value for the specified key in the KV store.
// It overwrites any existing value for the key.
func (kv *KV) Set(key, val []byte) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.data[string(key)] = []byte(val)
	return nil
}

// Get retrieves the value associated with the specified key from the KV store.
// It returns the value and a boolean indicating whether the key was found.
func (kv *KV) Get(key []byte) ([]byte, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	val, ok := kv.data[string(key)]
	return val, ok
}
