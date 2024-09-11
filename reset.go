package main

import "net/http"

func (cfg *apiConfig) resetHits(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText((http.StatusOK))))
	cfg.fileserverHits = 0
}
