package storage

import (
	"context"
	"errors"
	"strings"
)

type ListResult struct {
	IsObject bool
	Path     string
	Size     int
}

var InvalidKey = errors.New("invalid key")
var KeyNotFound = errors.New("key not found")
var KeyAlreadyExists = errors.New("key already exists")
var ValidationError = errors.New("validation error")

type BlobStore interface {
	List(ctx context.Context, prefix string, delimiter string) ([]ListResult, error)
	Get(ctx context.Context, key string) ([]byte, error)
	Put(ctx context.Context, key string, data []byte) error
	Delete(ctx context.Context, key string) error
}

type InMemoryBlobStore struct {
	store map[string][]byte
}

func NewInMemoryBlobStore() *InMemoryBlobStore {
	return &InMemoryBlobStore{make(map[string][]byte)}
}

func (s *InMemoryBlobStore) List(ctx context.Context, prefix string, delimiter string) ([]ListResult, error) {
	result := make([]ListResult, 0)
	for k, v := range s.store {
		if strings.HasPrefix(k, prefix) {
			base := strings.TrimPrefix(k, prefix)
			if delimiter == "" {
				result = append(result, ListResult{IsObject: true, Path: k, Size: len(v)})
				continue
			}
			index := strings.Index(base, delimiter)
			if index == -1 {
				result = append(result, ListResult{IsObject: true, Path: k, Size: len(v)})
				continue
			}
			dir := base[:index+1]
			result = append(result, ListResult{IsObject: false, Path: dir, Size: 0})
		}
	}
	return result, nil
}

func (s *InMemoryBlobStore) Get(ctx context.Context, key string) ([]byte, error) {
	data, ok := s.store[key]
	if !ok {
		return nil, KeyNotFound
	}
	return data, nil
}

func (s *InMemoryBlobStore) Put(ctx context.Context, key string, data []byte) error {
	s.store[key] = data
	return nil
}

func (s *InMemoryBlobStore) Delete(ctx context.Context, key string) error {
	delete(s.store, key)
	return nil
}
