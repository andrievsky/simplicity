package items

import (
	"encoding/json"
	"errors"
	"net/http"
	"simplicity/svc"
)

type ItemApi struct {
	registry Registry
}

func NewItemHandler(registry Registry) *http.ServeMux {
	router := http.NewServeMux()
	api := &ItemApi{registry: registry}

	router.HandleFunc("GET /item", api.list)
	router.HandleFunc("GET /item/", api.list)
	router.HandleFunc("POST /item", api.post)
	router.HandleFunc("GET /item/{id}", api.get)
	router.HandleFunc("PUT /item/{id}", api.put)
	router.HandleFunc("PATCH /item/{id}", api.patch)
	router.HandleFunc("DELETE /item/{id}", api.delete)

	return router
}

func (h *ItemApi) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.registry.List(r.Context())
	if err != nil {
		svc.WriteError(w, r, err)
		return
	}
	for i := range items {
		items[i] = ensureDefaults(items[i])
	}
	svc.WriteData(w, r, items, http.StatusOK)
}

func (h *ItemApi) post(w http.ResponseWriter, r *http.Request) {
	var item Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		svc.WriteError(w, r, err)
		return
	}
	id := item.ItemMetadata.ID
	if id == "" {
		svc.WriteError(w, r, errors.New("ID is required for item creation"))
		return
	}
	err = h.registry.Create(r.Context(), id, item.ItemData)
	if err != nil {
		svc.WriteError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *ItemApi) get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := h.registry.Read(r.Context(), id)
	if err != nil {
		svc.WriteError(w, r, err)
		return
	}
	item = ensureDefaults(item)
	svc.WriteData(w, r, item, http.StatusOK)
}

func ensureDefaults(item Item) Item {
	if item.Tags == nil {
		item.Tags = []string{}
	}
	if item.Images == nil {
		item.Images = []string{}
	}
	return item
}

func (h *ItemApi) put(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var item ItemData
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		svc.WriteError(w, r, err)
		return
	}
	err = h.registry.Update(r.Context(), id, item)
	if err != nil {
		svc.WriteError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *ItemApi) patch(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var fields map[string]any
	err := json.NewDecoder(r.Body).Decode(&fields)
	if err != nil {
		svc.WriteError(w, r, err)
		return
	}
	if len(fields) == 0 {
		svc.WriteError(w, r, errors.New("no fields to update"))
		return
	}
	for key := range fields {
		switch key {
		case "title", "description", "tags", "images":
			continue
		default:
			svc.WriteError(w, r, errors.New("invalid field: "+key))
			return
		}
	}

	item, err := h.registry.Read(r.Context(), id)
	if err != nil {
		svc.WriteError(w, r, err)
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
		svc.WriteError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *ItemApi) delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.registry.Delete(r.Context(), id)
	if err != nil {
		svc.WriteError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
