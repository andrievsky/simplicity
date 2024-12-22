package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"simplicity/storage"
)

type ItemHandler struct {
	registry storage.ItemRegistry
}

func NewItemHandler(registry storage.ItemRegistry) *http.ServeMux {
	router := http.NewServeMux()
	handler := &ItemHandler{registry: registry}

	router.HandleFunc("GET /item", handler.list)
	router.HandleFunc("GET /item/", handler.list)
	router.HandleFunc("POST /item/{id}", handler.post)
	router.HandleFunc("GET /item/{id}", handler.get)
	router.HandleFunc("PUT /item/{id}", handler.put)
	router.HandleFunc("DELETE /item/{id}", handler.delete)

	return router
}

func (h *ItemHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.registry.List(r.Context())
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeData(w, r, items)
}

func (h *ItemHandler) post(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var item storage.ItemData
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		writeError(w, r, err)
		return
	}
	err = h.registry.Create(r.Context(), id, item)
	if err != nil {
		writeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *ItemHandler) get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := h.registry.Read(r.Context(), id)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeData(w, r, item)
}

func (h *ItemHandler) put(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var item storage.ItemData
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		writeError(w, r, err)
		return
	}
	err = h.registry.Update(r.Context(), id, item)
	if err != nil {
		writeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *ItemHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.registry.Delete(r.Context(), id)
	if err != nil {
		writeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func writeData(w http.ResponseWriter, r *http.Request, data any) {
	h := w.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	b, err := json.Marshal(data)
	if err != nil {
		writeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		slog.Error("ItemHandler", "Request:", r, "Error:", err.Error())
	}
}

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("ItemHandler", "Request:", r, "Error:", err.Error())
	h := w.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(resolveErrorCode(err))
	msg, _ := json.Marshal(map[string]string{"error": err.Error()})
	w.Write(msg)
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func resolveErrorCode(err error) int {
	if errors.Is(err, storage.InvalidKey) || errors.Is(err, storage.ValidationError) || errors.Is(err, storage.KeyAlreadyExists) {
		return http.StatusBadRequest
	}
	if errors.Is(err, storage.KeyNotFound) {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}
