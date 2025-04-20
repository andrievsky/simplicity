package storage

import (
	"context"
	"io"
	"simplicity/oops"
)

type StripPrefixBlobStore struct {
	store  BlobStore
	prefix string
}

func NewPrefixBlobStore(store BlobStore, prefix string) *StripPrefixBlobStore {
	return &StripPrefixBlobStore{store: store, prefix: prefix}
}

func (s *StripPrefixBlobStore) List(ctx context.Context, prefix string, delimiter string) ([]ListResult, error) {
	return s.store.List(ctx, s.prefix+prefix, delimiter)
}

func (s *StripPrefixBlobStore) Get(ctx context.Context, key string) (io.ReadCloser, map[string]string, error) {
	if key == "" {
		return nil, nil, oops.InvalidKey
	}
	return s.store.Get(ctx, s.prefix+key)
}

func (s *StripPrefixBlobStore) Put(ctx context.Context, key string, reader io.Reader, metadata map[string]string) error {
	if key == "" {
		return oops.InvalidKey
	}
	return s.store.Put(ctx, s.prefix+key, reader, metadata)
}

func (s *StripPrefixBlobStore) Delete(ctx context.Context, key string) error {
	if key == "" {
		return oops.InvalidKey
	}
	return s.store.Delete(ctx, s.prefix+key)
}

func (s *StripPrefixBlobStore) DeleteAll(ctx context.Context, prefix string) error {
	if prefix == "" {
		return oops.InvalidKey
	}
	return s.store.DeleteAll(ctx, s.prefix+prefix)
}
