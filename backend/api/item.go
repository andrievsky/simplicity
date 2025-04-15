package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"simplicity/items"
	"simplicity/oops"
)

type ItemHandler struct {
	registry items.Registry
}

func NewItemHandler(registry items.Registry) *http.ServeMux {
	router := http.NewServeMux()
	handler := &ItemHandler{registry: registry}

	router.HandleFunc("GET /item", handler.list)
	router.HandleFunc("GET /item/", handler.list)
	router.HandleFunc("POST /item", handler.post)
	router.HandleFunc("GET /item/{id}", handler.get)
	router.HandleFunc("PUT /item/{id}", handler.put)
	router.HandleFunc("PATCH /item/{id}", handler.patch)
	router.HandleFunc("DELETE /item/{id}", handler.delete)

	return router
}

func (h *ItemHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.registry.List(r.Context())
	if err != nil {
		writeError(w, r, err)
		return
	}
	for i := range items {
		items[i] = ensureDefaults(items[i])
	}
	writeData(w, r, items)
}

func (h *ItemHandler) post(w http.ResponseWriter, r *http.Request) {
	var item items.Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		writeError(w, r, err)
		return
	}
	id := item.ItemMetadata.ID
	if id == "" {
		writeError(w, r, errors.New("ID is required for item creation"))
		return
	}
	err = h.registry.Create(r.Context(), id, item.ItemData)
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
	item = ensureDefaults(item)
	writeData(w, r, item)
}

func ensureDefaults(item items.Item) items.Item {
	if item.Tags == nil {
		item.Tags = []string{}
	}
	if item.Images == nil {
		item.Images = []string{}
	}
	return item
}

func (h *ItemHandler) put(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var item items.ItemData
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

func (h *ItemHandler) patch(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var fields map[string]any
	err := json.NewDecoder(r.Body).Decode(&fields)
	if err != nil {
		writeError(w, r, err)
		return
	}
	if len(fields) == 0 {
		writeError(w, r, errors.New("no fields to update"))
		return
	}
	for key := range fields {
		switch key {
		case "title", "description", "tags", "images":
			continue
		default:
			writeError(w, r, errors.New("invalid field: "+key))
			return
		}
	}

	item, err := h.registry.Read(r.Context(), id)
	if err != nil {
		writeError(w, r, err)
		return
	}
	if title, ok := fields["title"]; ok {
		item.Title = title.(string)
	}
	if description, ok := fields["description"]; ok {
		item.Description = description.(string)
	}
	if tags, ok := fields["tags"]; ok {
		item.Tags = make([]string, len(tags.([]any)))
		for i, tag := range tags.([]any) {
			item.Tags[i] = tag.(string)
		}
	}
	if images, ok := fields["images"]; ok {
		item.Images = make([]string, len(images.([]any)))
		for i, image := range images.([]any) {
			item.Images[i] = image.(string)
		}
	}
	err = h.registry.Update(r.Context(), id, item.ItemData)
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
	msg, _ := json.Marshal(map[string]string{"error": err.Error()})
	http.Error(w, string(msg), resolveErrorCode(err))
}

func resolveErrorCode(err error) int {
	if errors.Is(err, oops.KeyNotFound) {
		return http.StatusNotFound
	}
	//if errors.Is(err, oops.InvalidKey) || errors.Is(err, oops.ValidationError) || errors.Is(err, oops.KeyAlreadyExists) {
	//	return http.StatusBadRequest
	//}

	return http.StatusBadRequest
}
