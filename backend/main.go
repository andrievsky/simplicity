package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"log"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	_ "net/http/pprof"
	"os"
	"simplicity/config"
	"simplicity/genid"
	"simplicity/images"
	"simplicity/items"
	"simplicity/mock"
	"simplicity/storage"
	"simplicity/svc"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("cannot load config: %w", err))
	}
	if conf.EnableDebug {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}
	fmt.Printf("Backend %s Version %s\n", conf.BackendName, conf.BackendVersion)

	//store := storage.NewInMemoryBlobStore()
	store := storage.NewS3BlobStore(setupS3Client(conf), conf.AWS.Bucket)
	registry := items.NewPersistentRegistry(storage.NewPrefixBlobStore(store, "item/"), "items.js")
	err = registry.Init()
	if err != nil {
		panic(fmt.Errorf("cannot init registry: %w", err))
	}

	mux := setupServer(registry, store, conf)
	//populateWithMockData(registry, mux)

	slog.Info("Starting server on port", "Port", conf.Server.Port)
	http.ListenAndServe(":"+conf.Server.Port, mux)
}

func setupServer(registry items.Registry, store storage.BlobStore, conf *config.Config) *http.ServeMux {
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
		svc.WriteData(w, r, conf.BackendVersion, http.StatusOK)
	})

	return mux
}

func setupS3Client(conf *config.Config) *s3.Client {
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithSharedConfigProfile(conf.AWS.Profile),
	)
	if err != nil {
		log.Fatalf("unable to load AWS SDK config: %v", err)
	}

	// Create and return the S3 client
	return s3.NewFromConfig(cfg)
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
		ID    string `json:"id"`
		Error string `json:"error"`
	}
	var resp uploadResponse
	err = json.NewDecoder(rec.Body).Decode(&resp)
	if err != nil {
		panic(err)
	}
	if rec.Code != http.StatusCreated {
		panic(fmt.Sprintf("expected status code 201, got %d: %s", rec.Code, resp.Error))
	}

	return resp.ID
}
