package images

import (
	"fmt"
	"io"
	"net/http"
	"simplicity/svc"
)

func (h *ImageApi) post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		svc.Error(w, fmt.Errorf("failed to parse form: %w", err).Error(), http.StatusBadRequest)
		return
	}

	if r.MultipartForm == nil || len(r.MultipartForm.File) == 0 {
		http.Error(w, "missing file field", http.StatusBadRequest)
		return
	}

	if len(r.MultipartForm.File) > 1 {
		svc.Error(w, "only one file is supported", http.StatusBadRequest)
		return
	}

	fileHeader := r.MultipartForm.File["file"][0]
	if fileHeader == nil {
		svc.Error(w, "missing file field", http.StatusBadRequest)
		return
	}

	id := h.idProvider.Generate()

	metadata, err := buildMetadata(id, fileHeader)
	if err != nil {
		svc.Error(w, fmt.Errorf("failed to build metadata: %w", err).Error(), http.StatusBadRequest)
		return
	}

	sourcePath := storagePath(metadata.ID, Source)
	canonicalPath := storagePath(metadata.ID, Canonical)

	file, err := fileHeader.Open()
	if err != nil {
		svc.Error(w, "failed to open file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save the original file to storage (blocking)
	err = h.store.Put(r.Context(), sourcePath, file, metadata.Map())
	if err != nil {
		svc.Error(w, fmt.Errorf("failed to store image: %w", err).Error(), http.StatusBadRequest)
		return
	}

	tFile, err := fileHeader.Open()
	defer tFile.Close()
	tFormat := &Format{
		Name:   "",
		Ext:    metadata.Extension,
		Width:  0,
		Height: 0,
	}
	tr, err := transcodeFile(tFormat, Canonical, tFile)
	if err != nil {
		svc.Error(w, fmt.Errorf("failed to transcode file: %w", err).Error(), http.StatusBadRequest)
		return
	}

	err = h.store.Put(r.Context(), canonicalPath, tr, metadata.Map())
	if err != nil {
		svc.Error(w, fmt.Errorf("failed to store transcoded image: %w", err).Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", "/files/"+metadata.ID)
	svc.WriteData(w, r, Image{ID: metadata.ID}, http.StatusCreated)
}

func transcodeFile(in *Format, out *Format, r io.Reader) (io.Reader, error) {
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		err := Transcode(in, out, r, pw)
		if err != nil {
			pw.CloseWithError(fmt.Errorf("failed to transcode file: %w", err))
		}
	}()

	return pr, nil
}
