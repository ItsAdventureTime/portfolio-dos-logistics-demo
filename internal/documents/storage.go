package documents

import (
	"context"
	"fmt"
	"io"
	"sync"
)

// Storage is the interface for private document storage. Files are stored
// behind authenticated authorization checks and never served publicly.
type Storage interface {
	// Store saves a file at the given storage key.
	Store(ctx context.Context, key string, content []byte, contentType string) error
	// Retrieve returns the file content at the given key.
	Retrieve(ctx context.Context, key string) (content []byte, contentType string, err error)
	// Delete removes a file at the given key.
	Delete(ctx context.Context, key string) error
	// Exists reports whether a file exists at the given key.
	Exists(ctx context.Context, key string) bool
}

// MemStorage is an in-memory Storage implementation for development and testing.
type MemStorage struct {
	mu    sync.RWMutex
	files map[string]storedFile
}

type storedFile struct {
	content     []byte
	contentType string
}

func NewMemStorage() *MemStorage {
	return &MemStorage{files: make(map[string]storedFile)}
}

func (s *MemStorage) Store(ctx context.Context, key string, content []byte, contentType string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.files[key] = storedFile{content: content, contentType: contentType}
	return nil
}

func (s *MemStorage) Retrieve(ctx context.Context, key string) ([]byte, string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	f, ok := s.files[key]
	if !ok {
		return nil, "", fmt.Errorf("file not found: %s", key)
	}
	return f.content, f.contentType, nil
}

func (s *MemStorage) Delete(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.files, key)
	return nil
}

func (s *MemStorage) Exists(ctx context.Context, key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.files[key]
	return ok
}

// Count returns the number of stored files (for testing).
func (s *MemStorage) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.files)
}

// _ retains the io import for future streaming adapters.
var _ = io.EOF