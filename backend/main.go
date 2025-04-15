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
	mux := setupServer(time.Now)

	slog.Info("Starting server on port", "Port", config.BackendPort)
	http.ListenAndServe(fmt.Sprintf(":%s", config.BackendPort), mux)
}

func setupServer(now func() time.Time) *http.ServeMux {
	registry := items.NewInMemoryRegistry(now)
	itemsApi := api.NewItemHandler(registry)
	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", itemsApi))
	mux.Handle("/", api.WrapHandler(
		http.StripPrefix("/", http.FileServer(http.Dir("../ui/")))))
	mux.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	//populateWithMockData(registry)

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
