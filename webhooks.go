package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/DuganChandler/goserver/internal/auth"
)

func (cfg *apiConfig) upgradeUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

    apiKey, err := auth.GetAPIKey(req.Header) 
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "no api key present")
        return
    }

    polkaApiKey := os.Getenv("POLKA_API_KEY")
    if apiKey != polkaApiKey {
        respondWithError(w, http.StatusUnauthorized, "incorrect api key")
        return
    }

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding request body")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.DB.UpgradeUser(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "error upgrading user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
