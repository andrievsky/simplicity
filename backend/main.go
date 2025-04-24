package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"simplicity/config"
	"simplicity/genid"
	"simplicity/images"
	"simplicity/items"
	"simplicity/mock"
	"simplicity/storage"
	"simplicity/svc"
	"time"
)

func main() {
	fmt.Printf("Backend %s Version %s\n", config.BackendName, config.BackendVersion)
	registry := items.NewInMemoryRegistry(time.Now)
	store := storage.NewInMemoryBlobStore()
	populateWithMockData(registry)
	mux := setupServer(registry, store)

	slog.Info("Starting server on port", "Port", config.BackendPort)
	http.ListenAndServe(fmt.Sprintf(":%s", config.BackendPort), mux)
}

func setupServer(registry items.Registry, store storage.BlobStore) *http.ServeMux {
	idProvider, err := genid.NewSnowflakeProvider(1)
	if err != nil {
		panic(err)
	}
	itemsApi := items.NewItemHandler(registry)
	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", itemsApi))
	mux.Handle("/", svc.WrapHandler(
		http.StripPrefix("/", http.FileServer(http.Dir("../ui/")))))
	mux.Handle("/api/image/", svc.WrapHandler(http.StripPrefix("/api/image",
		images.NewImageApi(storage.NewPrefixBlobStore(store, "image/"), idProvider))))

	mux.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
		svc.WriteData(w, r, config.BackendInfo(), http.StatusOK)
	})

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
