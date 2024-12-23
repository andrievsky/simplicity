package main

import (
	"fmt"
	"net/http"
	"simplicity/api"
	"simplicity/config"
	"simplicity/items"
	"time"
)

func main() {
	fmt.Println(config.BackendInfo())
	mux := setupServer(time.Now)

	http.ListenAndServe(fmt.Sprintf(":%s", config.BackendPort), mux)
}

func setupServer(now func() time.Time) *http.ServeMux {
	registry := items.NewInMemoryRegistry(now)
	api := api.NewItemHandler(registry)
	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", api))
	mux.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return mux
}
