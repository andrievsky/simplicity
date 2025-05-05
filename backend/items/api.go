package items

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"simplicity/genid"
	"simplicity/svc"
)

type Api struct {
	registry   Registry
	idProvider genid.Provider
	logger     *slog.Logger
}

func NewApi(registry Registry, idProvider genid.Provider, logger *slog.Logger) *http.ServeMux {
	router := http.NewServeMux()
	api := &Api{registry: registry, idProvider: idProvider, logger: logger.With("component", "items")}

	router.HandleFunc("GET /", api.list)
	router.HandleFunc("POST /", api.post)
	router.HandleFunc("GET /{id}", api.get)
	router.HandleFunc("PUT /{id}", api.put)
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
	id := api.idProvider.Generate()
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
	if err := api.idProvider.Validate(id); err != nil {
		svc.Error(w, r, err)
		return
	}
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
	if err := api.idProvider.Validate(id); err != nil {
		svc.Error(w, r, err)
		return
	}
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

func (api *Api) delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := api.idProvider.Validate(id); err != nil {
		svc.Error(w, r, err)
		return
	}
	err := api.registry.Delete(r.Context(), id)
	if err != nil {
		svc.Error(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
