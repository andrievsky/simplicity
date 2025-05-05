package svc

import (
	"log/slog"
	"net/http"
	"runtime"
	"time"
)

type StatusRespWr struct {
	http.ResponseWriter // We embed http.ResponseWriter
	status              int
}

func (w *StatusRespWr) WriteHeader(status int) {
	w.status = status // Store the status for our own use
	w.ResponseWriter.WriteHeader(status)
}

func NewLoggingMiddleware(h http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		srw := &StatusRespWr{ResponseWriter: w}
		h.ServeHTTP(srw, r)
		if srw.status >= 400 {
			logger.Error("Error status code", "status", srw.status, "path", r.RequestURI)
		}
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		logger.Debug("Request stats", "method", r.Method, "path", r.RequestURI, "status", srw.status, "time", time.Since(startTime), "alloc", bToMb(m.Alloc), "totalAlloc", bToMb(m.TotalAlloc), "sys", bToMb(m.Sys))
	})
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
