package images

import (
	"net/http"
	"simplicity/svc"
)

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
