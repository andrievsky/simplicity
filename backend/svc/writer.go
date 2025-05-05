package svc

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"simplicity/oops"
)

func Data(w http.ResponseWriter, r *http.Request, data any, status int) {
	h := w.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Default().Error("Error encoding data", "Request:", r, "Error:", err.Error())
	}
}

func Error(w http.ResponseWriter, r *http.Request, err error) {
	ErrorWithCode(w, r, err, resolveErrorCode(err))
}

func ErrorWithCode(w http.ResponseWriter, r *http.Request, err error, code int) {
	h := w.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(resolveErrorCode(err))
	msg, err := json.Marshal(map[string]string{"error": err.Error()})
	if err != nil {
		slog.Default().Error("Error encoding error message", "Request:", r, "Error:", err.Error())
	}
	http.Error(w, string(msg), code)
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
