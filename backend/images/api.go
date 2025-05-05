package images

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"simplicity/genid"
	"simplicity/oops"
	"simplicity/storage"
	"simplicity/svc"
	"strings"
	"time"
)

type Api struct {
	store        storage.BlobStore
	deletedStore storage.BlobStore
	idProvider   genid.Provider
	logger       *slog.Logger
}

type Image struct {
	ID       string `json:"id"`
	Location string `json:"location"`
}

const maxUploadSize = 48 * 1024 * 1024 // 48MB

func NewApi(store storage.BlobStore, idProvider genid.Provider, logger *slog.Logger) *http.ServeMux {
	router := http.NewServeMux()
	api := &Api{
		storage.NewPrefixBlobStore(store, "images/files/"),
		storage.NewPrefixBlobStore(store, "images/deleted-files/"),
		idProvider,
		logger.With("component", "images"),
	}

	router.HandleFunc("GET /files/", api.list)
	router.HandleFunc("POST /upload", api.post)
	router.HandleFunc("GET /files/{id}", api.get)
	router.HandleFunc("DELETE /files/{id}", api.delete)

	return router
}

func storagePath(id string, format *Format) string {
	return storage.JoinPath(id, format.FileName())
}

func (api *Api) get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := api.idProvider.Validate(id); err != nil {
		svc.Error(w, r, err)
		return
	}
	format, err := resolveFormat(r.URL.Query().Get("format"))
	if err != nil {
		svc.Error(w, r, err)
		return
	}
	api.logger.DebugContext(r.Context(), "GET", "path", r.URL.Path)
	path := storagePath(id, format)
	reader, metadata, err := api.store.Get(r.Context(), path)
	if err == oops.KeyNotFound {
		err = api.createImageVariant(r.Context(), id, format)
		if err == nil {
			reader, metadata, err = api.store.Get(r.Context(), path)
		}
	}
	if err != nil {
		svc.Error(w, r, err)
		return
	}
	defer reader.Close()
	ext := format.Ext
	if format == Source {
		ext = MetadataReader{metadata}.Extension()
	}
	header := w.Header()
	header.Set("Content-Type", resolveMime(ext))
	if metadata != nil {
		for k, v := range metadata {
			header.Set("metadata-"+k, v)
		}
	}

	if _, err = io.Copy(w, reader); err != nil {
		api.logger.ErrorContext(r.Context(), "Error during response writing", "method", "GET", "Error:", err.Error())
		return
	}
}

func (api *Api) createImageVariant(ctx context.Context, id string, format *Format) error {
	reader, metadata, err := api.store.Get(ctx, storagePath(id, Canonical))
	if err != nil {
		return fmt.Errorf("failed to get canonical image: %w", err)
	}
	defer reader.Close()
	tReader, err := transcodeFile(Canonical, format, reader)
	if err != nil {
		return fmt.Errorf("failed to transcode image: %w", err)
	}
	return api.store.Put(ctx, storagePath(id, format), tReader, metadata)
}

func (api *Api) list(w http.ResponseWriter, r *http.Request) {
	api.logger.DebugContext(r.Context(), "LIST", "path", r.URL.Path)
	seq, err := api.store.List(r.Context(), "", storage.Delimiter)
	if err != nil {
		svc.Error(w, r, err)
		return
	}
	images := make([]string, 0, len(seq))
	for _, item := range seq {
		if !item.IsObject {
			images = append(images, strings.TrimSuffix(item.Key, "/"))
		}
	}
	svc.Data(w, r, images, http.StatusOK)
}

func (api *Api) post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		svc.ErrorWithCode(w, r, fmt.Errorf("failed to parse form: %w", err), http.StatusBadRequest)
		return
	}
	if r.MultipartForm == nil || len(r.MultipartForm.File) == 0 {
		svc.ErrorWithCode(w, r, errors.New("missing file field"), http.StatusBadRequest)
		return
	}
	if len(r.MultipartForm.File) > 1 {
		svc.Error(w, r, errors.New("only one file is supported"))
		return
	}
	fileHeader := r.MultipartForm.File["file"][0]
	if fileHeader == nil {
		svc.Error(w, r, errors.New("missing file field"))
		return
	}
	id := api.idProvider.Generate()
	metadata, err := buildMetadata(id, fileHeader)
	if err != nil {
		svc.Error(w, r, fmt.Errorf("failed to build metadata: %w", err))
		return
	}
	api.logger.DebugContext(r.Context(), "Creating image", "method", "POST", "id", id, "metadata", metadata)
	sourcePath := storagePath(metadata.ID, Source)
	canonicalPath := storagePath(metadata.ID, Canonical)

	file, err := fileHeader.Open()
	if err != nil {
		svc.Error(w, r, errors.New("failed to open file"))
		return
	}
	defer file.Close()

	err = api.store.Put(r.Context(), sourcePath, file, metadata.Map())
	if err != nil {
		svc.Error(w, r, fmt.Errorf("failed to store image: %w", err))
		return
	}

	tFile, err := fileHeader.Open()
	defer tFile.Close()
	tFormat := &Format{
		Name:   "",
		Ext:    metadata.Extension,
		Width:  0,
		Height: 0,
	}
	tr, err := transcodeFile(tFormat, Canonical, tFile)
	if err != nil {
		svc.Error(w, r, fmt.Errorf("failed to transcode file: %w", err))
		return
	}

	err = api.store.Put(r.Context(), canonicalPath, tr, metadata.Map())
	if err != nil {
		svc.Error(w, r, fmt.Errorf("failed to store transcoded image: %w", err))
		return
	}
	api.logger.InfoContext(r.Context(), "Image created", "method", "POST", "id", id, "metadata", metadata)
	w.Header().Set("Location", "images/files/"+metadata.ID)
	svc.Data(w, r, Image{ID: metadata.ID}, http.StatusCreated)
}

func transcodeFile(in *Format, out *Format, r io.Reader) (io.Reader, error) {
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		err := Transcode(in, out, r, pw)
		if err != nil {
			pw.CloseWithError(fmt.Errorf("failed to transcode file: %w", err))
		}
	}()

	return pr, nil
}

func buildMetadata(id string, fileHeader *multipart.FileHeader) (Metadata, error) {
	originalName := fileHeader.Filename
	ext, err := resolveExtFromFileName(originalName)
	if err != nil {
		return Metadata{}, fmt.Errorf("failed to resolve extension: %w", err)
	}
	return Metadata{
		ID:           id,
		Format:       Source.Name,
		Timestamp:    time.Now().Format(time.RFC3339),
		OriginalName: originalName,
		Extension:    ext,
	}, nil
}

func (api *Api) delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := api.idProvider.Validate(id); err != nil {
		svc.Error(w, r, err)
		return
	}
	api.logger.InfoContext(r.Context(), "Deleting image", "method", "DELETE", "id", id)
	sourcePath := storagePath(id, Source)
	sourceFile, metadata, err := api.store.Get(r.Context(), sourcePath)
	if err != nil {
		svc.Error(w, r, fmt.Errorf("failed to get source image: %w", err))
		return
	}
	defer sourceFile.Close()
	api.logger.DebugContext(r.Context(), "Coping source image to deleted store", "method", "DELETE", "path", sourcePath)
	err = api.deletedStore.Put(r.Context(), storagePath(id, Source), sourceFile, metadata)
	if err != nil {
		svc.Error(w, r, fmt.Errorf("failed to store deleted image: %w", err))
		return
	}
	api.logger.DebugContext(r.Context(), "Deleting image from store", "method", "DELETE", "id", id)
	err = api.store.DeleteAll(r.Context(), id)
	if err != nil {
		svc.Error(w, r, err)
		return
	}
	api.logger.DebugContext(r.Context(), "Image deleted", "method", "DELETE", "id", id)
	w.WriteHeader(http.StatusOK)
}
