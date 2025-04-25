package images

import (
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
	slog.Info("ImageApi", "Deleting image", "ID", id)
	err := h.store.DeleteAll(r.Context(), id)
	if err != nil {
		svc.WriteError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
