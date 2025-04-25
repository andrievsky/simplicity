package images

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"simplicity/oops"
	"simplicity/svc"
)

func (h *ImageApi) get(w http.ResponseWriter, r *http.Request) {
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
	if err == oops.KeyNotFound {
		err = h.createImageVariant(r.Context(), id, format)
		if err == nil {
			reader, metadata, err = h.store.Get(r.Context(), path)
		}
	}
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

func (h *ImageApi) createImageVariant(ctx context.Context, id string, format *Format) error {
	reader, metadata, err := h.store.Get(ctx, storagePath(id, Canonical))
	if err != nil {
		return fmt.Errorf("failed to get canonical image: %w", err)
	}
	defer reader.Close()
	tReader, err := transcodeFile(Canonical, format, reader)
	if err != nil {
		return fmt.Errorf("failed to transcode image: %w", err)
	}
	return h.store.Put(ctx, storagePath(id, format), tReader, metadata)
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
