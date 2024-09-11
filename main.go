package main

import (
	"log"
	"net/http"

	"github.com/DuganChandler/goserver/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
	const filePathRoot = "."
	const port = "8080"

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := &apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	mux := http.NewServeMux()
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))
	mux.Handle("/app/", apiCfg.middlwareMetricsInc(handler))

    // API
	mux.HandleFunc("GET /api/healthz", getHealth)
	mux.HandleFunc("GET /api/reset", apiCfg.resetHits)

	mux.HandleFunc("GET /api/chirps", apiCfg.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirpByIDHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.createChirpsHandler)

	mux.HandleFunc("POST /api/users", apiCfg.createUsersHandler)

    mux.HandleFunc("POST /api/login", apiCfg.loginUsersHadler)

    // ADMIN
	mux.HandleFunc("GET /admin/metrics", apiCfg.getHits)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
