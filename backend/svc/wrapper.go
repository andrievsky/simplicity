package svc

import (
	"log"
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

func WrapHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		srw := &StatusRespWr{ResponseWriter: w}
		h.ServeHTTP(srw, r)
		if srw.status >= 400 { // 400+ codes are the error codes
			log.Printf("Error status code: %d when serving path: %s",
				srw.status, r.RequestURI)
		}
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		log.Printf("Request %s %s took %s [Alloc = %v MiB, TotalAlloc = %v MiB, Sys = %v MiB]", r.Method, r.RequestURI, time.Since(startTime), bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys))
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
