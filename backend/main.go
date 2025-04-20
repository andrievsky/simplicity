package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"simplicity/config"
	"simplicity/images"
	"simplicity/items"
	"simplicity/mock"
	"simplicity/storage"
	"simplicity/svc"
	"time"
)

func main() {
	fmt.Println(config.BackendInfo())
	registry := items.NewInMemoryRegistry(time.Now)
	store := storage.NewInMemoryBlobStore()
	populateWithMockData(registry)
	mux := setupServer(registry, store)

	slog.Info("Starting server on port", "Port", config.BackendPort)
	http.ListenAndServe(fmt.Sprintf(":%s", config.BackendPort), mux)
}

func setupServer(registry items.Registry, store storage.BlobStore) *http.ServeMux {

	itemsApi := items.NewItemHandler(registry)
	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", itemsApi))
	mux.Handle("/", svc.WrapHandler(
		http.StripPrefix("/", http.FileServer(http.Dir("../ui/")))))
	mux.Handle("/api/image/", http.StripPrefix("/api/image",
		images.NewImageApi(storage.NewPrefixBlobStore(store, "image/"))))

	return mux
}

func populateWithMockData(registry items.Registry) {
	ctx := context.Background()
	data := mock.GenerateMockData()
	for i, item := range data {
		err := registry.Create(ctx, fmt.Sprintf("%d", i), item)
		if err != nil {
			slog.Error("Error creating item", "Error", err)
		}
	}
}
