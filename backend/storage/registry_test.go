package storage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryRegistry_Create(t *testing.T) {
	r := newTestRegistry()
	ctx := context.Background()
	id := "id"
	value := newImageData()
	err := r.Create(ctx, id, value)
	assert.Nil(t, err)
}

func TestInMemoryRegistry_Create_Duplicate(t *testing.T) {
	r := newTestRegistry()
	ctx := context.Background()
	id := "id"
	value := newImageData()
	err := r.Create(ctx, id, value)
	assert.Nil(t, err)

	err = r.Create(ctx, id, value)
	assert.Equal(t, KeyAlreadyExists, err)
}

func TestInMemoryRegistry_Create_EmptyID(t *testing.T) {
	r := newTestRegistry()
	ctx := context.Background()
	id := ""
	value := newImageData()
	err := r.Create(ctx, id, value)
	assert.Equal(t, InvalidKey, err)
}

func TestInMemoryRegistry_Read(t *testing.T) {
	r := newTestRegistry()
	ctx := context.Background()
	id := "id"
	value := newImageData()
	err := r.Create(ctx, id, value)
	assert.Nil(t, err)

	item, err := r.Read(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, id, item.ID)
	assert.Equal(t, value, item.ItemData)
}

func TestInMemoryRegistry_Read_EmptyID(t *testing.T) {
	r := newTestRegistry()
	ctx := context.Background()
	id := ""
	_, err := r.Read(ctx, id)
	assert.Equal(t, InvalidKey, err)
}

func TestInMemoryRegistry_Read_NotFound(t *testing.T) {
	r := newTestRegistry()
	ctx := context.Background()
	id := "id"
	_, err := r.Read(ctx, id)
	assert.Equal(t, KeyNotFound, err)
}

func TestInMemoryRegistry_List(t *testing.T) {
	r := newTestRegistry()
	ctx := context.Background()
	id1 := "id1"
	value1 := newImageData()
	err := r.Create(ctx, id1, value1)
	assert.Nil(t, err)

	id2 := "id2"
	value2 := newImageData()
	err = r.Create(ctx, id2, value2)
	assert.Nil(t, err)

	items, err := r.List(ctx)
	assert.Nil(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, id1, items[0].ID)
	assert.Equal(t, value1, items[0].ItemData)
	assert.Equal(t, id2, items[1].ID)
	assert.Equal(t, value2, items[1].ItemData)
}

func TestInMemoryRegistry_Update(t *testing.T) {
	r := newTestRegistry()
	ctx := context.Background()
	id := "id"
	value := newImageData()
	err := r.Create(ctx, id, value)
	assert.Nil(t, err)

	newValue := newImageData()
	newValue.Title = "new title"
	newValue.Description = "new description"
	newValue.Images = []string{"new image1", "new image2"}
	newValue.Tags = []string{"new tag1", "new tag2"}
	err = r.Update(ctx, id, newValue)
	assert.Nil(t, err)

	item, err := r.Read(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, id, item.ID)
	assert.Equal(t, newValue, item.ItemData)
}

func TestInMemoryRegistry_Update_EmptyID(t *testing.T) {
	r := newTestRegistry()
	ctx := context.Background()
	id := ""
	value := newImageData()
	err := r.Update(ctx, id, value)
	assert.Equal(t, InvalidKey, err)
}

func TestInMemoryRegistry_Delete(t *testing.T) {
	r := newTestRegistry()
	ctx := context.Background()
	id := "id"
	value := newImageData()
	err := r.Create(ctx, id, value)
	assert.Nil(t, err)

	err = r.Delete(ctx, id)
	assert.Nil(t, err)

	_, err = r.Read(ctx, id)
	assert.Equal(t, KeyNotFound, err)
}

func TestInMemoryRegistry_Delete_EmptyID(t *testing.T) {
	r := newTestRegistry()
	ctx := context.Background()
	id := ""
	err := r.Delete(ctx, id)
	assert.Equal(t, InvalidKey, err)
}

func TestInMemoryRegistry_Delete_NotFound(t *testing.T) {
	r := newTestRegistry()
	ctx := context.Background()
	id := "id"
	err := r.Delete(ctx, id)
	assert.Equal(t, KeyNotFound, err)
}

func newTestRegistry() *InMemoryItemRegistry {
	return NewInMemoryItemRegistry(time.Now)
}

func newImageData() ItemData {
	return ItemData{
		Title:       "title",
		Description: "description",
		Images:      []string{"image1", "image2"},
		Tags:        []string{"tag1", "tag2"},
	}
}
