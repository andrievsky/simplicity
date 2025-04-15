package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"simplicity/api"
	"simplicity/config"
	"simplicity/items"
	"simplicity/mock"
	"time"
)

func main() {
	fmt.Println(config.BackendInfo())
	registry := items.NewInMemoryRegistry(time.Now)
	populateWithMockData(registry)
	mux := setupServer(registry)

	slog.Info("Starting server on port", "Port", config.BackendPort)
	http.ListenAndServe(fmt.Sprintf(":%s", config.BackendPort), mux)
}

func setupServer(registry items.Registry) *http.ServeMux {

	itemsApi := api.NewItemHandler(registry)
	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", itemsApi))
	mux.Handle("/", api.WrapHandler(
		http.StripPrefix("/", http.FileServer(http.Dir("../ui/")))))
	mux.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
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
