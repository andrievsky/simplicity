package items

import (
	"errors"
	"time"
)

type Item struct {
	ItemMetadata
	ItemData
}

type ItemMetadata struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ItemData struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
	Tags        []string `json:"tags"`
}

type Tag string

const TagSeparator = ":"

func NewTag(key, value string) Tag {
	bytes := make([]byte, len(key)+len(value)+len(TagSeparator))
	copy(bytes, key)
	copy(bytes[len(key):], TagSeparator)
	copy(bytes[len(key)+len(TagSeparator):], value)
	return Tag(bytes)
}

func validateItemData(item ItemData) error {
	if item.Title == "" {
		return errors.New("title is required")
	}
	return nil
}
