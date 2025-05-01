package images

import (
	"fmt"
	"log/slog"
	"net/http"
	"simplicity/svc"
)

func (h *ImageApi) delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.idProvider.Validate(id); err != nil {
		svc.WriteError(w, r, err)
		return
	}
	sourceFile, metadata, err := h.store.Get(r.Context(), storagePath(id, Source))
	if err != nil {
		svc.WriteError(w, r, fmt.Errorf("failed to get source image: %w", err))
		return
	}
	defer sourceFile.Close()
	err = h.deletedStore.Put(r.Context(), storagePath(id, Source), sourceFile, metadata)
	if err != nil {
		svc.WriteError(w, r, fmt.Errorf("failed to store deleted image: %w", err))
		return
	}
	slog.Info("ImageApi", "Deleting image", "ID", id)
	err = h.store.DeleteAll(r.Context(), id)
	if err != nil {
		svc.WriteError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
