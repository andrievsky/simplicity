package items

import (
	"context"
	"errors"
	"simplicity/oops"
	"sort"
	"time"
)

type Registry interface {
	Create(ctx context.Context, id string, value ItemData) error
	Read(ctx context.Context, id string) (Item, error)
	List(ctx context.Context) ([]Item, error)
	Update(ctx context.Context, id string, value ItemData) error
	Delete(ctx context.Context, id string) error
}

type InMemoryRegistry struct {
	store map[string]Item
	now   func() time.Time
}

func NewInMemoryRegistry(now func() time.Time) *InMemoryRegistry {
	return &InMemoryRegistry{make(map[string]Item), now}
}

func (r *InMemoryRegistry) Create(ctx context.Context, id string, value ItemData) error {
	if id == "" {
		return oops.InvalidKey
	}
	err := validateItemData(value)
	if err != nil {
		return errors.Join(oops.ValidationError, err)
	}
	if _, ok := r.store[id]; ok {
		return oops.KeyAlreadyExists
	}
	r.store[id] = Item{
		ItemMetadata: ItemMetadata{
			ID:        id,
			CreatedAt: r.now(),
			UpdatedAt: r.now(),
		},
		ItemData: value,
	}
	return nil
}

func (r *InMemoryRegistry) Read(ctx context.Context, id string) (Item, error) {
	if id == "" {
		return Item{}, oops.InvalidKey
	}
	item, ok := r.store[id]
	if !ok {
		return Item{}, oops.KeyNotFound
	}
	return item, nil
}

func (r *InMemoryRegistry) List(ctx context.Context) ([]Item, error) {
	items := make([]Item, len(r.store))
	i := 0
	for _, item := range r.store {
		items[i] = item
		i++
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})
	return items, nil
}

func (r *InMemoryRegistry) Update(ctx context.Context, id string, value ItemData) error {
	if id == "" {
		return oops.InvalidKey
	}
	err := validateItemData(value)
	if err != nil {
		return err
	}
	item, ok := r.store[id]
	if !ok {
		return oops.KeyNotFound
	}
	item.ItemData = value
	item.UpdatedAt = r.now()
	r.store[id] = item
	return nil
}

func (r *InMemoryRegistry) Delete(ctx context.Context, id string) error {
	if id == "" {
		return oops.InvalidKey
	}
	if _, ok := r.store[id]; !ok {
		return oops.KeyNotFound
	}
	delete(r.store, id)
	return nil
}
