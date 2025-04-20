package svc

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"simplicity/oops"
)

/*
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
*/
func WriteData(w http.ResponseWriter, r *http.Request, data any, status int) {
	h := w.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("ItemHandler", "Request:", r, "Error:", err.Error())
	}
}

func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("ItemHandler", "Request:", r, "Error:", err.Error())
	h := w.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	slog.Error("ItemHandler", "Request:", r, "Error:", err.Error())
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
