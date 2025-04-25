package images

import (
	"net/http"
	"simplicity/genid"
	"simplicity/storage"
)

type ImageApi struct {
	store      storage.BlobStore
	idProvider genid.Provider
}

type Image struct {
	ID       string `json:"id"`
	Location string `json:"location"`
}

const maxUploadSize = 48 * 1024 * 1024 // 48MB

func NewImageApi(store storage.BlobStore, idProvider genid.Provider) *http.ServeMux {
	router := http.NewServeMux()
	api := &ImageApi{store, idProvider}

	router.HandleFunc("GET /files/", api.list)
	router.HandleFunc("POST /upload", api.post)
	router.HandleFunc("GET /files/{id}", api.get)
	router.HandleFunc("DELETE /files/{id}", api.delete)

	return router
}

func storagePath(id string, format *Format) string {
	return storage.JoinPath(id, format.FileName())
}
