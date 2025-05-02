package storage

import (
	"context"
	"io"
	"strings"
)

type ListResult struct {
	IsObject bool
	Key      string
	Size     int
}

const delimiter = "/"

type BlobStore interface {
	List(ctx context.Context, prefix string, delimiter string) ([]ListResult, error)
	Get(ctx context.Context, key string) (io.ReadCloser, map[string]string, error)
	Put(ctx context.Context, key string, reader io.Reader, metadata map[string]string) error
	Delete(ctx context.Context, key string) error
	DeleteAll(ctx context.Context, prefix string) error
}

func JoinPath(elem ...string) string {
	if len(elem) == 0 {
		return ""
	}
	return strings.Join(elem, delimiter)
}
