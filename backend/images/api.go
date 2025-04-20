package images

import (
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"simplicity/genid"
	"simplicity/storage"
	"simplicity/svc"
	"strings"
	"time"
)

type ImageApi struct {
	store      storage.BlobStore
	idProvider genid.Provider
}

type Image struct {
	ID       string `json:"id"`
	Location string `json:"location"`
}

const maxUploadSize = 48 * 1024 * 1024 // 48MB

func NewImageApi(store storage.BlobStore, idProvider genid.Provider) *http.ServeMux {
	router := http.NewServeMux()
	api := &ImageApi{store, idProvider}

	router.HandleFunc("GET /files/", api.list)
	router.HandleFunc("POST /upload", api.post)
	router.HandleFunc("GET /files/{id}", api.get)
	router.HandleFunc("DELETE /files/{id}", api.delete)

	return router
}

func (h *ImageApi) list(w http.ResponseWriter, r *http.Request) {
	seq, err := h.store.List(r.Context(), "", "/")
	if err != nil {
		svc.WriteError(w, r, err)
		return
	}

	images := make([]string, 0, len(seq))
	for _, item := range seq {
		if !item.IsObject {
			images = append(images, strings.TrimSuffix(item.Path, "/"))
		}
	}

	svc.WriteData(w, r, images, http.StatusOK)
}

func (h *ImageApi) post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		svc.Error(w, fmt.Errorf("failed to parse form: %w", err).Error(), http.StatusBadRequest)
		return
	}

	if r.MultipartForm == nil || len(r.MultipartForm.File) == 0 {
		http.Error(w, "missing file field", http.StatusBadRequest)
		return
	}

	if len(r.MultipartForm.File) > 1 {
		svc.Error(w, "only one file is supported", http.StatusBadRequest)
		return
	}

	if len(r.MultipartForm.Value) > 0 {
		svc.Error(w, "only one file is supported", http.StatusBadRequest)
		return
	}

	fileHeader := r.MultipartForm.File["file"][0]
	if fileHeader == nil {
		svc.Error(w, "missing file field", http.StatusBadRequest)
		return
	}

	id := h.idProvider.Generate()

	metadata, err := buildMetadata(id, fileHeader)
	if err != nil {
		svc.Error(w, fmt.Errorf("failed to build metadata: %w", err).Error(), http.StatusBadRequest)
		return
	}

	path := storagePath(metadata.ID, Source)

	file, err := fileHeader.Open()
	if err != nil {
		svc.Error(w, "failed to open file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = h.store.Put(r.Context(), path, file, metadata.Map())
	if err != nil {
		svc.Error(w, fmt.Errorf("failed to store image: %w", err).Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", "/files/"+metadata.ID)
	svc.WriteData(w, r, Image{ID: metadata.ID}, http.StatusCreated)
}

func (h *ImageApi) get(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get", r.URL.Path)
	id := r.PathValue("id")
	if err := h.idProvider.Validate(id); err != nil {
		svc.WriteError(w, r, err)
		return
	}

	format, err := resolveFormat(r.URL.Query().Get("format"))
	if err != nil {
		svc.WriteError(w, r, err)
		return
	}

	path := storagePath(id, format)
	reader, metadata, err := h.store.Get(r.Context(), path)
	if err != nil {
		svc.WriteError(w, r, err)
		return
	}
	defer reader.Close()
	header := w.Header()
	ext := format.Ext
	if format == Source {
		ext = getExtOrElse(metadata, ext)
	}
	header.Set("Content-Type", resolveMime(ext))
	if metadata != nil {
		for k, v := range metadata {
			header.Set("metadata-"+k, v)
		}
	}

	if _, err := io.Copy(w, reader); err != nil {
		slog.Error("ImageApi", "Request:", r, "Error:", err.Error())
		return
	}
}

func (h *ImageApi) delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.idProvider.Validate(id); err != nil {
		svc.WriteError(w, r, err)
		return
	}

	maybeFormat := r.URL.Query().Get("format")
	if maybeFormat == "" {
		err := h.store.DeleteAll(r.Context(), id)
		if err != nil {
			svc.WriteError(w, r, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	format, err := resolveFormat(r.URL.Query().Get("format"))
	if err != nil {
		svc.WriteError(w, r, err)
		return
	}
	err = h.store.DeleteAll(r.Context(), storagePath(id, format))
	if err != nil {
		svc.WriteError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func storagePath(id string, format *Format) string {
	return storage.JoinPath(id, format.FileName())
}

func buildMetadata(id string, fileHeader *multipart.FileHeader) (Metadata, error) {
	originalName := fileHeader.Filename
	ext := filepath.Ext(originalName)
	if ext == "" {
		return Metadata{}, fmt.Errorf("missing file extension")
	}
	_, err := resolveExt(ext[1:])
	if err != nil {
		return Metadata{}, fmt.Errorf("unsupported file format: %w", err)
	}
	return Metadata{
		ID:           id,
		Format:       Source.Name,
		Timestamp:    time.Now().Format(time.RFC3339),
		OriginalName: originalName,
	}, nil
}

func getExtOrElse(metadata map[string]string, fallback string) string {
	if metadata == nil {
		return fallback
	}
	m := metadata["original_name"]
	if m == "" {
		return fallback
	}
	ext := filepath.Ext(m)
	if len(ext) < 2 {
		return fallback
	}
	return ext[1:]
}
