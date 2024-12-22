package storage

import (
	"context"
	"testing"
)

func TestInMemoryBlobStore_List(t *testing.T) {
	type PutObject struct {
		key  string
		data string
	}
	type ListObject struct {
		prefix    string
		delimiter string
		data      []ListResult
	}
	tests := []struct {
		put  []PutObject
		list []ListObject
	}{
		{
			put: []PutObject{{key: "prefix1", data: "data1"}},
			list: []ListObject{
				{prefix: "prefix", delimiter: "/", data: []ListResult{{IsObject: true, Path: "prefix1", Size: 5}}},
				{prefix: "prefix", delimiter: "", data: []ListResult{{IsObject: true, Path: "prefix1", Size: 5}}},
				{prefix: "prefix", delimiter: "////", data: []ListResult{{IsObject: true, Path: "prefix1", Size: 5}}},
			},
		},
	}

	for _, tt := range tests {
		store := NewInMemoryBlobStore()
		ctx := context.Background()
		for _, o := range tt.put {
			store.Put(ctx, o.key, []byte(o.data))
		}
		for _, o := range tt.list {
			result, err := store.List(ctx, o.prefix, o.delimiter)
			if err != nil {
				t.Errorf("List() error = %v", err)
				return
			}
			if len(result) != len(o.data) {
				t.Errorf("List() = %v, want %v", result, o.data)
				return
			}
			for i, r := range result {
				if r.IsObject != o.data[i].IsObject {
					t.Errorf("List() = %v, want %v", result, o.data)
					return
				}
				if r.Path != o.data[i].Path {
					t.Errorf("List() = %v, want %v", result, o.data)
					return
				}
				if r.Size != o.data[i].Size {
					t.Errorf("List() = %v, want %v", result, o.data)
					return
				}
			}
		}
	}
}
