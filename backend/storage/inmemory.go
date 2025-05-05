package storage

import (
	"context"
	"io"
	"simplicity/oops"
	"strings"
)

type InMemoryBlobStore struct {
	store    map[string][]byte
	metadata map[string]map[string]string
}

func NewInMemoryBlobStore() *InMemoryBlobStore {
	return &InMemoryBlobStore{make(map[string][]byte), make(map[string]map[string]string)}
}

func (s *InMemoryBlobStore) List(ctx context.Context, prefix string, delimiter string) ([]ListResult, error) {
	result := make([]ListResult, 0)
	for k, v := range s.store {
		if strings.HasPrefix(k, prefix) {
			base := strings.TrimPrefix(k, prefix)
			if delimiter == "" {
				result = append(result, ListResult{IsObject: true, Key: k, Size: len(v)})
				continue
			}
			index := strings.Index(base, delimiter)
			if index == -1 {
				result = append(result, ListResult{IsObject: true, Key: k, Size: len(v)})
				continue
			}
			dir := base[:index+1]
			result = append(result, ListResult{IsObject: false, Key: dir, Size: 0})
		}
	}
	return result, nil
}

func (s *InMemoryBlobStore) Get(ctx context.Context, key string) (io.ReadCloser, map[string]string, error) {
	data, ok := s.store[key]
	if !ok {
		return nil, nil, oops.KeyNotFound
	}
	reader := io.NopCloser(strings.NewReader(string(data)))
	metadata, ok := s.metadata[key]
	if !ok {
		metadata = make(map[string]string)
	}
	return reader, metadata, nil
}
func (s *InMemoryBlobStore) Put(ctx context.Context, key string, reader io.Reader, metadata map[string]string) error {
	if key == "" {
		return oops.InvalidKey
	}
	if metadata == nil {
		metadata = make(map[string]string)
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	s.store[key] = data
	s.metadata[key] = metadata
	return nil
}

func (s *InMemoryBlobStore) Delete(ctx context.Context, key string) error {
	delete(s.store, key)
	return nil
}

func (s *InMemoryBlobStore) DeleteAll(ctx context.Context, prefix string) error {
	if prefix == "" {
		return oops.InvalidKey
	}
	if !strings.HasSuffix(prefix, Delimiter) {
		prefix += Delimiter
	}
	for k := range s.store {
		if strings.HasPrefix(k, prefix) {
			delete(s.store, k)
		}
	}
	return nil
}
