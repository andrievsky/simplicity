package images

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"simplicity/genid"
	"testing"

	"simplicity/storage"
)

func createMultipartFormFile(t *testing.T, fieldName, filename string, content []byte) (*bytes.Buffer, string) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldName, filename)
	require.NoError(t, err)
	_, err = part.Write(content)
	require.NoError(t, err)
	require.NoError(t, writer.Close())
	return body, writer.FormDataContentType()
}

func TestImageApi_HappyPath(t *testing.T) {
	store := storage.NewPrefixBlobStore(storage.NewInMemoryBlobStore(), "image/")
	idProvider, err := genid.NewSnowflakeProvider(1)
	require.NoError(t, err)
	router := NewApi(store, idProvider)

	var imageID string
	var imageData = []byte("\xFF\xD8\xFF\xE0" + "fake jpg data")

	t.Run("POST /upload", func(t *testing.T) {
		body, contentType := createMultipartFormFile(t, "file", "pic.jpg", imageData)
		req := httptest.NewRequest(http.MethodPost, "/upload", body)
		req.Header.Set("Content-Type", contentType)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		require.Equal(t, http.StatusCreated, resp.Code, "Response: %s", resp.Body.String())

		var img Image
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &img))
		assert.NotEmpty(t, img.ID)
		imageID = img.ID
	})

	t.Run("GET /files/", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/files/", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Code)

		var files []string
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &files))
		assert.Contains(t, files[0], imageID)
	})

	t.Run("GET /files/{id}?format=source", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/files/"+imageID+"?format=source", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Code, "Response: %s", resp.Body.String())
		assert.Equal(t, "image/jpeg", resp.Header().Get("Content-Type"))

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, imageData, body)
	})

	t.Run("DELETE /files/{id}", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/files/"+imageID, nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Code, "Response: %s", resp.Body.String())
	})
}

func TestImageApi_UnhappyPath(t *testing.T) {
	store := storage.NewInMemoryBlobStore()
	idProvider, err := genid.NewSnowflakeProvider(1)
	require.NoError(t, err)
	router := NewApi(store, idProvider)

	t.Run("POST /upload with no file", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		require.NoError(t, writer.Close())

		req := httptest.NewRequest(http.MethodPost, "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("GET /files/{id} with invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/files/invalid!", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("DELETE /files/{id} with invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/files/badid__", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("GET /files/{id}?format=unsupported", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/files/20240101120000_abcd1234?format=unsupported", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}
