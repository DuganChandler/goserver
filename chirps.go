package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/DuganChandler/goserver/internal/auth"
	"github.com/DuganChandler/goserver/internal/database"
)

func (cfg *apiConfig) createChirpsHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no token provided")
		return
	}

	subject, err := auth.VerifyJWT(token, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to verify jwt token")
		return
	}

	userID, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to turn subject to user id")
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters")
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	params.Body = checkBadWords(params.Body)

	chirp, err := cfg.DB.CreateChirp(params.Body, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseWithJSON(w, http.StatusCreated, database.Chirp{
		Id:       chirp.Id,
		Body:     chirp.Body,
		AuthorID: userID,
	})
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, req *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	chirps := []database.Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, database.Chirp{
			Id:       dbChirp.Id,
			Body:     dbChirp.Body,
			AuthorID: dbChirp.AuthorID,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})

	responseWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) getChirpByIDHandler(w http.ResponseWriter, req *http.Request) {
	chirpID, err := strconv.Atoi(req.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
        return
	}

	chirp, err := cfg.DB.GetChirpByID(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
        return
	}

	responseWithJSON(w, http.StatusOK, chirp)
}

func (cfg *apiConfig) deleteChirpByIDHandler(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no token provided")
		return
	}

	subject, err := auth.VerifyJWT(token, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to verify jwt token")
		return
	}

	userID, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to turn subject to user id")
		return
	}

	chirpID, err := strconv.Atoi(req.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
        return
	}

    chirp, err := cfg.DB.GetChirpByID(chirpID)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "unable to find chirp with desired id")
        return
    }

    if chirp.AuthorID != userID {
        respondWithError(w, http.StatusForbidden, "you do not have authorization to delete provided chirp")
        return
    }

    err = cfg.DB.DeleteChirpByID(chirpID)

    w.WriteHeader(http.StatusNoContent)
}

func checkBadWords(body string) string {
	var badWords = map[string]bool{"kerfuffle": true, "sharbert": true, "fornax": true}
	words := strings.Split(body, " ")
	for i, word := range words {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
