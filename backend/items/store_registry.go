package items

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"simplicity/oops"
	"simplicity/storage"
	"time"
)

type StoreRegistry struct {
	store    storage.BlobStore
	key      string
	registry *InMemoryRegistry
}

func NewPersistentRegistry(store storage.BlobStore, key string) *StoreRegistry {
	return &StoreRegistry{store, key, NewInMemoryRegistry(func() time.Time {
		return time.Now()
	})}
}

func (r *StoreRegistry) Init() error {
	reader, _, err := r.store.Get(context.Background(), r.key)
	if err != nil {
		if err == oops.KeyNotFound {
			return nil
		}
		return err
	}
	if reader == nil {
		return nil
	}
	defer reader.Close()

	var items map[string]Item
	err = json.NewDecoder(reader).Decode(&items)
	if err != nil {
		return fmt.Errorf("failed to decode blob: %w", err)
	}

	r.registry.store = items
	return nil
}

func (r *StoreRegistry) Create(ctx context.Context, id string, value ItemData) error {
	err := r.registry.Create(ctx, id, value)
	if err != nil {
		return err
	}
	return r.flush(ctx)
}

func (r *StoreRegistry) Read(ctx context.Context, id string) (Item, error) {
	return r.registry.Read(ctx, id)
}

func (r *StoreRegistry) List(ctx context.Context) ([]Item, error) {
	return r.registry.List(ctx)
}

func (r *StoreRegistry) Update(ctx context.Context, id string, value ItemData) error {
	err := r.registry.Update(ctx, id, value)
	if err != nil {
		return err
	}
	return r.flush(ctx)
}

func (r *StoreRegistry) Delete(ctx context.Context, id string) error {
	err := r.registry.Delete(ctx, id)
	if err != nil {
		return err
	}
	return r.flush(ctx)
}

func (r *StoreRegistry) flush(ctx context.Context) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(r.registry.store)
	if err != nil {
		return fmt.Errorf("failed to encode blob: %w", err)
	}

	return r.store.Put(ctx, r.key, &buf, nil)
}
