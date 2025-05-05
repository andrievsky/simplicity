package items

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"simplicity/svc"
)

type Api struct {
	registry Registry
	logger   *slog.Logger
}

func NewApi(registry Registry, logger *slog.Logger) *http.ServeMux {
	router := http.NewServeMux()
	api := &Api{registry: registry, logger: logger.With("component", "items")}

	router.HandleFunc("GET /", api.list)
	router.HandleFunc("POST /", api.post)
	router.HandleFunc("GET /{id}", api.get)
	router.HandleFunc("PUT /{id}", api.put)
	router.HandleFunc("PATCH /{id}", api.patch)
	router.HandleFunc("DELETE /{id}", api.delete)

	return router
}

func (api *Api) list(w http.ResponseWriter, r *http.Request) {
	items, err := api.registry.List(r.Context())
	if err != nil {
		svc.Error(w, r, err)
		return
	}
	for i := range items {
		items[i] = ensureDefaults(items[i])
	}
	svc.Data(w, r, items, http.StatusOK)
}

func (api *Api) post(w http.ResponseWriter, r *http.Request) {
	var item Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		svc.Error(w, r, err)
		return
	}
	id := item.ItemMetadata.ID
	if id == "" {
		svc.Error(w, r, errors.New("ID is required for item creation"))
		return
	}
	api.logger.Info("Creating item", "ID", id, "data", item, "method", "POST")
	err = api.registry.Create(r.Context(), id, item.ItemData)
	if err != nil {
		svc.Error(w, r, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (api *Api) get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := api.registry.Read(r.Context(), id)
	if err != nil {
		svc.Error(w, r, err)
		return
	}
	item = ensureDefaults(item)
	svc.Data(w, r, item, http.StatusOK)
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

func (api *Api) put(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var item ItemData
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		svc.Error(w, r, err)
		return
	}
	err = api.registry.Update(r.Context(), id, item)
	if err != nil {
		svc.Error(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (api *Api) patch(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var fields map[string]any
	err := json.NewDecoder(r.Body).Decode(&fields)
	if err != nil {
		svc.Error(w, r, err)
		return
	}
	if len(fields) == 0 {
		svc.Error(w, r, errors.New("no fields to update"))
		return
	}
	for key := range fields {
		switch key {
		case "title", "description", "tags", "images":
			continue
		default:
			svc.Error(w, r, errors.New("invalid field: "+key))
			return
		}
	}

	item, err := api.registry.Read(r.Context(), id)
	if err != nil {
		svc.Error(w, r, err)
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
	err = api.registry.Update(r.Context(), id, item.ItemData)
	if err != nil {
		svc.Error(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (api *Api) delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := api.registry.Delete(r.Context(), id)
	if err != nil {
		svc.Error(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
