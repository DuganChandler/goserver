package main

import (
	"log"
	"net/http"
	"os"

	"github.com/DuganChandler/goserver/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	JWTSecret      string
}

func main() {
	const filePathRoot = "."
	const port = "8080"
	godotenv.Load()

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	apiCfg := &apiConfig{
		fileserverHits: 0,
		DB:             db,
		JWTSecret:      jwtSecret,
	}

	mux := http.NewServeMux()
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))
	mux.Handle("/app/", apiCfg.middlwareMetricsInc(handler))

	// API
	mux.HandleFunc("GET /api/healthz", getHealth)
	mux.HandleFunc("GET /api/reset", apiCfg.resetHits)
	mux.HandleFunc("POST /api/refresh", apiCfg.refreshTokenHandler)
	mux.HandleFunc("POST /api/revoke", apiCfg.revokeTokenHandler)

	mux.HandleFunc("GET /api/chirps", apiCfg.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirpByIDHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.createChirpsHandler)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.deleteChirpByIDHandler)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.upgradeUser)

	mux.HandleFunc("POST /api/users", apiCfg.createUsersHandler)
	mux.HandleFunc("PUT /api/users", apiCfg.updateUsersLoginHandler)

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
