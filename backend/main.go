package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
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

	mux := setupServer(registry, store)
	populateWithMockData(registry, mux)

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

func populateWithMockData(registry items.Registry, mux http.Handler) {
	ctx := context.Background()
	data := mock.GenerateMockData()
	imageIDs := []string{uploadTestImage(mux, "./files/test-image-1.png"), uploadTestImage(mux, "./files/test-image-2.png")}
	for i, item := range data {
		item.Images = imageIDs

		err := registry.Create(ctx, fmt.Sprintf("%d", i), item)
		if err != nil {
			slog.Error("Error creating item", "Error", err)
		}
	}
}

func uploadTestImage(mux http.Handler, filepath string) string {
	file, err := os.Open(filepath)
	if err != nil {
		panic(fmt.Errorf("cannot open image file %s: %w", filepath, err))
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		panic(err)
	}

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/image/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	type uploadResponse struct {
		ID string `json:"id"`
	}
	var resp uploadResponse
	err = json.NewDecoder(rec.Body).Decode(&resp)
	if err != nil {
		panic(err)
	}

	return resp.ID
}
