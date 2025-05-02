package images

import (
	"net/http"
	"simplicity/svc"
	"strings"
)

func (h *ImageApi) list(w http.ResponseWriter, r *http.Request) {
	seq, err := h.store.List(r.Context(), "", "/")
	if err != nil {
		svc.WriteError(w, r, err)
		return
	}

	images := make([]string, 0, len(seq))
	for _, item := range seq {
		if !item.IsObject {
			images = append(images, strings.TrimSuffix(item.Key, "/"))
		}
	}

	svc.WriteData(w, r, images, http.StatusOK)
}
