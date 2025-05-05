package main

import (
	"context"
	"fmt"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"
	"simplicity/config"
	"simplicity/genid"
	"simplicity/images"
	"simplicity/items"
	"simplicity/loggers"
	"simplicity/storage"
	"simplicity/svc"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("cannot load config: %w", err))
	}
	logger := loggers.NewLogger(conf)
	slog.SetDefault(logger)
	if conf.EnableDebug {
		go func() {
			logger.Debug("Starting pprof server", "Port", "6060")
			http.ListenAndServe("localhost:6060", nil)
		}()
	}

	logger.Info("Backend info", "Version", conf.BackendVersion, "Name", conf.BackendName)
	buildInfo, _ := debug.ReadBuildInfo()
	logger.Debug("Build info", "Version", buildInfo.Main.Version, "Path", buildInfo.Main.Path, "GoVersion", buildInfo.GoVersion, "Settings", buildInfo.Settings)

	//store := storage.NewInMemoryBlobStore()
	s3Client, err := setupS3Client(conf)
	if err != nil {
		panic(fmt.Errorf("cannot create S3 client: %w", err))
	}
	store := storage.NewS3BlobStore(s3Client, conf.AWS.Bucket)
	registry := items.NewPersistentRegistry(store, "item/items.js")
	err = registry.Init()
	if err != nil {
		panic(fmt.Errorf("cannot init registry: %w", err))
	}

	handler := svc.NewLoggingMiddleware(setupServer(registry, store, conf, logger), logger)

	//populateWithMockData(registry, mux)

	logger.Info("Starting backend service", "Port", conf.Server.Port)
	http.ListenAndServe(":"+conf.Server.Port, handler)
}

func setupServer(registry items.Registry, store storage.BlobStore, conf *config.Config, logger *slog.Logger) http.Handler {
	idProvider, err := genid.NewSnowflakeProvider(1)
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("../ui/"))))
	mux.Handle("/api/item/", http.StripPrefix("/api/item", items.NewApi(registry, idProvider, logger)))
	mux.Handle("/api/image/", http.StripPrefix("/api/image", images.NewApi(store, idProvider, logger)))
	mux.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
		svc.Data(w, r, conf.BackendVersion, http.StatusOK)
	})

	return mux
}

func setupS3Client(conf *config.Config) (*s3.Client, error) {
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithSharedConfigProfile(conf.AWS.Profile),
	)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(cfg), nil
}
